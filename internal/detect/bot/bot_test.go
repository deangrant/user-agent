package bot

import (
	"testing"

	"github.com/deangrant/user-agent/internal/detect"
	"github.com/deangrant/user-agent/internal/tokenize"
)

func TestDetectGooglebot(t *testing.T) {
	ua := "Mozilla/5.0 (compatible; Googlebot/2.1; " +
		"+http://www.google.com/bot.html)"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if !state.BotMatched {
		t.Fatal("BotMatched = false, want true")
	}
	if state.AgentName != "Googlebot" {
		t.Fatalf("agent = %q, want Googlebot", state.AgentName)
	}
	if state.AgentClass != "Robot" {
		t.Fatalf("agent class = %q, want Robot", state.AgentClass)
	}
	if state.DeviceClass != "Robot" {
		t.Fatalf("device class = %q, want Robot", state.DeviceClass)
	}
	if state.EngineClass != "Robot" {
		t.Fatalf("engine class = %q, want Robot", state.EngineClass)
	}
	if state.OSClass != "Cloud" {
		t.Fatalf("os class = %q, want Cloud", state.OSClass)
	}
}

func TestDetectBrowserWithBotSubstringNotBot(t *testing.T) {
	// "robot" contains heuristic "bot", but a normal Chrome UA must not match.
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 " +
		"(KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 RobotLab/1.0"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if state.BotMatched {
		t.Fatalf("BotMatched = true, want false; agent=%#v", state)
	}
}

func TestDetectHeuristicGenericBot(t *testing.T) {
	ua := "SomethingBot/1.0"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if !state.BotMatched {
		t.Fatal("BotMatched = false, want true")
	}
	if state.AgentName != "Generic Bot" {
		t.Fatalf("agent = %q, want Generic Bot", state.AgentName)
	}
}

func TestDetectWhatsAppPreviewIsBot(t *testing.T) {
	ua := "WhatsApp/2.23.25.76 A"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if !state.BotMatched {
		t.Fatal("BotMatched = false, want true")
	}
	if state.AgentName != "WhatsApp" {
		t.Fatalf("agent = %q, want WhatsApp", state.AgentName)
	}
}

func TestDetectChromeWhatsAppNotHeuristicBot(t *testing.T) {
	ua := "Mozilla/5.0 (Linux; Android 13; Pixel 7) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) " +
		"Chrome/120.0.0.0 Mobile Safari/537.36 WhatsApp/2.23.25.76"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if state.BotMatched {
		t.Fatalf("BotMatched = true, want false; agent=%#v", state)
	}
}

func TestDetectEmptyUA(t *testing.T) {
	state := &detect.State{}
	New().Detect(state)
	if state.BotMatched {
		t.Fatal("BotMatched = true on empty UA")
	}
	New().Detect(nil) // must not panic
}

func TestDetectHackerSetsSecurity(t *testing.T) {
	ua := "sqlmap/1.4.2#stable (http://sqlmap.org)"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if !state.BotMatched {
		t.Fatal("BotMatched = false, want true")
	}
	if state.AgentClass != "Hacker" {
		t.Fatalf("agent class = %q, want Hacker", state.AgentClass)
	}
	if state.AgentSecurity != "Hacker" {
		t.Fatalf("security = %q, want Hacker", state.AgentSecurity)
	}
}

func TestDetectGooglebotLeavesSecurityEmpty(t *testing.T) {
	ua := "Mozilla/5.0 (compatible; Googlebot/2.1; " +
		"+http://www.google.com/bot.html)"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if !state.BotMatched {
		t.Fatal("BotMatched = false, want true")
	}
	if state.AgentSecurity != "" {
		t.Fatalf("security = %q, want empty for Robot", state.AgentSecurity)
	}
}

func TestDetectSogouBrowserNotBot(t *testing.T) {
	ua := "Mozilla/5.0 (Linux; Android 13; Pixel 7) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) " +
		"Version/4.0 Chrome/120.0.0.0 Mobile Safari/537.36 " +
		"SogouMobileBrowser/6.0"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if state.BotMatched {
		t.Fatalf("BotMatched = true, want false; agent=%#v", state)
	}
}

func TestDetectFlipboardAppNotBot(t *testing.T) {
	ua := "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) " +
		"AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 " +
		"Flipboard/4.2.0"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if state.BotMatched {
		t.Fatalf("BotMatched = true, want false; agent=%#v", state)
	}
}

func TestDetectFlipboardProxyIsBot(t *testing.T) {
	ua := "Mozilla/5.0 (compatible; FlipboardProxy/1.0; " +
		"+http://flipboard.com/browserproxy)"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if !state.BotMatched {
		t.Fatal("BotMatched = false, want true")
	}
	if state.AgentName != "FlipboardProxy" {
		t.Fatalf("agent = %q, want FlipboardProxy", state.AgentName)
	}
}

func TestDetectWordPressPingbackIsBot(t *testing.T) {
	ua := "WordPress/6.4; https://example.com"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if !state.BotMatched {
		t.Fatal("BotMatched = false, want true")
	}
	if state.AgentName != "WordPress" {
		t.Fatalf("agent = %q, want WordPress", state.AgentName)
	}
}

func TestDetectWordPressAppNotBot(t *testing.T) {
	ua := "wp-android/22.0 (wordpress for android)"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if state.BotMatched {
		t.Fatalf("BotMatched = true, want false; agent=%#v", state)
	}
}

func TestDetectUptimeRobotIsBot(t *testing.T) {
	ua := "Mozilla/5.0+(compatible; UptimeRobot/2.0; " +
		"http://www.uptimerobot.com/)"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if !state.BotMatched {
		t.Fatal("BotMatched = false, want true")
	}
	if state.AgentName != "Uptime Monitor" {
		t.Fatalf("agent = %q, want Uptime Monitor", state.AgentName)
	}
}

func TestDetectUnrelatedUptimeNotBot(t *testing.T) {
	ua := "MyUptimeDashboard/1.0"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if state.BotMatched {
		t.Fatalf("BotMatched = true, want false; agent=%#v", state)
	}
}

func TestDetectBrowserMastodonNotBot(t *testing.T) {
	ua := "Mozilla/5.0 (Linux; Android 13; Pixel 7) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) " +
		"Chrome/120.0.0.0 Mobile Safari/537.36 Mastodon/4.0"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if state.BotMatched {
		t.Fatalf("BotMatched = true, want false; agent=%#v", state)
	}
}

func TestDetectBareMastodonIsBot(t *testing.T) {
	ua := "http.rb/5.1.0 (Mastodon/4.0; +https://example.social/)"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if !state.BotMatched {
		t.Fatal("BotMatched = false, want true")
	}
	if state.AgentName != "Mastodon" {
		t.Fatalf("agent = %q, want Mastodon", state.AgentName)
	}
}

func TestDetectBrowserTumblrNotBot(t *testing.T) {
	ua := "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) " +
		"AppleWebKit/605.1.15 (KHTML, like Gecko) " +
		"Version/16.0 Mobile/15E148 Safari/604.1 Tumblr/20.0"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if state.BotMatched {
		t.Fatalf("BotMatched = true, want false; agent=%#v", state)
	}
}
