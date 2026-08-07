// Package platform maps Client Hints platform versions to OS names.
package platform

import (
	"strconv"
	"strings"

	"github.com/deangrant/user-agent/internal/uaclass"
	"github.com/deangrant/user-agent/internal/version"
)

// ResolveWindows maps Sec-CH-UA-Platform-Version to a Windows version.
// See Microsoft Edge docs on detecting Windows 11 via Client Hints:
// https://learn.microsoft.com/microsoft-edge/web-platform/how-to-detect-win11
func ResolveWindows(platformVersion string) (name, ver string) {
	platformVersion = strings.TrimSpace(platformVersion)
	if platformVersion == "" {
		return "Windows", ""
	}
	major := version.Major(platformVersion)
	n, err := strconv.Atoi(major)
	if err != nil {
		return "Windows", platformVersion
	}
	switch {
	case n >= 13:
		return "Windows", "11"
	case n == 0:
		// Historical mapping used by Chromium for older Windows.
		return resolveWindowsZero(platformVersion)
	default:
		// 1–12 historically map to Windows 10.
		return "Windows", "10"
	}
}

func resolveWindowsZero(platformVersion string) (name, ver string) {
	parts := strings.Split(platformVersion, ".")
	if len(parts) < 2 {
		return "Windows", "10"
	}
	minor, _ := strconv.Atoi(parts[1])
	switch minor {
	case 1:
		return "Windows", "7"
	case 2:
		return "Windows", "8"
	case 3:
		return "Windows", "8.1"
	default:
		return "Windows", "10"
	}
}

// ResolveMacOS maps a macOS / Mac OS X version string to display form.
func ResolveMacOS(ver string) (name, out string) {
	ver = version.NormalizeSeparators(strings.TrimSpace(ver))
	if ver == "" {
		return "macOS", ""
	}
	major := version.Major(ver)
	n, _ := strconv.Atoi(major)
	// Modern macOS versions are 11+.
	if n >= 11 {
		return "macOS", ver
	}
	return "Mac OS X", ver
}

// ResolveFromCH picks OS name/version from Client Hints platform fields.
func ResolveFromCH(platform, platformVersion string) (name, ver, class string) {
	p := strings.TrimSpace(platform)
	switch strings.ToLower(p) {
	case "windows":
		name, ver = ResolveWindows(platformVersion)
		return name, ver, uaclass.Desktop
	case "macos", "mac os x":
		name, ver = ResolveMacOS(platformVersion)
		return name, ver, uaclass.Desktop
	case "android":
		ver = strings.TrimSpace(platformVersion)
		if ver == "" {
			return "Android", "", uaclass.Mobile
		}
		return "Android", version.FirstNumber(ver), uaclass.Mobile
	case "ios":
		ver = strings.TrimSpace(platformVersion)
		return "iOS", version.NormalizeSeparators(ver), uaclass.Mobile
	case "ipados":
		ver = strings.TrimSpace(platformVersion)
		return "iPadOS", version.NormalizeSeparators(ver), uaclass.Mobile
	case "linux":
		return "Linux", strings.TrimSpace(platformVersion), uaclass.Desktop
	case "chrome os", "chromeos":
		return "Chrome OS", strings.TrimSpace(platformVersion), uaclass.Desktop
	case "":
		return "", "", ""
	default:
		return p, strings.TrimSpace(platformVersion), ""
	}
}

// MapWindowsNT maps Windows NT token versions to marketing versions.
func MapWindowsNT(nt string) string {
	nt = version.NormalizeSeparators(strings.TrimSpace(nt))
	switch nt {
	case "5.0":
		return "2000"
	case "5.1", "5.2":
		return "XP"
	case "6.0":
		return "Vista"
	case "6.1":
		return "7"
	case "6.2":
		return "8"
	case "6.3":
		return "8.1"
	case "10.0":
		return "10"
	default:
		return nt
	}
}
