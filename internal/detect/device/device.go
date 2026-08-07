// Package device detects device class, brand, and model.
package device

import (
	"strings"
	"unicode"

	"github.com/deangrant/user-agent/internal/data"
	"github.com/deangrant/user-agent/internal/detect"
	"github.com/deangrant/user-agent/internal/uaclass"
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

	if state.DeviceClass == "" || state.DeviceClass == uaclass.Unknown {
		state.DeviceClass = d.classify(state, lower)
	}
	if state.DeviceBrand == "" || state.DeviceName == "" {
		d.applyBrand(state, lower)
	}
	if state.DeviceName == "" {
		d.extractAndroidModel(state, lower)
	}
	// Refine Unknown/Mobile Android UAs that include "Mobile" toward Phone.
	// Do not override Watch, TV, tablets, or other specific classes.
	if state.OSName == "Android" &&
		(state.DeviceClass == uaclass.Mobile ||
			state.DeviceClass == uaclass.Unknown ||
			state.DeviceClass == "") {
		if strings.Contains(lower, "mobile") {
			state.DeviceClass = uaclass.Phone
		}
	}
}

func (d *Detector) classify(state *detect.State, lower string) string {
	if state.BotMatched && state.DeviceClass != "" {
		return state.DeviceClass
	}
	switch {
	case hasAll(lower, "googlebot", "mobile"):
		return uaclass.RobotMobile
	case strings.Contains(lower, "iphone"):
		return uaclass.Phone
	case strings.Contains(lower, "ipad"):
		return uaclass.Tablet
	case strings.Contains(lower, "ipod"):
		return uaclass.Mobile
	case hasAny(lower, "smart-tv", "smarttv", "hbbtv",
		"bravia", "appletv", "googletv", "android tv", "androidtv"):
		return uaclass.TV
	case hasAny(lower, "chromecast", "crkey", "roku"):
		return uaclass.SetTopBox
	case containsToken(lower, "aft"):
		return uaclass.SetTopBox
	case hasAny(lower, "playstation", "xbox", "nintendo"):
		if hasAny(lower, "3ds", "new nintendo 3ds") {
			return uaclass.HandheldGameConsole
		}
		return uaclass.GameConsole
	case containsToken(lower, "watch") ||
		hasAny(lower, "wear os", "wearos") ||
		strings.Contains(lower, "sm-r"):
		return uaclass.Watch
	case strings.Contains(lower, "tesla"):
		return uaclass.Car
	case hasAny(lower, "oculus", "quest"):
		return uaclass.VirtualReality
	case strings.Contains(lower, "glass"):
		return uaclass.AugmentedReality
	case hasAny(lower, "tablet", "kindle"):
		return uaclass.Tablet
	case hasAll(lower, "android", "mobile"):
		return uaclass.Phone
	case strings.Contains(lower, "android"):
		return uaclass.Phone
	case strings.Contains(lower, "windows phone"):
		return uaclass.Phone
	case strings.Contains(lower, "mobile"):
		return uaclass.Mobile
	case hasAny(lower, "macintosh", "windows nt", "x11",
		"cros ", "linux"):
		return uaclass.Desktop
	case state.OSClass == uaclass.Desktop:
		return uaclass.Desktop
	case state.OSClass == uaclass.Mobile:
		return uaclass.Mobile
	default:
		return uaclass.Unknown
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
		if b.Class != "" && brandClassOverrides(state.DeviceClass, b.Class) {
			state.DeviceClass = b.Class
		}
		return
	}
}

// brandClassOverrides reports whether a brand-derived class should
// replace the current device class (unknown/mobile, or Phone upgraded
// to Tablet/Watch/TV).
func brandClassOverrides(current, brandClass string) bool {
	if current == "" || current == uaclass.Unknown ||
		current == uaclass.Mobile {
		return true
	}
	if current != uaclass.Phone {
		return false
	}
	switch brandClass {
	case uaclass.Tablet, uaclass.Watch, uaclass.TV:
		return true
	default:
		return false
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
						state.DeviceBrand = detect.BrandFromModel(model)
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
				state.DeviceBrand = detect.BrandFromModel(c)
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
