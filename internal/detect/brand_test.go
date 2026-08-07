package detect_test

import (
	"testing"

	"github.com/deangrant/user-agent/internal/detect"
)

func TestBrandFromModel(t *testing.T) {
	tests := []struct {
		model string
		want  string
	}{
		{"iPhone14,2", "Apple"},
		{"iPad13,1", "Apple"},
		{"Pixel 7", "Google"},
		{"Nexus 5", "Google"},
		{"SM-G991B", "Samsung"},
		{"GT-I9300", "Samsung"},
		{"samsung galaxy", "Samsung"},
		{"Moto G", "Motorola"},
		{"Nokia 7.2", "Nokia"},
		{"ANA-NX9", "Huawei"},
		{"LYA-L29", "Huawei"},
		{"Huawei P30", "Huawei"},
		{"Redmi Note 10", "Xiaomi"},
		{"Mi 11", "Xiaomi"},
		{"POCO F3", "Xiaomi"},
		{"POCOPHONE F1", "Xiaomi"},
		{"OnePlus 9", "OnePlus"},
		{"UnknownPhone", ""},
	}
	for _, tt := range tests {
		if got := detect.BrandFromModel(tt.model); got != tt.want {
			t.Errorf("BrandFromModel(%q) = %q, want %q",
				tt.model, got, tt.want)
		}
	}
}
