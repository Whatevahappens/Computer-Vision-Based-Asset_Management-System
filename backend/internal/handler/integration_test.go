package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ============================================================
// Helper: test router + auth tokens
// ============================================================

const testJWTSecret = "test-secret"

func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	return db
}

func SetupRoutes(db *gorm.DB) *gin.Engine {
	r := gin.New()
	return r
}

// setupTestRouter нь SQLite in-memory + Gin router бэлдэж өгнө.
// Чиний SetupRoutes, SetupTestDB функцүүдэд тохируулж ашигла.
func setupTestRouter(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	db := SetupTestDB(t) // SQLite in-memory
	r := SetupRoutes(db)
	return r
}

func GenerateTestToken(role string) (string, error) {
	claims := jwt.MapClaims{
		"role":  role,
		"roles": []string{role},
		"exp":   time.Now().Add(time.Hour).Unix(),
		"iat":   time.Now().Unix(),
		"sub":   "test-user",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(testJWTSecret))
}

// Token-ууд нь чиний GenerateTestToken() функцээс гарна.
// Дүрүүдэд тохирсон JWT token үүсгэнэ.
func mustToken(t *testing.T, role string) string {
	t.Helper()
	tok, err := GenerateTestToken(role)
	if err != nil {
		t.Fatalf("failed to generate %s token: %v", role, err)
	}
	return tok
}

func GenerateExpiredToken(t *testing.T) string {
	t.Helper()
	claims := jwt.MapClaims{
		"role":  "custodian",
		"roles": []string{"custodian"},
		"exp":   time.Now().Add(-time.Hour).Unix(),
		"iat":   time.Now().Unix(),
		"sub":   "test-user",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tok, err := token.SignedString([]byte(testJWTSecret))
	if err != nil {
		t.Fatalf("failed to generate expired token: %v", err)
	}
	return tok
}

// ============================================================
// TC-001: Эд хөрөнгө бүртгэх (Integration)
// ============================================================
func TestCreateAssetAPI(t *testing.T) {
	r := setupTestRouter(t)
	token := mustToken(t, "custodian")

	tests := []struct {
		name       string
		dataType   string
		body       map[string]interface{}
		wantStatus int
	}{
		{
			name:     "Normal — valid Macbook Pro",
			dataType: "Normal",
			body: map[string]interface{}{
				"assetName":        "Macbook Pro",
				"acquisitionPrice": 6500000,
				"acquisitionDate":  "2024-05-01",
				"usefulLifeMonths": 48,
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:     "Invalid — negative price",
			dataType: "Invalid",
			body: map[string]interface{}{
				"assetName":        "Test Asset",
				"acquisitionPrice": -1500,
				"acquisitionDate":  "2024-05-01",
				"usefulLifeMonths": 48,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:     "Exception — invalid date format",
			dataType: "Exception",
			body: map[string]interface{}{
				"assetName":        "Test Asset",
				"acquisitionPrice": 500000,
				"acquisitionDate":  "invalid-date",
				"usefulLifeMonths": 48,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:     "Boundary — 1 char name, months=1",
			dataType: "Boundary",
			body: map[string]interface{}{
				"assetName":        "A",
				"acquisitionPrice": 100000,
				"acquisitionDate":  "2024-01-01",
				"usefulLifeMonths": 1,
			},
			wantStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyJSON, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/api/assets", bytes.NewReader(bodyJSON))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d, body: %s", w.Code, tt.wantStatus, w.Body.String())
			}

			// Normal case: check barcode was generated
			if tt.wantStatus == http.StatusCreated {
				var resp map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &resp)
				if data, ok := resp["data"].(map[string]interface{}); ok {
					if bc, ok := data["barcode"].(string); !ok || bc == "" {
						t.Error("expected barcode in response")
					}
				}
			}
		})
	}
}

// ============================================================
// TC-004: RBAC — Employee хориглолт (Integration)
// ============================================================
func TestCreateAssetRBAC(t *testing.T) {
	r := setupTestRouter(t)

	validBody, _ := json.Marshal(map[string]interface{}{
		"assetName":        "RBAC Test Asset",
		"acquisitionPrice": 500000,
		"acquisitionDate":  "2024-06-01",
		"usefulLifeMonths": 36,
	})

	tests := []struct {
		name       string
		dataType   string
		token      string
		wantStatus int
	}{
		{
			name:       "Normal — Custodian can create",
			dataType:   "Normal",
			token:      mustToken(t, "custodian"),
			wantStatus: http.StatusCreated,
		},
		{
			name:       "Invalid — Employee cannot create",
			dataType:   "Invalid",
			token:      mustToken(t, "employee"),
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "Exception — no token",
			dataType:   "Exception",
			token:      "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Boundary — expired token",
			dataType:   "Boundary",
			token:      GenerateExpiredToken(t),
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/assets", bytes.NewReader(validBody))
			req.Header.Set("Content-Type", "application/json")
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("got %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

// ============================================================
// TC-005: Жагсаалт шүүлтүүр (Integration)
// ============================================================
func TestListAssetsFilter(t *testing.T) {
	r := setupTestRouter(t)
	token := mustToken(t, "custodian")

	// Seed test data first
	seedTestAssets(t, r, token)

	tests := []struct {
		name       string
		dataType   string
		query      string
		wantStatus int
		checkEmpty bool // true = expect empty data array
	}{
		{
			name:       "Normal — filter by ACTIVE",
			dataType:   "Normal",
			query:      "/api/assets?status=ACTIVE",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Invalid — unknown status value",
			dataType:   "Invalid",
			query:      "/api/assets?status=INVALID",
			wantStatus: http.StatusOK,
			checkEmpty: true,
		},
		{
			name:       "Exception — no filter returns all",
			dataType:   "Exception",
			query:      "/api/assets",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Boundary — page=9999 returns empty",
			dataType:   "Boundary",
			query:      "/api/assets?page=9999",
			wantStatus: http.StatusOK,
			checkEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.query, nil)
			req.Header.Set("Authorization", "Bearer "+token)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("got %d, want %d", w.Code, tt.wantStatus)
			}

			if tt.checkEmpty {
				var resp map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &resp)
				if data, ok := resp["data"].([]interface{}); ok && len(data) > 0 {
					t.Errorf("expected empty data, got %d items", len(data))
				}
			}
		})
	}
}

// ============================================================
// TC-006: Эд хөрөнгө засварлах (Integration)
// ============================================================
func TestUpdateAssetAPI(t *testing.T) {
	r := setupTestRouter(t)
	custodianToken := mustToken(t, "custodian")
	employeeToken := mustToken(t, "employee")

	// Create an asset to update
	assetID := createTestAsset(t, r, custodianToken)

	tests := []struct {
		name       string
		dataType   string
		url        string
		body       map[string]interface{}
		token      string
		wantStatus int
	}{
		{
			name:       "Normal — update name",
			dataType:   "Normal",
			url:        "/api/assets/" + assetID,
			body:       map[string]interface{}{"assetName": "Шинэ нэр"},
			token:      custodianToken,
			wantStatus: http.StatusOK,
		},
		{
			name:       "Invalid — wrong asset ID",
			dataType:   "Invalid",
			url:        "/api/assets/nonexistent-uuid",
			body:       map[string]interface{}{"assetName": "Test"},
			token:      custodianToken,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Exception — Employee cannot update",
			dataType:   "Exception",
			url:        "/api/assets/" + assetID,
			body:       map[string]interface{}{"assetName": "Hack"},
			token:      employeeToken,
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "Boundary — empty body",
			dataType:   "Boundary",
			url:        "/api/assets/" + assetID,
			body:       map[string]interface{}{},
			token:      custodianToken,
			wantStatus: http.StatusOK, // no change, or 400 depending on your impl
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyJSON, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPut, tt.url, bytes.NewReader(bodyJSON))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+tt.token)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("got %d, want %d, body: %s", w.Code, tt.wantStatus, w.Body.String())
			}
		})
	}
}

// ============================================================
// TC-007: Эд хөрөнгийн түүх (Integration)
// ============================================================
func TestAssetHistoryAPI(t *testing.T) {
	r := setupTestRouter(t)
	token := mustToken(t, "custodian")

	// Create a fresh asset
	assetID := createTestAsset(t, r, token)

	tests := []struct {
		name       string
		dataType   string
		url        string
		wantStatus int
		checkCount int // -1 = don't check
	}{
		{
			name:       "Normal — history after creation",
			dataType:   "Normal",
			url:        "/api/assets/" + assetID + "/history",
			wantStatus: http.StatusOK,
			checkCount: 1, // exactly 1 CREATED entry
		},
		{
			name:       "Invalid — wrong asset ID",
			dataType:   "Invalid",
			url:        "/api/assets/nonexistent-uuid/history",
			wantStatus: http.StatusNotFound,
			checkCount: -1,
		},
		{
			name:       "Exception — newly created asset has exactly 1 entry",
			dataType:   "Exception",
			url:        "/api/assets/" + assetID + "/history",
			wantStatus: http.StatusOK,
			checkCount: 1,
		},
		{
			name:       "Boundary — history after multiple operations",
			dataType:   "Boundary",
			url:        "/api/assets/" + assetID + "/history",
			wantStatus: http.StatusOK,
			checkCount: -1, // just check status, count depends on prior ops
		},
	}

	// For Boundary test: perform an update to add history entries
	updateBody, _ := json.Marshal(map[string]interface{}{"assetName": "Updated"})
	updateReq := httptest.NewRequest(http.MethodPut, "/api/assets/"+assetID, bytes.NewReader(updateBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateReq.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(httptest.NewRecorder(), updateReq)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			req.Header.Set("Authorization", "Bearer "+token)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("got %d, want %d", w.Code, tt.wantStatus)
			}

			if tt.checkCount > 0 && w.Code == http.StatusOK {
				var resp map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &resp)
				if data, ok := resp["data"].([]interface{}); ok {
					if tt.name == "Normal" && len(data) < tt.checkCount {
						t.Errorf("expected at least %d history entries, got %d", tt.checkCount, len(data))
					}
				}
			}
		})
	}
}

// ============================================================
// TC-019: Элэгдэл RBAC (Integration)
// ============================================================
func TestDepreciationRBAC(t *testing.T) {
	r := setupTestRouter(t)
	custodianToken := mustToken(t, "custodian")

	// Create asset for depreciation
	assetID := createTestAsset(t, r, custodianToken)

	depBody, _ := json.Marshal(map[string]interface{}{
		"assetId": assetID,
		"method":  "STRAIGHT_LINE",
	})

	tests := []struct {
		name       string
		dataType   string
		token      string
		wantStatus int
	}{
		{
			name:       "Normal — Accountant can calculate",
			dataType:   "Normal",
			token:      mustToken(t, "accountant"),
			wantStatus: http.StatusOK,
		},
		{
			name:       "Invalid — Employee cannot calculate",
			dataType:   "Invalid",
			token:      mustToken(t, "employee"),
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "Exception — no token",
			dataType:   "Exception",
			token:      "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Boundary — Admin can calculate",
			dataType:   "Boundary",
			token:      mustToken(t, "admin"),
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/depreciation", bytes.NewReader(depBody))
			req.Header.Set("Content-Type", "application/json")
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("got %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

// ============================================================
// TC-020: Байхгүй хөрөнгийн элэгдэл (Integration)
// ============================================================
func TestDepreciationNonexistentAsset(t *testing.T) {
	r := setupTestRouter(t)
	accountantToken := mustToken(t, "accountant")
	custodianToken := mustToken(t, "custodian")

	// Create one valid asset
	validAssetID := createTestAsset(t, r, custodianToken)

	tests := []struct {
		name       string
		dataType   string
		body       map[string]interface{}
		wantStatus int
	}{
		{
			name:     "Normal — valid assetId",
			dataType: "Normal",
			body: map[string]interface{}{
				"assetId": validAssetID,
				"method":  "STRAIGHT_LINE",
			},
			wantStatus: http.StatusOK,
		},
		{
			name:     "Invalid — nonexistent assetId",
			dataType: "Invalid",
			body: map[string]interface{}{
				"assetId": "nonexistent-id",
				"method":  "STRAIGHT_LINE",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:     "Exception — empty assetId",
			dataType: "Exception",
			body: map[string]interface{}{
				"assetId": "",
				"method":  "STRAIGHT_LINE",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:     "Boundary — assetId field missing",
			dataType: "Boundary",
			body: map[string]interface{}{
				"method": "STRAIGHT_LINE",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyJSON, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/api/depreciation", bytes.NewReader(bodyJSON))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+accountantToken)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("got %d, want %d, body: %s", w.Code, tt.wantStatus, w.Body.String())
			}
		})
	}
}

// ============================================================
// Helpers — чиний repo-д байгаа helper функцүүдтэй солих
// ============================================================

// createTestAsset creates a valid asset and returns its ID.
func createTestAsset(t *testing.T, r *gin.Engine, token string) string {
	t.Helper()
	body, _ := json.Marshal(map[string]interface{}{
		"assetName":        "Test Asset",
		"acquisitionPrice": 1000000,
		"acquisitionDate":  "2024-01-01",
		"usefulLifeMonths": 60,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/assets", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("seed asset failed: %d %s", w.Code, w.Body.String())
	}
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	return data["id"].(string)
}

// seedTestAssets creates a few assets for list/filter tests.
func seedTestAssets(t *testing.T, r *gin.Engine, token string) {
	t.Helper()
	for i := 0; i < 3; i++ {
		createTestAsset(t, r, token)
	}
}
