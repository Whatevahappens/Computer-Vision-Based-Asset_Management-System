package service

import (
	"testing"
)

// ============================================================
// TC-002: Баркод давтагдашгүй байдал (Unit)
// ============================================================
func TestGenerateBarcode(t *testing.T) {
	tests := []struct {
		name     string
		dataType string
		run      func(t *testing.T)
	}{
		{
			name:     "Normal — single call returns BC- prefix",
			dataType: "Normal",
			run: func(t *testing.T) {
				bc := GenerateBarcode()
				if len(bc) < 3 || bc[:3] != "BC-" {
					t.Errorf("expected BC- prefix, got %s", bc)
				}
			},
		},
		{
			name:     "Invalid — two assets with identical data get different barcodes",
			dataType: "Invalid",
			run: func(t *testing.T) {
				bc1 := GenerateBarcode()
				bc2 := GenerateBarcode()
				if bc1 == bc2 {
					t.Errorf("expected unique barcodes, got same: %s", bc1)
				}
			},
		},
		{
			name:     "Exception — 1000 consecutive calls produce zero duplicates",
			dataType: "Exception",
			run: func(t *testing.T) {
				seen := make(map[string]bool)
				for i := 0; i < 1000; i++ {
					bc := GenerateBarcode()
					if seen[bc] {
						t.Fatalf("duplicate barcode at iteration %d: %s", i, bc)
					}
					seen[bc] = true
				}
			},
		},
		{
			name:     "Boundary — 100 calls all non-empty",
			dataType: "Boundary",
			run: func(t *testing.T) {
				for i := 0; i < 100; i++ {
					bc := GenerateBarcode()
					if bc == "" {
						t.Fatalf("empty barcode at iteration %d", i)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.run(t)
		})
	}
}

// ============================================================
// TC-003: Заавал талбарын validation (Unit)
// ============================================================
func TestValidateAssetInput(t *testing.T) {
	tests := []struct {
		name      string
		dataType  string
		input     CreateAssetRequest
		wantError bool
	}{
		{
			name:     "Normal — all fields valid",
			dataType: "Normal",
			input: CreateAssetRequest{
				AssetName:        "Dell Monitor P2422H",
				AcquisitionPrice: 450000,
				AcquisitionDate:  "2024-03-15",
				UsefulLifeMonths: 60,
			},
			wantError: false,
		},
		{
			name:     "Invalid — empty name and zero price",
			dataType: "Invalid",
			input: CreateAssetRequest{
				AssetName:        "",
				AcquisitionPrice: 0,
				AcquisitionDate:  "2024-03-15",
				UsefulLifeMonths: 60,
			},
			wantError: true,
		},
		{
			name:     "Exception — only assetName sent, others zero-value",
			dataType: "Exception",
			input: CreateAssetRequest{
				AssetName: "Keyboard",
			},
			wantError: true,
		},
		{
			name:     "Boundary — usefulLifeMonths = 0",
			dataType: "Boundary",
			input: CreateAssetRequest{
				AssetName:        "Mouse",
				AcquisitionPrice: 25000,
				AcquisitionDate:  "2024-01-01",
				UsefulLifeMonths: 0,
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAssetInput(&tt.input)
			if tt.wantError && err == nil {
				t.Error("expected validation error, got nil")
			}
			if !tt.wantError && err != nil {
				t.Errorf("expected no error, got: %v", err)
			}
		})
	}
}

// ============================================================
// TC-015: Шулуун шугамын элэгдэл (Unit)
//
// Signature: ComputeMonthly(method string, price, currentValue, lifeMonths int)
// StraightLine uses price/lifeMonths
// ============================================================
func TestComputeMonthlyStraightLine(t *testing.T) {
	tests := []struct {
		name         string
		dataType     string
		price        int
		currentValue int
		lifeMonths   int
		want         float64
		wantErr      bool
	}{
		{
			name:         "Normal — 1,200,000 / 12 months",
			dataType:     "Normal",
			price:        1200000,
			currentValue: 1200000,
			lifeMonths:   12,
			want:         100000.00,
		},
		{
			name:         "Invalid — lifeMonths=0 (division by zero)",
			dataType:     "Invalid",
			price:        1200000,
			currentValue: 1200000,
			lifeMonths:   0,
			wantErr:      true,
		},
		{
			name:         "Exception — negative price",
			dataType:     "Exception",
			price:        -500000,
			currentValue: 0,
			lifeMonths:   12,
			wantErr:      true,
		},
		{
			name:         "Boundary — price=1, lifeMonths=120",
			dataType:     "Boundary",
			price:        1,
			currentValue: 1,
			lifeMonths:   120,
			want:         0.01,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ComputeMonthly("STRAIGHT_LINE", tt.price, tt.currentValue, tt.lifeMonths)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			diff := got - tt.want
			if diff < -0.01 || diff > 0.01 {
				t.Errorf("got %.2f, want %.2f", got, tt.want)
			}
		})
	}
}

// ============================================================
// TC-016: Буурах үлдэгдлийн элэгдэл (Unit)
//
// DecliningBalance uses: currentValue * (2 / (lifeMonths/12)) / 12
// ============================================================
func TestComputeMonthlyDecliningBalance(t *testing.T) {
	tests := []struct {
		name         string
		dataType     string
		price        int
		currentValue int
		lifeMonths   int
		want         float64
		wantErr      bool
	}{
		{
			name:         "Normal — cv=2,500,000, lifeMonths=60",
			dataType:     "Normal",
			price:        2500000,
			currentValue: 2500000,
			lifeMonths:   60,
			want:         83333.33,
		},
		{
			name:         "Invalid — lifeMonths=-60 (negative)",
			dataType:     "Invalid",
			price:        2500000,
			currentValue: 1250000,
			lifeMonths:   -60,
			wantErr:      true,
		},
		{
			name:         "Exception — cv=0 (fully depreciated)",
			dataType:     "Exception",
			price:        2500000,
			currentValue: 0,
			lifeMonths:   60,
			want:         0.00,
		},
		{
			name:         "Boundary — cv=900,000, lifeMonths=36",
			dataType:     "Boundary",
			price:        2500000,
			currentValue: 900000,
			lifeMonths:   36,
			want:         50000.00,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ComputeMonthly("DECLINING_BALANCE", tt.price, tt.currentValue, tt.lifeMonths)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			diff := got - tt.want
			if diff < -0.01 || diff > 0.01 {
				t.Errorf("got %.2f, want %.2f", got, tt.want)
			}
		})
	}
}

// ============================================================
// TC-017: currentValue хамгаалалт (Unit)
// ============================================================
func TestCurrentValueProtection(t *testing.T) {
	tests := []struct {
		name         string
		dataType     string
		currentValue float64
		depreciation float64
		wantNewValue float64
		wantErr      bool
	}{
		{
			name:         "Normal — cv=500,000 - 100,000 = 400,000",
			dataType:     "Normal",
			currentValue: 500000,
			depreciation: 100000,
			wantNewValue: 400000,
		},
		{
			name:         "Invalid — depreciation > cv, should clamp to 0",
			dataType:     "Invalid",
			currentValue: 50000,
			depreciation: 100000,
			wantNewValue: 0,
		},
		{
			name:         "Exception — cv is negative (invalid state)",
			dataType:     "Exception",
			currentValue: -1,
			depreciation: 100000,
			wantErr:      true,
		},
		{
			name:         "Boundary — cv=0, depreciation=0",
			dataType:     "Boundary",
			currentValue: 0,
			depreciation: 0,
			wantNewValue: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newVal, err := ApplyDepreciation(tt.currentValue, tt.depreciation)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if newVal < 0 {
				t.Errorf("newValue went negative: %.2f", newVal)
			}
			if newVal != tt.wantNewValue {
				t.Errorf("got %.2f, want %.2f", newVal, tt.wantNewValue)
			}
		})
	}
}

// ============================================================
// TC-018: Буруу элэгдлийн арга (Unit)
// ============================================================
func TestInvalidDepreciationMethod(t *testing.T) {
	tests := []struct {
		name     string
		dataType string
		method   string
		wantErr  bool
	}{
		{
			name:     "Normal — STRAIGHT_LINE accepted",
			dataType: "Normal",
			method:   "STRAIGHT_LINE",
			wantErr:  false,
		},
		{
			name:     "Invalid — MAGIC_MATH rejected",
			dataType: "Invalid",
			method:   "MAGIC_MATH",
			wantErr:  true,
		},
		{
			name:     "Exception — numeric string rejected",
			dataType: "Exception",
			method:   "123",
			wantErr:  true,
		},
		{
			name:     "Boundary — empty string rejected",
			dataType: "Boundary",
			method:   "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ComputeMonthly(tt.method, 1000000, 1000000, 12)
			if tt.wantErr && err == nil {
				t.Error("expected error for invalid method, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
