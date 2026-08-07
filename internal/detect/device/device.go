// Package device detects device class, brand, and model.
package device

import (
	"strings"
	"unicode"

	"github.com/deangrant/user-agent/internal/data"
	"github.com/deangrant/user-agent/internal/detect"
)

// Detector identifies device class and brand/model.
type Detector struct {
	brands []data.BrandPattern
}

// New returns a device Detector.
func New() *Detector {
	c := data.MustLoad()
	return &Detector{brands: c.Brands}
}

// Detect implements detect.Detector.
func (d *Detector) Detect(state *detect.State) {
	if state == nil || state.UA == "" {
		return
	}
	lower := strings.ToLower(state.UA)

	if state.DeviceClass == "" || state.DeviceClass == "Unknown" {
		state.DeviceClass = d.classify(state, lower)
	}
	if state.DeviceBrand == "" || state.DeviceName == "" {
		d.applyBrand(state, lower)
	}
	if state.DeviceName == "" {
		d.extractAndroidModel(state, lower)
	}
	// Refine phone vs tablet for Android.
	if state.OSName == "Android" &&
		(state.DeviceClass == "Mobile" || state.DeviceClass == "Phone" ||
			state.DeviceClass == "Unknown" || state.DeviceClass == "") {
		if strings.Contains(lower, "mobile") {
			state.DeviceClass = "Phone"
		}
	}
}

func (d *Detector) classify(state *detect.State, lower string) string {
	if state.BotMatched && state.DeviceClass != "" {
		return state.DeviceClass
	}
	switch {
	case hasAll(lower, "googlebot", "mobile"):
		return "Robot Mobile"
	case strings.Contains(lower, "iphone"):
		return "Phone"
	case strings.Contains(lower, "ipad"):
		return "Tablet"
	case strings.Contains(lower, "ipod"):
		return "Mobile"
	case hasAny(lower, "smart-tv", "smarttv", "hbbtv",
		"bravia", "appletv", "googletv"):
		return "TV"
	case hasAny(lower, "chromecast", "crkey", "roku"):
		return "Set-top box"
	case containsToken(lower, "aft"):
		return "Set-top box"
	case hasAny(lower, "playstation", "xbox", "nintendo"):
		if hasAny(lower, "3ds", "new nintendo 3ds") {
			return "Handheld Game Console"
		}
		return "Game Console"
	case containsToken(lower, "watch"):
		return "Watch"
	case strings.Contains(lower, "tesla"):
		return "Car"
	case hasAny(lower, "oculus", "quest"):
		return "Virtual Reality"
	case strings.Contains(lower, "glass"):
		return "Augmented Reality"
	case hasAny(lower, "tablet", "kindle"):
		return "Tablet"
	case hasAll(lower, "android", "mobile"):
		return "Phone"
	case strings.Contains(lower, "android"):
		return "Phone"
	case strings.Contains(lower, "windows phone"):
		return "Phone"
	case strings.Contains(lower, "mobile"):
		return "Mobile"
	case hasAny(lower, "macintosh", "windows nt", "x11",
		"cros ", "linux"):
		return "Desktop"
	case state.OSClass == "Desktop":
		return "Desktop"
	case state.OSClass == "Mobile":
		return "Mobile"
	default:
		return "Unknown"
	}
}

func hasAll(s string, parts ...string) bool {
	for _, p := range parts {
		if !strings.Contains(s, p) {
			return false
		}
	}
	return true
}

func hasAny(s string, parts ...string) bool {
	for _, p := range parts {
		if strings.Contains(s, p) {
			return true
		}
	}
	return false
}

// containsToken reports whether key appears at the start of an
// alphanumeric token in s (ASCII letter/digit runs).
// "aft" matches "aftmm" and "; aft ", not "draftphone".
func containsToken(s, key string) bool {
	if key == "" || len(key) > len(s) {
		return false
	}
	key = strings.ToLower(key)
	s = strings.ToLower(s)
	for i := 0; i+len(key) <= len(s); i++ {
		if s[i:i+len(key)] != key {
			continue
		}
		if i > 0 && isASCIIAlnum(s[i-1]) {
			continue
		}
		return true
	}
	return false
}

func isASCIIAlnum(b byte) bool {
	return (b >= 'a' && b <= 'z') ||
		(b >= 'A' && b <= 'Z') ||
		(b >= '0' && b <= '9')
}

func (d *Detector) applyBrand(state *detect.State, lower string) {
	for _, b := range d.brands {
		key := b.Prefix
		if key == "" {
			key = b.Pattern
		}
		if key == "" {
			continue
		}
		if !containsToken(lower, key) {
			continue
		}
		if state.DeviceBrand == "" && b.Brand != "" {
			state.DeviceBrand = b.Brand
		}
		if state.DeviceName == "" && b.Name != "" {
			state.DeviceName = b.Name
		}
		if b.Class != "" &&
			(state.DeviceClass == "" || state.DeviceClass == "Unknown" ||
				state.DeviceClass == "Mobile" ||
				(b.Class == "Tablet" && state.DeviceClass == "Phone")) {
			state.DeviceClass = b.Class
		}
		return
	}
}

func (d *Detector) extractAndroidModel(state *detect.State, lower string) {
	if state.OSName != "Android" && !strings.Contains(lower, "android") {
		return
	}
	// Typical: Linux; Android 13; Pixel 7 Build/TQ...
	for _, c := range state.Tokens.Comments {
		cl := strings.ToLower(c)
		if strings.HasPrefix(cl, "android ") ||
			strings.HasPrefix(cl, "linux") ||
			strings.HasPrefix(cl, "u") ||
			cl == "wv" {
			continue
		}
		if strings.HasPrefix(cl, "build/") {
			continue
		}
		if strings.Contains(cl, "build/") {
			idx := strings.Index(cl, " build/")
			if idx > 0 {
				model := strings.TrimSpace(c[:idx])
				if isPlausibleModel(model) {
					state.DeviceName = model
					if state.DeviceBrand == "" {
						state.DeviceBrand = brandFromModel(model)
					}
					return
				}
			}
			continue
		}
		// Model-only comment fragment between Android version and Build.
		if isPlausibleModel(c) && !strings.Contains(cl, "android") {
			state.DeviceName = c
			if state.DeviceBrand == "" {
				state.DeviceBrand = brandFromModel(c)
			}
			return
		}
	}
}

func isPlausibleModel(s string) bool {
	s = strings.TrimSpace(s)
	if len(s) < 2 || len(s) > 60 {
		return false
	}
	lower := strings.ToLower(s)
	switch lower {
	case "mobile", "tablet", "wv", "en-us", "en-gb", "zh-cn", "k":
		return false
	}
	if strings.HasPrefix(lower, "http") {
		return false
	}
	letters := 0
	for _, r := range s {
		if unicode.IsLetter(r) {
			letters++
		}
	}
	return letters > 0
}

func brandFromModel(model string) string {
	m := strings.ToLower(model)
	switch {
	case strings.HasPrefix(m, "pixel"), strings.HasPrefix(m, "nexus"):
		return "Google"
	case strings.HasPrefix(m, "sm-"), strings.HasPrefix(m, "gt-"):
		return "Samsung"
	case strings.HasPrefix(m, "moto"):
		return "Motorola"
	case strings.HasPrefix(m, "nokia"):
		return "Nokia"
	case strings.HasPrefix(m, "redmi"), strings.HasPrefix(m, "mi "),
		strings.HasPrefix(m, "poco"):
		return "Xiaomi"
	case strings.HasPrefix(m, "oneplus"):
		return "OnePlus"
	case strings.HasPrefix(m, "huawei"):
		return "Huawei"
	default:
		return ""
	}
}
