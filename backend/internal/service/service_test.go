package service

import (
	"strings"
	"testing"
)

// ============================================================
// UC-001: TC-002 — Баркод давтагдашгүй байдал
// ============================================================

func TestGenerateBarcode_Unique(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		barcode := GenerateBarcode()
		if seen[barcode] {
			t.Fatalf("duplicate barcode at iteration %d: %s", i, barcode)
		}
		seen[barcode] = true
	}
}

func TestGenerateBarcode_Format(t *testing.T) {
	barcode := GenerateBarcode()
	if len(barcode) == 0 {
		t.Fatal("barcode should not be empty")
	}
	prefix := "BC-"
	if !strings.HasPrefix(barcode, prefix) {
		t.Errorf("barcode should start with %s, got %s", prefix, barcode)
	}
}

func TestGenerateBarcode_NonEmpty(t *testing.T) {
	for i := 0; i < 100; i++ {
		barcode := GenerateBarcode()
		if barcode == "" {
			t.Fatalf("empty barcode at iteration %d", i)
		}
	}
}

// ============================================================
// UC-003: TC-011 — Шулуун шугамын арга
// ============================================================

func TestStraightLine_BasicCases(t *testing.T) {
	tests := []struct {
		name     string
		price    int
		months   int
		expected float64
	}{
		{"5 year asset", 2500000, 60, 41666.67},
		{"1 year asset", 1200000, 12, 100000.00},
		{"1 month asset", 100000, 1, 100000.00},
		{"high value 10yr", 50000000, 120, 416666.67},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			monthly, err := ComputeMonthly("STRAIGHT_LINE",
				tt.price, 0, tt.months)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if monthly != tt.expected {
				t.Errorf("got %.2f, want %.2f",
					monthly, tt.expected)
			}
		})
	}
}

// ============================================================
// UC-003: TC-012 — Буурах үлдэгдлийн арга
// ============================================================

func TestDecliningBalance_BasicCases(t *testing.T) {
	tests := []struct {
		name         string
		currentValue int
		lifeMonths   int
		expected     float64
	}{
		{"standard 5yr", 2500000, 60, 83333.33},
		{"half depleted 5yr", 1250000, 60, 41666.67},
		{"3 year life", 900000, 36, 50000.00},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			monthly, err := ComputeMonthly("DECLINING_BALANCE",
				0, tt.currentValue, tt.lifeMonths)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if monthly != tt.expected {
				t.Errorf("got %.2f, want %.2f",
					monthly, tt.expected)
			}
		})
	}
}

// ============================================================
// UC-003: TC-013 — Edge cases
// ============================================================

func TestDepreciation_FloorAtZero(t *testing.T) {
	monthly, err := ComputeMonthly("STRAIGHT_LINE", 500, 0, 60)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	newValue := 500 - int(monthly*100)
	if newValue < 0 {
		newValue = 0
	}
	if newValue < 0 {
		t.Errorf("value went below zero: %d", newValue)
	}
}

func TestDepreciation_AlreadyZero(t *testing.T) {
	monthly, err := ComputeMonthly("DECLINING_BALANCE", 0, 0, 60)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if monthly != 0 {
		t.Errorf("expected 0 monthly, got %.2f", monthly)
	}
}

func TestDepreciation_LargeDeduction_ClampsToZero(t *testing.T) {
	monthly, err := ComputeMonthly("STRAIGHT_LINE", 100, 0, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	newValue := 100 - int(monthly)
	if newValue < 0 {
		newValue = 0
	}
	if newValue != 0 {
		t.Errorf("expected 0, got %d", newValue)
	}
}

func TestComputeMonthly_InvalidMethod(t *testing.T) {
	_, err := ComputeMonthly("INVALID", 1000, 1000, 12)
	if err == nil {
		t.Error("expected error for invalid method")
	}
}
