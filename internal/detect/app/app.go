// Package app detects mobile apps, webviews, and in-app browsers.
package app

import (
	"regexp"
	"strings"

	"github.com/deangrant/user-agent/internal/data"
	"github.com/deangrant/user-agent/internal/detect"
)

// entry is a compiled app pattern ready for matching.
type entry struct {
	pattern data.AppPattern
	needle  string // lowercased substring for Contains matches
	re      *regexp.Regexp
}

// Detector identifies applications and webviews.
type Detector struct {
	apps []entry
}

// New returns an app Detector using embedded patterns.
func New() *Detector {
	c := data.MustLoad()
	apps := make([]entry, 0, len(c.Apps))
	for i := range c.Apps {
		p := c.Apps[i]
		e := entry{
			pattern: p,
			needle:  strings.ToLower(p.Pattern),
		}
		if strings.Contains(p.Pattern, ".*") {
			e.re = regexp.MustCompile("(?i)" + p.Pattern)
		}
		apps = append(apps, e)
	}
	return &Detector{apps: apps}
}

// Detect implements detect.Detector.
func (d *Detector) Detect(state *detect.State) {
	if state == nil || state.UA == "" || state.BotMatched {
		return
	}
	ua := strings.ToLower(state.UA)

	for i := range d.apps {
		e := &d.apps[i]
		matched := false
		if e.re != nil {
			matched = e.re.MatchString(state.UA)
		} else if strings.Contains(ua, e.needle) {
			matched = true
		}
		if !matched {
			continue
		}
		// Skip patterns that are primarily browser identifiers;
		// those belong to the agent detector unless they set a
		// non-browser class.
		if isBrowserOnly(&e.pattern) {
			continue
		}
		applyApp(state, &e.pattern)
		return
	}

	// Generic Android WebView signal.
	if strings.Contains(ua, "; wv)") || strings.Contains(ua, "; wv ") {
		state.AppMatched = true
		state.SetAgent("Browser Webview", "Android WebView", "")
	}
}

func isBrowserOnly(p *data.AppPattern) bool {
	if p.AgentClass != "Browser" {
		return false
	}
	switch strings.ToLower(p.Name) {
	case "edge", "opera", "opera touch", "brave", "vivaldi",
		"yandex browser", "qq browser", "miuibrowser",
		"huaweibrowser", "heytapbrowser", "samsung browser",
		"uc browser", "duckduckgo", "ecosia", "mobile safari":
		return true
	default:
		return false
	}
}

func applyApp(state *detect.State, p *data.AppPattern) {
	state.AppMatched = true
	ver := ""
	if p.Product != "" {
		if prod, ok := state.Tokens.FindProduct(p.Product); ok {
			ver = prod.Version
		}
	}
	if ver == "" {
		for _, prod := range state.Tokens.Products {
			if strings.EqualFold(prod.Name, p.Name) {
				ver = prod.Version
				break
			}
		}
	}
	state.SetAgent(p.AgentClass, p.Name, ver)
	if p.DeviceClass != "" {
		state.DeviceClass = p.DeviceClass
	}
	if state.EngineClass == "" {
		switch p.AgentClass {
		case "Mobile App":
			state.EngineClass = "Mobile App"
		case "Desktop App":
			state.EngineClass = "Desktop App"
		case "Browser Webview":
			state.EngineClass = "Browser"
		case "Voice":
			state.EngineClass = "Special"
		case "Email Client":
			state.EngineClass = "Special"
		}
	}
}
