// Package agent detects browsers and related user agents.
package agent

import (
	"strings"

	"github.com/deangrant/user-agent/internal/detect"
	"github.com/deangrant/user-agent/internal/uaclass"
)

// Detector identifies the browsing agent.
type Detector struct{}

// New returns an agent Detector.
func New() *Detector { return &Detector{} }

type rule struct {
	product string
	name    string
	class   string
}

// Order matters: more specific agents before generic Chrome/Safari.
var rules = []rule{
	{product: "Edg", name: "Edge", class: uaclass.Browser},
	{product: "EdgA", name: "Edge", class: uaclass.Browser},
	{product: "EdgiOS", name: "Edge", class: uaclass.Browser},
	{product: "Edge", name: "Edge", class: uaclass.Browser},
	{product: "OPR", name: "Opera", class: uaclass.Browser},
	{product: "OPT", name: "Opera Touch", class: uaclass.Browser},
	{product: "Opera", name: "Opera", class: uaclass.Browser},
	{
		product: "SamsungBrowser",
		name:    "Samsung Browser",
		class:   uaclass.Browser,
	},
	{product: "UCBrowser", name: "UC Browser", class: uaclass.Browser},
	{product: "YaBrowser", name: "Yandex Browser", class: uaclass.Browser},
	{product: "Vivaldi", name: "Vivaldi", class: uaclass.Browser},
	{product: "Brave", name: "Brave", class: uaclass.Browser},
	{product: "QQBrowser", name: "QQ Browser", class: uaclass.Browser},
	{product: "MiuiBrowser", name: "MiuiBrowser", class: uaclass.Browser},
	{product: "HuaweiBrowser", name: "HuaweiBrowser", class: uaclass.Browser},
	{product: "HeyTapBrowser", name: "HeyTapBrowser", class: uaclass.Browser},
	{product: "DuckDuckGo", name: "DuckDuckGo", class: uaclass.Browser},
	{product: "CriOS", name: "Chrome", class: uaclass.Browser},
	{product: "FxiOS", name: "Firefox", class: uaclass.Browser},
	{product: "Firefox", name: "Firefox", class: uaclass.Browser},
	{product: "Fennec", name: "Firefox", class: uaclass.Browser},
	{product: "Chrome", name: "Chrome", class: uaclass.Browser},
	{product: "Chromium", name: "Chromium", class: uaclass.Browser},
	{product: "Version", name: "Safari", class: uaclass.Browser},
	{product: "MSIE", name: "IE", class: uaclass.Browser},
	{product: "Trident", name: "IE", class: uaclass.Browser},
	{product: "Silk", name: "Silk", class: uaclass.Browser},
	{product: "Coast", name: "Opera Coast", class: uaclass.Browser},
	{product: "Puffin", name: "Puffin", class: uaclass.Browser},
	{product: "Sraf", name: "Sraf Browser", class: uaclass.Browser},
	{product: "Hisense", name: "Hisense HiBrowser", class: uaclass.Browser},
}

// Detect implements detect.Detector.
func (d *Detector) Detect(state *detect.State) {
	if state == nil || state.UA == "" {
		return
	}
	if state.BotMatched {
		return
	}
	defer detectSecurity(state)

	if state.AppMatched && state.AgentName != "" &&
		state.AgentClass != uaclass.Browser &&
		state.AgentClass != uaclass.BrowserWebview {
		return
	}
	// Webview: still try to detect underlying browser engine agent.
	if state.AppMatched && state.AgentClass == uaclass.BrowserWebview {
		detectBrowser(state, true)
		return
	}
	if state.AgentName != "" && state.AgentClass != "" &&
		state.AgentClass != uaclass.Unknown {
		return
	}
	detectBrowser(state, false)
}

// detectSecurity sets AgentSecurity from classic Mozilla comment
// encryption tokens (N/I/U). Does not overwrite a prior value.
func detectSecurity(state *detect.State) {
	if state.AgentSecurity != "" {
		return
	}
	for _, c := range state.Tokens.Comments {
		switch strings.ToUpper(strings.TrimSpace(c)) {
		case "N":
			state.AgentSecurity = uaclass.SecurityNone
			return
		case "I":
			state.AgentSecurity = uaclass.SecurityWeak
			return
		case "U":
			state.AgentSecurity = uaclass.SecurityStrong
			return
		}
	}
}

func detectBrowser(state *detect.State, webview bool) {
	ua := state.UA
	lower := strings.ToLower(ua)

	// IE 11: Trident/7.0 with rv:11.0
	if strings.Contains(lower, "trident/") && strings.Contains(lower, "rv:") {
		ver := extractRV(ua)
		class := uaclass.Browser
		if webview {
			class = uaclass.BrowserWebview
		}
		state.SetAgent(class, "IE", ver)
		return
	}

	for _, r := range rules {
		p, ok := state.Tokens.FindProduct(r.product)
		if !ok {
			continue
		}
		if r.product == "Version" {
			if !hasSafari(state) {
				continue
			}
			// Prefer Safari when Version + Safari present.
			class := uaclass.Browser
			if webview {
				class = uaclass.BrowserWebview
			}
			state.SetAgent(class, "Safari", p.Version)
			return
		}
		if r.product == "Trident" {
			continue // handled above / via MSIE
		}
		name := r.name
		ver := p.Version
		class := r.class
		if webview && class == uaclass.Browser {
			// Keep discovered browser name but mark webview class
			// only when we didn't already set a more specific app.
			if state.AgentClass != uaclass.BrowserWebview {
				class = uaclass.Browser
			} else {
				// Fill name/version under existing webview class.
				state.SetAgent(state.AgentClass, name, ver)
				return
			}
		}
		state.SetAgent(class, name, ver)
		return
	}

	// Fallback: last non-Mozilla product.
	for i := len(state.Tokens.Products) - 1; i >= 0; i-- {
		p := state.Tokens.Products[i]
		if strings.EqualFold(p.Name, "Mozilla") ||
			strings.EqualFold(p.Name, "AppleWebKit") ||
			strings.EqualFold(p.Name, "Safari") ||
			strings.EqualFold(p.Name, "Mobile") ||
			strings.EqualFold(p.Name, "Gecko") {
			continue
		}
		class := uaclass.Browser
		if webview {
			class = uaclass.BrowserWebview
		}
		state.SetAgent(class, p.Name, p.Version)
		return
	}
}

func hasSafari(state *detect.State) bool {
	_, ok := state.Tokens.FindProduct("Safari")
	return ok
}

func extractRV(ua string) string {
	lower := strings.ToLower(ua)
	i := strings.Index(lower, "rv:")
	if i < 0 {
		return ""
	}
	i += 3
	j := i
	for j < len(ua) {
		c := ua[j]
		if (c >= '0' && c <= '9') || c == '.' {
			j++
			continue
		}
		break
	}
	return ua[i:j]
}
