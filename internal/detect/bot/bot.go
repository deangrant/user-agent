// Package bot detects robots, crawlers, hackers, and HTTP clients.
package bot

import (
	"strings"

	"github.com/deangrant/user-agent/internal/data"
	"github.com/deangrant/user-agent/internal/detect"
)

// Detector identifies bots and similar non-browser clients.
type Detector struct {
	bots []data.BotPattern
}

// New returns a bot Detector using embedded patterns.
func New() *Detector {
	c := data.MustLoad()
	return &Detector{bots: c.Bots}
}

// Detect implements detect.Detector.
func (d *Detector) Detect(state *detect.State) {
	if state == nil || state.UA == "" {
		return
	}
	ua := strings.ToLower(state.UA)

	var heuristic *data.BotPattern
	for i := range d.bots {
		p := &d.bots[i]
		if !strings.Contains(ua, strings.ToLower(p.Pattern)) {
			continue
		}
		if p.Heuristic {
			if heuristic == nil {
				heuristic = p
			}
			continue
		}
		applyBot(state, p)
		return
	}
	if heuristic != nil && looksLikeBot(ua) {
		applyBot(state, heuristic)
	}
}

func applyBot(state *detect.State, p *data.BotPattern) {
	state.BotMatched = true
	state.SetAgent(p.AgentClass, p.Name, productVersion(state, p.Name))
	if p.DeviceClass != "" {
		state.DeviceClass = p.DeviceClass
	}
	if state.EngineClass == "" {
		switch p.AgentClass {
		case "Robot", "Testclient", "Hacker", "Server", "Cloud Application":
			state.EngineClass = mapEngineClass(p.AgentClass)
		}
	}
	if state.OSClass == "" {
		switch p.DeviceClass {
		case "Hacker":
			state.OSClass = "Hacker"
		case "Cloud", "Robot", "Robot Mobile", "Robot Imitator":
			state.OSClass = "Cloud"
		}
	}
}

func mapEngineClass(agentClass string) string {
	switch agentClass {
	case "Robot", "Testclient":
		return "Robot"
	case "Hacker":
		return "Hacker"
	case "Server", "Cloud Application":
		return "Cloud"
	default:
		return "Unknown"
	}
}

func productVersion(state *detect.State, name string) string {
	for _, p := range state.Tokens.Products {
		if strings.EqualFold(p.Name, name) ||
			strings.Contains(strings.ToLower(p.Name), strings.ToLower(name)) {
			return p.Version
		}
	}
	// Try first product after Mozilla.
	for _, p := range state.Tokens.Products {
		if !strings.EqualFold(p.Name, "Mozilla") && p.Version != "" {
			return p.Version
		}
	}
	return ""
}

func looksLikeBot(ua string) bool {
	// Avoid classifying normal browsers that happen to include "bot"
	// substrings in unrelated tokens.
	if strings.Contains(ua, "mozilla/") &&
		(strings.Contains(ua, "chrome/") ||
			strings.Contains(ua, "safari/") ||
			strings.Contains(ua, "firefox/")) &&
		!strings.Contains(ua, "googlebot") &&
		!strings.Contains(ua, "bingbot") &&
		!strings.Contains(ua, "adsbot") {
		return false
	}
	return true
}
