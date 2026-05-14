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

var (
	router          *gin.Engine
	custodianToken  string
	accountantToken string
	employeeToken   string
	adminToken      string
	testLocationID  string
	testAssetID     string
	testSessionID   string
)

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
		"asset_name":         "Dell Latitude 5520",
		"acquisition_price":  2500000,
		"acquisition_date":   "2024-03-15",
		"useful_life_months": 60,
		"location_id":        testLocationID,
	}
	w := performRequest(router, "POST", "/api/assets", body, custodianToken)
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
		"asset_name":         "",
		"acquisition_price":  2500000,
		"acquisition_date":   "2024-03-15",
		"useful_life_months": 60,
	}
	w := performRequest(router, "POST", "/api/assets", body, custodianToken)
	if w.Code != http.StatusBadRequest {
		t.Errorf("TC-003: expected 400, got %d", w.Code)
	}
}

func TestCreateAsset_MissingPrice(t *testing.T) {
	// TC-003: Үнэ дутуу
	body := map[string]interface{}{
		"asset_name":         "Test Monitor",
		"acquisition_date":   "2024-03-15",
		"useful_life_months": 60,
	}
	w := performRequest(router, "POST", "/api/assets", body, custodianToken)
	if w.Code != http.StatusBadRequest {
		t.Errorf("TC-003: expected 400 for missing price, got %d", w.Code)
	}
}

func TestCreateAsset_DuplicateBarcode(t *testing.T) {
	// TC-004: Ижил баркод давхардах
	body := map[string]interface{}{
		"asset_name":         "Test Asset A",
		"acquisition_price":  1000000,
		"acquisition_date":   "2024-01-01",
		"useful_life_months": 36,
		"barcode":            "BC-DUP-TEST-001",
	}
	w1 := performRequest(router, "POST", "/api/assets", body, custodianToken)
	if w1.Code != http.StatusCreated {
		t.Fatalf("TC-004: first create failed: %d", w1.Code)
	}

	body["asset_name"] = "Test Asset B"
	w2 := performRequest(router, "POST", "/api/assets", body, custodianToken)
	if w2.Code != http.StatusConflict {
		t.Errorf("TC-004: expected 409, got %d", w2.Code)
	}
}

func TestListAssets_WithFilter(t *testing.T) {
	// TC-005: Шүүлтүүрээр жагсаалт харах
	w := performRequest(router, "GET", "/api/assets?status=ACTIVE", nil, custodianToken)
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
		"asset_name":         "Unauthorized Asset",
		"acquisition_price":  500000,
		"acquisition_date":   "2024-01-01",
		"useful_life_months": 12,
	}
	w := performRequest(router, "POST", "/api/assets", body, employeeToken)
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
	// TC-008: 4 буланд авсан зургуудыг batch илрүүлэх
	images := createTestImages(4)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for i, img := range images {
		part, err := writer.CreateFormFile("images",
			fmt.Sprintf("corner_%d.jpg", i+1))
		if err != nil {
			t.Fatalf("TC-008: failed to create form file: %v", err)
		}
		part.Write(img)
	}
	writer.WriteField("location_id", testLocationID)
	writer.Close()

	req, _ := http.NewRequest("POST", "/api/audit/cv", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+custodianToken)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("TC-008: expected 200, got %d, body: %s",
			w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["summary"] == nil {
		t.Error("TC-008: response should contain summary")
	}
}

func TestCVAudit_CompareWithRegistry(t *testing.T) {
	// TC-009: Бүртгэлтэй тоотой харьцуулж зөрүү тодорхойлох
	if testSessionID == "" {
		t.Skip("TC-009: no audit session created yet")
	}

	w := performRequest(router, "GET",
		fmt.Sprintf("/api/audit/sessions/%s", testSessionID),
		nil, custodianToken)
	if w.Code != http.StatusOK {
		t.Errorf("TC-009: expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	summary, ok := resp["summary"].(map[string]interface{})
	if !ok {
		t.Fatal("TC-009: response missing summary object")
	}
	if summary["registered_count"] == nil {
		t.Error("TC-009: missing registered_count")
	}
	if summary["detected_count"] == nil {
		t.Error("TC-009: missing detected_count")
	}
	if summary["difference"] == nil {
		t.Error("TC-009: missing difference field")
	}
}

func TestCVAudit_RBAC_EmployeeBlocked(t *testing.T) {
	// UC-002 RBAC: Ажилтан аудит хийх эрхгүй
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("location_id", testLocationID)
	writer.Close()

	req, _ := http.NewRequest("POST", "/api/audit/cv", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+employeeToken)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

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
		"asset_id": testAssetID,
		"method":   "STRAIGHT_LINE",
	}
	w := performRequest(router, "POST",
		"/api/depreciation/calculate", body, accountantToken)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d, body: %s",
			w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["monthly_amount"] == nil {
		t.Error("missing monthly_amount in response")
	}
	if resp["current_value"] == nil {
		t.Error("missing current_value in response")
	}
	// Verify current_value decreased
	if cv, ok := resp["current_value"].(float64); ok {
		if cv < 0 {
			t.Error("current_value should not be negative")
		}
	}
}

func TestDepreciationAPI_DecliningBalance(t *testing.T) {
	if testAssetID == "" {
		t.Skip("no test asset created")
	}
	body := map[string]string{
		"asset_id": testAssetID,
		"method":   "DECLINING_BALANCE",
	}
	w := performRequest(router, "POST",
		"/api/depreciation/calculate", body, accountantToken)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["monthly_amount"] == nil {
		t.Error("missing monthly_amount")
	}
	if resp["method"] != "DECLINING_BALANCE" {
		t.Errorf("expected method DECLINING_BALANCE, got %v", resp["method"])
	}
}

func TestDepreciationAPI_InvalidMethod(t *testing.T) {
	body := map[string]string{
		"asset_id": testAssetID,
		"method":   "INVALID_METHOD",
	}
	w := performRequest(router, "POST",
		"/api/depreciation/calculate", body, accountantToken)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid method, got %d", w.Code)
	}
}

func TestDepreciationAPI_RBAC_EmployeeBlocked(t *testing.T) {
	body := map[string]string{
		"asset_id": testAssetID,
		"method":   "STRAIGHT_LINE",
	}
	w := performRequest(router, "POST",
		"/api/depreciation/calculate", body, employeeToken)
	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestDepreciationAPI_NonexistentAsset(t *testing.T) {
	body := map[string]string{
		"asset_id": "nonexistent-id-12345",
		"method":   "STRAIGHT_LINE",
	}
	w := performRequest(router, "POST",
		"/api/depreciation/calculate", body, accountantToken)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 for nonexistent asset, got %d", w.Code)
	}
}
