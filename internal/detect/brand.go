package detect

import "strings"

// BrandFromModel guesses a device brand from a model string
// using common manufacturer prefixes.
func BrandFromModel(model string) string {
	m := strings.ToLower(model)
	switch {
	case strings.HasPrefix(m, "iphone"), strings.HasPrefix(m, "ipad"),
		strings.HasPrefix(m, "ipod"), strings.HasPrefix(m, "mac"):
		return "Apple"
	case strings.HasPrefix(m, "pixel"), strings.HasPrefix(m, "nexus"):
		return "Google"
	case strings.HasPrefix(m, "sm-"), strings.HasPrefix(m, "gt-"),
		strings.HasPrefix(m, "samsung"):
		return "Samsung"
	case strings.HasPrefix(m, "moto"):
		return "Motorola"
	case strings.HasPrefix(m, "nokia"):
		return "Nokia"
	case strings.HasPrefix(m, "huawei"), strings.HasPrefix(m, "ana-"),
		strings.HasPrefix(m, "lya-"):
		return "Huawei"
	case strings.HasPrefix(m, "redmi"), strings.HasPrefix(m, "mi "),
		strings.HasPrefix(m, "pocophone"), strings.HasPrefix(m, "poco"):
		return "Xiaomi"
	case strings.HasPrefix(m, "oneplus"):
		return "OnePlus"
	default:
		return ""
	}
}
