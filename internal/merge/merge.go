// Package merge applies Client Hints over UA-derived analysis state.
package merge

import (
	"strings"

	"github.com/deangrant/user-agent/internal/detect"
	"github.com/deangrant/user-agent/internal/hintparse"
	"github.com/deangrant/user-agent/internal/platform"
	"github.com/deangrant/user-agent/internal/uaclass"
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
	refineIPadOS(state)
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
		if state.AgentClass == "" || state.AgentClass == uaclass.Unknown {
			state.AgentClass = uaclass.Browser
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
	return existing == "" || existing == uaclass.Unknown
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
	if name == "iOS" && preferIPadOS(state) {
		name = "iPadOS"
	}

	uaName := state.OSName
	uaEmpty := uaName == "" || uaName == uaclass.Unknown
	if !uaEmpty && name != "" {
		uaFam := osFamily(uaName)
		chFam := osFamily(name)
		if uaFam != "" && chFam != "" && uaFam != chFam {
			state.ClientHintsMismatch = true
			return
		}
	}

	if name != "" {
		// Keep UA-derived iPadOS when CH reports the broader iOS platform.
		if state.OSName != "iPadOS" || name != "iOS" {
			state.OSName = name
		}
	}
	if ver != "" {
		state.OSVersion = ver
	}
	if class != "" &&
		(state.OSClass == "" || state.OSClass == uaclass.Unknown) {
		state.OSClass = class
	}
}

// osFamily maps an OS display name to a coarse platform family.
func osFamily(name string) string {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "windows", "windows phone":
		return "windows"
	case "android":
		return "android"
	case "ios", "ipados":
		return "ios"
	case "macos", "mac os x":
		return "macos"
	case "linux", "ubuntu", "debian", "fedora", "unix":
		return "linux"
	case "chrome os", "chromeos":
		return "chromeos"
	case "":
		return ""
	default:
		return strings.ToLower(strings.TrimSpace(name))
	}
}

// preferIPadOS reports whether hints or device fields already indicate an iPad.
func preferIPadOS(state *detect.State) bool {
	if strings.Contains(strings.ToLower(state.Hints.Model), "ipad") {
		return true
	}
	if strings.Contains(strings.ToLower(state.DeviceName), "ipad") {
		return true
	}
	return false
}

// refineIPadOS promotes iOS to iPadOS after device model hints are applied.
func refineIPadOS(state *detect.State) {
	if state.OSName != "iOS" {
		return
	}
	if preferIPadOS(state) {
		state.OSName = "iPadOS"
	}
}

func applyDeviceHints(state *detect.State) {
	h := state.Hints
	if h.Model != "" {
		state.DeviceName = h.Model
		if state.DeviceBrand == "" {
			state.DeviceBrand = detect.BrandFromModel(h.Model)
		}
	}
	if len(h.FormFactors) > 0 {
		if class := formFactorClass(h.FormFactors); class != "" {
			state.DeviceClass = class
		}
	} else if h.Mobile != nil {
		if *h.Mobile {
			if state.DeviceClass == "" ||
				state.DeviceClass == uaclass.Unknown ||
				state.DeviceClass == uaclass.Desktop ||
				state.DeviceClass == uaclass.Tablet {
				state.DeviceClass = uaclass.Phone
			}
		} else if state.DeviceClass == "" ||
			state.DeviceClass == uaclass.Unknown {
			state.DeviceClass = uaclass.Desktop
		}
	}
}

func formFactorClass(factors []string) string {
	for _, f := range factors {
		switch strings.ToLower(strings.TrimSpace(f)) {
		case "mobile", "phone":
			return uaclass.Phone
		case "tablet":
			return uaclass.Tablet
		case "desktop", "computer":
			return uaclass.Desktop
		case "xr", "xr-compatible", "immersive-xr":
			return uaclass.VirtualReality
		case "automotive", "car":
			return uaclass.Car
		case "tv", "television":
			return uaclass.TV
		case "watch", "wristband":
			return uaclass.Watch
		case "ereader":
			return uaclass.EReader
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
		state.DeviceClass = uaclass.Unknown
	}
	if state.OSClass == "" {
		state.OSClass = uaclass.Unknown
	}
	if state.EngineClass == "" {
		state.EngineClass = uaclass.Unknown
	}
	if state.AgentClass == "" {
		state.AgentClass = uaclass.Unknown
	}
	if state.AgentSecurity == "" {
		state.AgentSecurity = uaclass.Unknown
	}
}
