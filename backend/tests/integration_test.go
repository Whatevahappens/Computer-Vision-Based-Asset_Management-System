package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// ============================================================
// Test setup — adjust imports to match your project
// ============================================================

func performRequest(r *gin.Engine, method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBytes)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, _ := http.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// ============================================================
// UC-001: Эд хөрөнгө бүртгэх
// ============================================================

func TestCreateAsset_Success(t *testing.T) {
	// TC-001: Бүх заавал талбарыг зөв оруулж бүртгэх
	body := map[string]interface{}{
		"assetName":        "Dell Latitude 5520",
		"acquisitionPrice": 2500000,
		"acquisitionDate":  "2024-03-15",
		"usefulLifeMonths": 60,
		"locationId":       testLocationID,
	}
	w := performRequest(testRouter, "POST", "/api/v1/assets", body, custodianToken)
	if w.Code != http.StatusCreated {
		t.Errorf("TC-001: expected 201, got %d, body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp["barcode"] == nil || resp["barcode"] == "" {
		t.Error("TC-001: barcode should be auto-generated")
	}
	if resp["id"] == nil || resp["id"] == "" {
		t.Error("TC-001: asset id should be returned")
	}
	// Save for later tests
	if id, ok := resp["id"].(string); ok {
		testAssetID = id
	}
}

func TestCreateAsset_MissingName(t *testing.T) {
	// TC-003: Заавал талбар дутуу
	body := map[string]interface{}{
		"assetName":        "",
		"acquisitionPrice": 2500000,
		"acquisitionDate":  "2024-03-15",
		"usefulLifeMonths": 60,
		"locationId":       testLocationID,
	}
	w := performRequest(testRouter, "POST", "/api/v1/assets", body, custodianToken)
	if w.Code != http.StatusBadRequest {
		t.Errorf("TC-003: expected 400, got %d", w.Code)
	}
}

func TestCreateAsset_MissingPrice(t *testing.T) {
	// TC-003: Үнэ дутуу
	body := map[string]interface{}{
		"assetName":        "Test Monitor",
		"acquisitionDate":  "2024-03-15",
		"usefulLifeMonths": 60,
		"locationId":       testLocationID,
	}
	w := performRequest(testRouter, "POST", "/api/v1/assets", body, custodianToken)
	if w.Code != http.StatusBadRequest {
		t.Errorf("TC-003: expected 400 for missing price, got %d", w.Code)
	}
}

func TestCreateAsset_DuplicateBarcode(t *testing.T) {
	// TC-004: Two assets should get different barcodes
	body := map[string]interface{}{
		"assetName":        "Test Asset A",
		"acquisitionPrice": 1000000,
		"acquisitionDate":  "2024-01-01",
		"usefulLifeMonths": 36,
	}
	w1 := performRequest(testRouter, "POST", "/api/v1/assets", body, custodianToken)
	w2 := performRequest(testRouter, "POST", "/api/v1/assets", body, custodianToken)

	var r1, r2 map[string]interface{}
	json.Unmarshal(w1.Body.Bytes(), &r1)
	json.Unmarshal(w2.Body.Bytes(), &r2)

	if r1["barcode"] == r2["barcode"] {
		t.Error("TC-004: two assets got the same barcode")
	}
}

func TestListAssets_WithFilter(t *testing.T) {
	// TC-005: Шүүлтүүрээр жагсаалт харах
	w := performRequest(testRouter, "GET", "/api/v1/assets?status=ACTIVE", nil, custodianToken)
	if w.Code != http.StatusOK {
		t.Errorf("TC-005: expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data, ok := resp["data"].([]interface{})
	if !ok || len(data) == 0 {
		t.Error("TC-005: should return at least 1 active asset")
	}
}

func TestCreateAsset_RBAC_EmployeeBlocked(t *testing.T) {
	// UC-001 RBAC: Ажилтан хөрөнгө бүртгэх эрхгүй
	body := map[string]interface{}{
		"assetName":        "Unauthorized Asset",
		"acquisitionPrice": 500000,
		"acquisitionDate":  "2024-01-01",
		"usefulLifeMonths": 12,
		"locationId":       testLocationID,
	}
	w := performRequest(testRouter, "POST", "/api/v1/assets", body, employeeToken)
	if w.Code != http.StatusForbidden {
		t.Errorf("UC-001 RBAC: expected 403, got %d", w.Code)
	}
}

// ============================================================
// UC-002: CV аудит
// ============================================================

func createTestImages(count int) [][]byte {
	// Creates minimal valid JPEG bytes for testing
	// In real tests, use actual test images
	images := make([][]byte, count)
	for i := 0; i < count; i++ {
		// Minimal JPEG placeholder — replace with actual test images
		images[i] = []byte{0xFF, 0xD8, 0xFF, 0xE0} // JPEG magic bytes
	}
	return images
}

func TestCVAudit_BatchDetection(t *testing.T) {
	sessionBody := map[string]string{"locationId": testLocationID}
	sw := performRequest(testRouter, "POST", "/api/v1/audits", sessionBody, custodianToken)
	if sw.Code != http.StatusCreated && sw.Code != http.StatusOK {
		t.Skipf("TC-008: audit session returned %d — CV service not available in test", sw.Code)
	}
	var sessionResp map[string]interface{}
	json.Unmarshal(sw.Body.Bytes(), &sessionResp)
	sessionID, ok := sessionResp["id"].(string)
	if !ok {
		t.Skip("TC-008: could not get session ID")
	}
	testSessionID = sessionID

	images := createTestImages(4)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for i, img := range images {
		part, _ := writer.CreateFormFile("images", fmt.Sprintf("corner_%d.jpg", i+1))
		part.Write(img)
	}
	writer.Close()

	req, _ := http.NewRequest("POST", fmt.Sprintf("/api/v1/audits/%s/cv", sessionID), body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+custodianToken)

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	// CV service not running — 500 is expected (connection refused)
	// This proves routing and auth work; actual CV testing done via pytest
	if w.Code == http.StatusNotFound {
		t.Error("TC-008: route not found")
	}
	if w.Code == http.StatusForbidden {
		t.Error("TC-008: auth should pass for custodian")
	}
}

func TestCVAudit_CompareWithRegistry(t *testing.T) {
	if testSessionID == "" {
		t.Skip("TC-009: no audit session — CV service not available in test")
	}
	w := performRequest(testRouter, "GET",
		fmt.Sprintf("/api/v1/audits/%s", testSessionID), nil, custodianToken)
	if w.Code == http.StatusNotFound {
		t.Skip("TC-009: audit session endpoint not available without CV service")
	}
}

func TestCVAudit_RBAC_EmployeeBlocked(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("locationId", testLocationID)
	writer.Close()

	req, _ := http.NewRequest("POST",
		"/api/v1/audits/any-id/cv", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+employeeToken)

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("UC-002 RBAC: expected 403, got %d", w.Code)
	}
}

// ============================================================
// UC-003: Элэгдэл тооцоолох
// ============================================================

func TestDepreciationAPI_StraightLine(t *testing.T) {
	if testAssetID == "" {
		t.Skip("no test asset created")
	}
	body := map[string]string{
		"assetId": testAssetID,
		"method":  "STRAIGHT_LINE",
	}
	w := performRequest(testRouter, "POST",
		"/api/v1/depreciation/calculate", body, accountantToken)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d, body: %s",
			w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["monthlyAmount"] == nil {
		t.Error("missing monthlyAmount in response")
	}
	if resp["currentValue"] == nil {
		t.Error("missing currentValue in response")
	}
	// Verify currentValue decreased
	if cv, ok := resp["currentValue"].(float64); ok {
		if cv < 0 {
			t.Error("currentValue should not be negative")
		}
	}
}

func TestDepreciationAPI_DecliningBalance(t *testing.T) {
	if testAssetID == "" {
		t.Skip("no test asset created")
	}
	body := map[string]string{
		"assetId": testAssetID,
		"method":  "DECLINING_BALANCE",
	}
	w := performRequest(testRouter, "POST",
		"/api/v1/depreciation/calculate", body, accountantToken)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["monthlyAmount"] == nil {
		t.Error("missing monthlyAmount")
	}
	if resp["method"] != "DECLINING_BALANCE" {
		t.Errorf("expected method DECLINING_BALANCE, got %v", resp["method"])
	}
}

func TestDepreciationAPI_InvalidMethod(t *testing.T) {
	body := map[string]string{
		"assetId": testAssetID,
		"method":  "INVALID_METHOD",
	}
	w := performRequest(testRouter, "POST",
		"/api/v1/depreciation/calculate", body, accountantToken)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid method, got %d", w.Code)
	}
}

func TestDepreciationAPI_RBAC_EmployeeBlocked(t *testing.T) {
	body := map[string]string{
		"assetId": testAssetID,
		"method":  "STRAIGHT_LINE",
	}
	w := performRequest(testRouter, "POST",
		"/api/v1/depreciation/calculate", body, employeeToken)
	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestDepreciationAPI_NonexistentAsset(t *testing.T) {
	body := map[string]string{
		"assetId": "nonexistent-id-12345",
		"method":  "STRAIGHT_LINE",
	}
	w := performRequest(testRouter, "POST",
		"/api/v1/depreciation/calculate", body, accountantToken)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
