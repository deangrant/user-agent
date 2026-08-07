// Package engine detects layout/rendering engines.
package engine

import (
	"strings"

	"github.com/deangrant/user-agent/internal/detect"
)

// Detector identifies the layout engine.
type Detector struct{}

// New returns an engine Detector.
func New() *Detector { return &Detector{} }

// Detect implements detect.Detector.
func (d *Detector) Detect(state *detect.State) {
	if state == nil || state.UA == "" {
		return
	}
	if state.EngineName != "" {
		return
	}

	lower := strings.ToLower(state.UA)

	switch {
	case strings.Contains(lower, "trident/"):
		p, _ := state.Tokens.FindProduct("Trident")
		state.SetEngine(engineClass(state), "Trident", p.Version)
	case strings.Contains(lower, "edge/"):
		// Legacy EdgeHTML
		if !strings.Contains(lower, "edg/") {
			p, _ := state.Tokens.FindProduct("Edge")
			state.SetEngine(engineClass(state), "EdgeHTML", p.Version)
		} else {
			setBlink(state)
		}
	case strings.Contains(lower, "goanna/"):
		p, _ := state.Tokens.FindProduct("Goanna")
		state.SetEngine(engineClass(state), "Goanna", p.Version)
	case strings.Contains(lower, "gecko/") &&
		(strings.Contains(lower, "firefox/") ||
			strings.Contains(lower, "fxios")):
		ver := ""
		if p, ok := state.Tokens.FindProduct("Firefox"); ok {
			ver = p.Version
		} else if p, ok := state.Tokens.FindProduct("FxiOS"); ok {
			ver = p.Version
		}
		state.SetEngine(engineClass(state), "Gecko", ver)
	case strings.Contains(lower, "presto/"):
		p, _ := state.Tokens.FindProduct("Presto")
		state.SetEngine(engineClass(state), "Presto", p.Version)
	case strings.Contains(lower, "applewebkit/"):
		// Blink vs WebKit: Chromium family uses AppleWebKit but is Blink.
		if isBlink(state, lower) {
			setBlink(state)
		} else {
			p, _ := state.Tokens.FindProduct("AppleWebKit")
			state.SetEngine(engineClass(state), "WebKit", p.Version)
		}
	case strings.Contains(lower, "khtml/"):
		p, _ := state.Tokens.FindProduct("KHTML")
		state.SetEngine(engineClass(state), "KHTML", p.Version)
	case strings.Contains(lower, "netfront/"):
		p, _ := state.Tokens.FindProduct("NetFront")
		state.SetEngine(engineClass(state), "NetFront", p.Version)
	default:
		if state.EngineClass == "" {
			state.EngineClass = "Unknown"
		}
	}
}

func isBlink(state *detect.State, lower string) bool {
	if strings.Contains(lower, "chrome/") ||
		strings.Contains(lower, "crios/") ||
		strings.Contains(lower, "chromium/") ||
		strings.Contains(lower, "edg/") ||
		strings.Contains(lower, "opr/") ||
		strings.Contains(lower, "samsungbrowser/") {
		return true
	}
	if state.AgentName != "" {
		switch strings.ToLower(state.AgentName) {
		case "chrome", "chromium", "edge", "opera", "opera touch",
			"samsung browser", "brave", "vivaldi", "yandex browser",
			"qq browser", "uc browser", "huaweibrowser", "heytapbrowser",
			"miuibrowser":
			return true
		}
	}
	return false
}

func setBlink(state *detect.State) {
	ver := ""
	if p, ok := state.Tokens.FindProduct("Chrome"); ok {
		ver = p.Version
	} else if p, ok := state.Tokens.FindProduct("CriOS"); ok {
		ver = p.Version
	} else if p, ok := state.Tokens.FindProduct("Chromium"); ok {
		ver = p.Version
	} else if p, ok := state.Tokens.FindProduct("Edg"); ok {
		ver = p.Version
	}
	state.SetEngine(engineClass(state), "Blink", ver)
}

func engineClass(state *detect.State) string {
	if state.EngineClass != "" {
		return state.EngineClass
	}
	switch state.AgentClass {
	case "Mobile App":
		return "Mobile App"
	case "Desktop App":
		return "Desktop App"
	case "Robot", "Robot Mobile":
		return "Robot"
	case "Hacker":
		return "Hacker"
	case "Server", "Cloud Application":
		return "Cloud"
	case "Voice", "Email Client", "Special", "Testclient":
		return "Special"
	default:
		return "Browser"
	}
}
