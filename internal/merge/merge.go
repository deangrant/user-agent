// Package merge applies Client Hints over UA-derived analysis state.
package merge

import (
	"strings"

	"github.com/deangrant/user-agent/internal/detect"
	"github.com/deangrant/user-agent/internal/hintparse"
	"github.com/deangrant/user-agent/internal/platform"
	"github.com/deangrant/user-agent/internal/version"
)

// Apply enriches state with Client Hints, preferring CH when present.
func Apply(state *detect.State) {
	if state == nil {
		return
	}
	applyAgentHints(state)
	applyOSHints(state)
	applyDeviceHints(state)
	applyCPUHints(state)
	finalizeDerived(state)
}

func applyAgentHints(state *detect.State) {
	h := state.Hints
	list := h.FullVersionList
	if len(list) == 0 {
		list = h.Brands
	}
	if brand, ok := hintparse.SignificantBrand(list); ok {
		name := normalizeBrandName(brand.Brand)
		ver := brand.Version
		if ver == "" && h.FullVersion != "" {
			ver = h.FullVersion
		}
		if name != "" {
			override := shouldOverrideAgent(state.AgentName, name)
			if state.AgentName == "" || override {
				state.AgentName = name
			}
		}
		if ver != "" {
			// Prefer full version from CH over frozen UA majors.
			if state.AgentVersion == "" ||
				isFrozenChromiumVersion(state.AgentVersion) ||
				version.CompareMajor(ver, state.AgentVersion) != 0 ||
				len(ver) > len(state.AgentVersion) {
				state.AgentVersion = ver
			}
		}
		if state.AgentClass == "" || state.AgentClass == "Unknown" {
			state.AgentClass = "Browser"
		}
	}
}

func shouldOverrideAgent(existing, fromCH string) bool {
	e := strings.ToLower(existing)
	c := strings.ToLower(fromCH)
	if e == c {
		return false
	}
	// Chrome UA often says Chrome while CH says Google Chrome.
	if e == "chrome" && strings.Contains(c, "chrome") {
		return true
	}
	if e == "chromium" && c != "chromium" {
		return true
	}
	if e == "safari" || e == "firefox" {
		return false
	}
	return existing == "" || existing == "Unknown"
}

func normalizeBrandName(brand string) string {
	switch strings.ToLower(strings.TrimSpace(brand)) {
	case "google chrome":
		return "Chrome"
	case "microsoft edge":
		return "Edge"
	case "chromium":
		return "Chromium"
	default:
		return strings.TrimSpace(brand)
	}
}

func isFrozenChromiumVersion(ver string) bool {
	// Reduced UA uses X.0.0.0 style versions.
	parts := strings.Split(ver, ".")
	if len(parts) < 3 {
		return false
	}
	for i := 1; i < len(parts); i++ {
		if parts[i] != "0" {
			return false
		}
	}
	return true
}

func applyOSHints(state *detect.State) {
	h := state.Hints
	if h.Platform == "" && h.PlatformVersion == "" {
		return
	}
	name, ver, class := platform.ResolveFromCH(h.Platform, h.PlatformVersion)
	if name != "" {
		state.OSName = name
	}
	if ver != "" {
		state.OSVersion = ver
	}
	if class != "" && (state.OSClass == "" || state.OSClass == "Unknown") {
		state.OSClass = class
	}
}

func applyDeviceHints(state *detect.State) {
	h := state.Hints
	if h.Model != "" {
		state.DeviceName = h.Model
		if state.DeviceBrand == "" {
			state.DeviceBrand = guessBrand(h.Model)
		}
	}
	if len(h.FormFactors) > 0 {
		if class := formFactorClass(h.FormFactors); class != "" {
			state.DeviceClass = class
		}
	} else if h.Mobile != nil {
		if *h.Mobile {
			if state.DeviceClass == "" ||
				state.DeviceClass == "Unknown" ||
				state.DeviceClass == "Desktop" {
				state.DeviceClass = "Phone"
			}
		} else if state.DeviceClass == "" || state.DeviceClass == "Unknown" {
			state.DeviceClass = "Desktop"
		}
	}
}

func formFactorClass(factors []string) string {
	for _, f := range factors {
		switch strings.ToLower(strings.TrimSpace(f)) {
		case "mobile", "phone":
			return "Phone"
		case "tablet":
			return "Tablet"
		case "desktop", "computer":
			return "Desktop"
		case "xr", "xr-compatible", "immersive-xr":
			return "Virtual Reality"
		case "automotive", "car":
			return "Car"
		case "tv", "television":
			return "TV"
		case "watch", "wristband":
			return "Watch"
		case "ereader":
			return "eReader"
		}
	}
	return ""
}

func applyCPUHints(state *detect.State) {
	h := state.Hints
	var parts []string
	if h.Arch != "" {
		parts = append(parts, h.Arch)
	}
	if h.Bitness != "" {
		parts = append(parts, h.Bitness+"-bit")
	}
	if len(parts) > 0 {
		state.DeviceCPU = strings.Join(parts, " ")
	}
}

func finalizeDerived(state *detect.State) {
	an, av := state.AgentName, state.AgentVersion
	state.AgentVersionMajor = version.Major(av)
	state.AgentNameVersion = version.NameVersion(an, av)
	state.AgentNameVersionMajor = version.NameVersionMajor(an, av)

	en, ev := state.EngineName, state.EngineVersion
	state.EngineVersionMajor = version.Major(ev)
	state.EngineNameVersion = version.NameVersion(en, ev)
	state.EngineNameVersionMajor = version.NameVersionMajor(en, ev)

	state.OSNameVersion = version.NameVersion(
		state.OSName, state.OSVersion)

	if state.DeviceClass == "" {
		state.DeviceClass = "Unknown"
	}
	if state.OSClass == "" {
		state.OSClass = "Unknown"
	}
	if state.EngineClass == "" {
		state.EngineClass = "Unknown"
	}
	if state.AgentClass == "" {
		state.AgentClass = "Unknown"
	}
	if state.AgentSecurity == "" {
		state.AgentSecurity = "Unknown"
	}
}

func guessBrand(model string) string {
	m := strings.ToLower(model)
	switch {
	case strings.HasPrefix(m, "iphone"), strings.HasPrefix(m, "ipad"),
		strings.HasPrefix(m, "ipod"), strings.HasPrefix(m, "mac"):
		return "Apple"
	case strings.HasPrefix(m, "pixel"):
		return "Google"
	case strings.HasPrefix(m, "sm-"), strings.HasPrefix(m, "samsung"):
		return "Samsung"
	case strings.HasPrefix(m, "nexus"):
		return "Google"
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
