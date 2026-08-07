// Package os detects operating systems from User-Agent tokens.
package os

import (
	"strings"

	"github.com/deangrant/user-agent/internal/detect"
	"github.com/deangrant/user-agent/internal/platform"
	"github.com/deangrant/user-agent/internal/uaclass"
	"github.com/deangrant/user-agent/internal/version"
)

// Detector identifies the operating system.
type Detector struct{}

// New returns an OS Detector.
func New() *Detector { return &Detector{} }

// Detect implements detect.Detector.
func (d *Detector) Detect(state *detect.State) {
	if state == nil || state.UA == "" {
		return
	}
	if state.OSName != "" {
		return
	}

	ua := state.UA
	lower := strings.ToLower(ua)

	switch {
	case strings.Contains(lower, "windows phone"):
		ver := findAfter(lower, ua, "windows phone ")
		state.SetOS(uaclass.Mobile, "Windows Phone", version.FirstNumber(ver))
	case strings.Contains(lower, "windows nt"):
		nt := findAfter(lower, ua, "windows nt ")
		nt = trimToken(nt)
		mapped := platform.MapWindowsNT(nt)
		state.SetOS(uaclass.Desktop, "Windows", mapped)
		setCPUFromUA(state, lower)
	case strings.Contains(lower, "windows "):
		ver := findAfter(lower, ua, "windows ")
		state.SetOS(uaclass.Desktop, "Windows", trimToken(ver))
	case strings.Contains(lower, "android"):
		ver := findAfter(lower, ua, "android ")
		ver = trimToken(ver)
		state.SetOS(uaclass.Mobile, "Android", version.NormalizeSeparators(ver))
		if i := strings.Index(lower, "build/"); i >= 0 {
			build := trimToken(ua[i+6:])
			state.OSVersionBuild = build
		}
	case strings.Contains(lower, "ipad") &&
		(strings.Contains(lower, "cpu os ") ||
			strings.Contains(lower, "os ")):
		ver := iosVersion(lower, ua)
		state.SetOS(uaclass.Mobile, "iPadOS", ver)
	case strings.Contains(lower, "cpu iphone os"),
		strings.Contains(lower, "cpu os "),
		strings.Contains(lower, "iphone os"):
		ver := iosVersion(lower, ua)
		state.SetOS(uaclass.Mobile, "iOS", ver)
	case strings.Contains(lower, "mac os x"), strings.Contains(lower, "macos"):
		ver := findMacVersion(lower, ua)
		name, out := platform.ResolveMacOS(ver)
		state.SetOS(uaclass.Desktop, name, out)
	case strings.Contains(lower, "cros "):
		ver := chromeOSVersion(state)
		state.SetOS(uaclass.Desktop, "Chrome OS", ver)
	case strings.Contains(lower, "harmonyos"):
		ver := findAfter(lower, ua, "harmonyos ")
		state.SetOS(uaclass.Mobile, "HarmonyOS", trimToken(ver))
	case strings.Contains(lower, "watchos"):
		ver := findAfter(lower, ua, "watchos ")
		state.SetOS(uaclass.Mobile, "watchOS",
			version.NormalizeSeparators(trimToken(ver)))
	case strings.Contains(lower, "tvos") ||
		strings.Contains(lower, "apple tv"):
		ver := findAfter(lower, ua, "tvos ")
		state.SetOS(uaclass.Embedded, "tvOS",
			version.NormalizeSeparators(trimToken(ver)))
	case strings.Contains(lower, "freebsd"):
		state.SetOS(uaclass.Desktop, "FreeBSD", "")
	case strings.Contains(lower, "openbsd"):
		state.SetOS(uaclass.Desktop, "OpenBSD", "")
	case strings.Contains(lower, "netbsd"):
		state.SetOS(uaclass.Desktop, "NetBSD", "")
	case strings.Contains(lower, "sunos"):
		state.SetOS(uaclass.Desktop, "SunOS", "")
	case strings.Contains(lower, "webos") || strings.Contains(lower, "web0s"):
		state.SetOS(uaclass.Embedded, "webOS", "")
	case strings.Contains(lower, "tizen"):
		ver := findAfter(lower, ua, "tizen ")
		state.SetOS(uaclass.Embedded, "Tizen", trimToken(ver))
	case strings.Contains(lower, "kaios"):
		ver := findAfter(lower, ua, "kaios/")
		if ver == "" {
			ver = findAfter(lower, ua, "kaios ")
		}
		state.SetOS(uaclass.Mobile, "KaiOS", trimToken(ver))
	case strings.Contains(lower, "fedora"):
		state.SetOS(uaclass.Desktop, "Fedora", "")
	case strings.Contains(lower, "ubuntu"):
		state.SetOS(uaclass.Desktop, "Ubuntu", "")
	case strings.Contains(lower, "debian"):
		state.SetOS(uaclass.Desktop, "Debian", "")
	case strings.Contains(lower, "linux"):
		state.SetOS(uaclass.Desktop, "Linux", "")
		setCPUFromUA(state, lower)
	case strings.Contains(lower, "x11"):
		state.SetOS(uaclass.Desktop, "Unix", "")
	case state.BotMatched:
		// leave unknown unless already set
	}
}

func iosVersion(lower, ua string) string {
	for _, marker := range []string{
		"cpu iphone os ", "cpu os ", "iphone os ",
	} {
		if strings.Contains(lower, marker) {
			ver := findAfter(lower, ua, marker)
			return version.NormalizeSeparators(trimToken(ver))
		}
	}
	return ""
}

func findMacVersion(lower, ua string) string {
	for _, marker := range []string{"mac os x ", "macos "} {
		if strings.Contains(lower, marker) {
			return trimToken(findAfter(lower, ua, marker))
		}
	}
	return ""
}

func chromeOSVersion(state *detect.State) string {
	// CrOS token often in comments: CrOS x86_64 14541.0.0
	for _, c := range state.Tokens.Comments {
		if strings.HasPrefix(strings.ToLower(c), "cros ") {
			parts := strings.Fields(c)
			if len(parts) >= 3 {
				return parts[2]
			}
		}
	}
	return ""
}

func findAfter(lower, original, marker string) string {
	i := strings.Index(lower, marker)
	if i < 0 {
		return ""
	}
	return original[i+len(marker):]
}

func trimToken(s string) string {
	s = strings.TrimSpace(s)
	for i, r := range s {
		if r == ';' || r == ')' || r == '(' || r == ' ' || r == ',' {
			return strings.TrimSpace(s[:i])
		}
	}
	return s
}

func setCPUFromUA(state *detect.State, lower string) {
	if state.DeviceCPU != "" {
		return
	}
	switch {
	case strings.Contains(lower, "arm64") || strings.Contains(lower, "aarch64"):
		state.DeviceCPU = "arm64"
	case strings.Contains(lower, "win64") || strings.Contains(lower, "x64") ||
		strings.Contains(lower, "x86_64") || strings.Contains(lower, "amd64"):
		state.DeviceCPU = "x86_64"
	case strings.Contains(lower, "wow64"):
		state.DeviceCPU = "x86_64"
	case strings.Contains(lower, "i686") || strings.Contains(lower, "i386"):
		state.DeviceCPU = "x86"
	case strings.Contains(lower, "arm"):
		state.DeviceCPU = "arm"
	}
}
