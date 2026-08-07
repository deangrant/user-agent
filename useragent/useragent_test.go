package useragent_test

import (
	"net/http"
	"testing"

	"github.com/deangrant/user-agent/useragent"
)

func TestParseChromeDesktop(t *testing.T) {
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 " +
		"(KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	r := useragent.Parse(ua)
	if r.Agent.Name != "Chrome" {
		t.Fatalf("agent = %q", r.Agent.Name)
	}
	if r.Agent.Class != useragent.AgentClassBrowser {
		t.Fatalf("agent class = %q", r.Agent.Class)
	}
	if r.LayoutEngine.Name != "Blink" {
		t.Fatalf("engine = %q", r.LayoutEngine.Name)
	}
	if r.OperatingSystem.Name != "Windows" ||
		r.OperatingSystem.Version != "10" {
		t.Fatalf("os = %s %s", r.OperatingSystem.Name,
			r.OperatingSystem.Version)
	}
	if r.Device.Class != useragent.DeviceClassDesktop {
		t.Fatalf("device = %q", r.Device.Class)
	}
	if !r.IsDesktop() || r.IsBot() {
		t.Fatalf("flags desktop=%v bot=%v", r.IsDesktop(), r.IsBot())
	}
	if r.AgentSecurity != useragent.AgentSecurityUnknown {
		t.Fatalf("security = %q, want Unknown", r.AgentSecurity)
	}
}

func TestParseAndroidChrome(t *testing.T) {
	ua := "Mozilla/5.0 (Linux; Android 7.0; Nexus 6 Build/NBD90Z) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) " +
		"Chrome/53.0.2785.124 Mobile Safari/537.36"
	r := useragent.Parse(ua)
	if r.Agent.Name != "Chrome" {
		t.Fatalf("agent = %q", r.Agent.Name)
	}
	if r.OperatingSystem.Name != "Android" ||
		r.OperatingSystem.Version != "7.0" {
		t.Fatalf("os = %#v", r.OperatingSystem)
	}
	if r.OperatingSystem.VersionBuild != "NBD90Z" {
		t.Fatalf("build = %q", r.OperatingSystem.VersionBuild)
	}
	if r.Device.Class != useragent.DeviceClassPhone {
		t.Fatalf("device class = %q", r.Device.Class)
	}
	if r.Device.Name != "Nexus 6" {
		t.Fatalf("device name = %q", r.Device.Name)
	}
	if !r.IsMobile() || !r.IsPhone() {
		t.Fatalf("mobile flags failed")
	}
}

func TestParseSafariMac(t *testing.T) {
	ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) " +
		"AppleWebKit/605.1.15 (KHTML, like Gecko) " +
		"Version/18.0 Safari/605.1.15"
	r := useragent.Parse(ua)
	if r.Agent.Name != "Safari" {
		t.Fatalf("agent = %q want Safari", r.Agent.Name)
	}
	if r.LayoutEngine.Name != "WebKit" {
		t.Fatalf("engine = %q", r.LayoutEngine.Name)
	}
	if r.OperatingSystem.Name != "Mac OS X" {
		t.Fatalf("os name = %q", r.OperatingSystem.Name)
	}
}

func TestParseFirefox(t *testing.T) {
	ua := "Mozilla/5.0 (X11; Linux x86_64; rv:121.0) " +
		"Gecko/20100101 Firefox/121.0"
	r := useragent.Parse(ua)
	if r.Agent.Name != "Firefox" || r.Agent.Version != "121.0" {
		t.Fatalf("agent = %#v", r.Agent)
	}
	if r.LayoutEngine.Name != "Gecko" {
		t.Fatalf("engine = %q", r.LayoutEngine.Name)
	}
}

func TestParseIPhone(t *testing.T) {
	ua := "Mozilla/5.0 (iPhone; CPU iPhone OS 17_2 like Mac OS X) " +
		"AppleWebKit/605.1.15 (KHTML, like Gecko) " +
		"Version/17.2 Mobile/15E148 Safari/604.1"
	r := useragent.Parse(ua)
	if r.OperatingSystem.Name != "iOS" {
		t.Fatalf("os = %q", r.OperatingSystem.Name)
	}
	if r.Device.Class != useragent.DeviceClassPhone {
		t.Fatalf("device = %q", r.Device.Class)
	}
	if r.Device.Brand != "Apple" {
		t.Fatalf("brand = %q", r.Device.Brand)
	}
}

func TestParseWithClientHintsWindows7(t *testing.T) {
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 " +
		"(KHTML, like Gecko) Chrome/100.0.0.0 Safari/537.36"
	hints := useragent.ClientHints{
		Brands: []useragent.BrandVersion{
			{Brand: "Not A;Brand", Version: "99"},
			{Brand: "Chromium", Version: "100"},
			{Brand: "Google Chrome", Version: "100"},
		},
		FullVersionList: []useragent.BrandVersion{
			{Brand: "Not A;Brand", Version: "99.0.0.0"},
			{Brand: "Chromium", Version: "100.0.4896.75"},
			{Brand: "Google Chrome", Version: "100.0.4896.75"},
		},
		Platform:        "Windows",
		PlatformVersion: "0.1.0",
		Mobile:          boolPtr(false),
		Arch:            "x86",
		Bitness:         "64",
	}
	r := useragent.ParseWithHints(ua, hints)
	if r.OperatingSystem.Name != "Windows" ||
		r.OperatingSystem.Version != "7" {
		t.Fatalf("os = %s %s, want Windows 7",
			r.OperatingSystem.Name, r.OperatingSystem.Version)
	}
	if r.Agent.Name != "Chrome" {
		t.Fatalf("agent = %q", r.Agent.Name)
	}
	if r.Agent.Version != "100.0.4896.75" {
		t.Fatalf("agent version = %q", r.Agent.Version)
	}
}

func TestParseWithClientHintsWindows11(t *testing.T) {
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 " +
		"(KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	mobile := false
	r := useragent.ParseWithHints(ua, useragent.ClientHints{
		Platform:        "Windows",
		PlatformVersion: "15.0.0",
		Mobile:          &mobile,
		FullVersionList: []useragent.BrandVersion{
			{Brand: "Google Chrome", Version: "120.0.6099.109"},
			{Brand: "Chromium", Version: "120.0.6099.109"},
		},
	})
	if r.OperatingSystem.Version != "11" {
		t.Fatalf("windows version = %q", r.OperatingSystem.Version)
	}
}

func TestParseWithHintsMobileOverridesTabletUA(t *testing.T) {
	ua := "Mozilla/5.0 (Linux; Android 12; Lenovo TB-X606F) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) " +
		"Chrome/120.0.0.0 Safari/537.36 Tablet"
	mobile := true
	r := useragent.ParseWithHints(ua, useragent.ClientHints{
		Platform: "Android",
		Mobile:   &mobile,
	})
	if r.Device.Class != useragent.DeviceClassPhone {
		t.Fatalf("class = %q, want Phone", r.Device.Class)
	}
}

func TestParseHeadersAndRequest(t *testing.T) {
	h := http.Header{}
	h.Set("User-Agent", "Mozilla/5.0 (Linux; Android 13; Pixel 7) "+
		"AppleWebKit/537.36 (KHTML, like Gecko) "+
		"Chrome/120.0.0.0 Mobile Safari/537.36")
	h.Set("Sec-CH-UA-Platform", `"Android"`)
	h.Set("Sec-CH-UA-Platform-Version", `"13.0.0"`)
	h.Set("Sec-CH-UA-Mobile", "?1")
	h.Set("Sec-CH-UA-Model", `"Pixel 7"`)
	h.Set("Sec-CH-UA-Full-Version-List",
		`"Chromium";v="120.0.6099.43", "Google Chrome";v="120.0.6099.43"`)

	r := useragent.ParseHeaders(h)
	if r.Device.Name != "Pixel 7" {
		t.Fatalf("model = %q", r.Device.Name)
	}
	if r.OperatingSystem.Name != "Android" {
		t.Fatalf("os = %q", r.OperatingSystem.Name)
	}
	if r.Agent.Version != "120.0.6099.43" {
		t.Fatalf("version = %q", r.Agent.Version)
	}

	req, err := http.NewRequestWithContext(
		t.Context(), http.MethodGet, "https://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header = h
	r2 := useragent.ParseRequest(req)
	if r2.Device.Name != r.Device.Name {
		t.Fatalf("ParseRequest mismatch: %#v vs %#v", r2.Device, r.Device)
	}
}

func TestParseGooglebot(t *testing.T) {
	ua := "Mozilla/5.0 (compatible; Googlebot/2.1; " +
		"+http://www.google.com/bot.html)"
	r := useragent.Parse(ua)
	if !r.IsBot() {
		t.Fatalf("expected bot, got %#v", r)
	}
	if r.Agent.Name != "Googlebot" {
		t.Fatalf("agent = %q", r.Agent.Name)
	}
	if r.Device.Class != useragent.DeviceClassRobot {
		t.Fatalf("device = %q", r.Device.Class)
	}
}

func TestParseCurlIsBot(t *testing.T) {
	r := useragent.Parse("curl/7.68.0")
	if !r.IsBot() {
		t.Fatalf("IsBot = false, want true; agent=%#v device=%#v",
			r.Agent, r.Device)
	}
	if r.Agent.Class != useragent.AgentClassServer {
		t.Fatalf("agent class = %q, want Server", r.Agent.Class)
	}
}

func TestParseOkHttpIsBot(t *testing.T) {
	r := useragent.Parse("okhttp/4.9.3")
	if !r.IsBot() {
		t.Fatalf("IsBot = false, want true; agent=%#v device=%#v",
			r.Agent, r.Device)
	}
	if r.Agent.Class != useragent.AgentClassServer {
		t.Fatalf("agent class = %q, want Server", r.Agent.Class)
	}
}

func TestParseSqlmapIsBot(t *testing.T) {
	r := useragent.Parse("sqlmap/1.4.2#stable (http://sqlmap.org)")
	if !r.IsBot() {
		t.Fatalf("IsBot = false, want true; agent=%#v device=%#v",
			r.Agent, r.Device)
	}
	if r.Agent.Class != useragent.AgentClassHacker {
		t.Fatalf("agent class = %q, want Hacker", r.Agent.Class)
	}
}

func TestParseEdge(t *testing.T) {
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 " +
		"(KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0"
	r := useragent.Parse(ua)
	if r.Agent.Name != "Edge" {
		t.Fatalf("agent = %q", r.Agent.Name)
	}
}

func TestParseInstagramApp(t *testing.T) {
	ua := "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) " +
		"AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 " +
		"Instagram 269.0.0.18.75"
	r := useragent.Parse(ua)
	if r.Agent.Name != "Instagram" {
		t.Fatalf("agent = %q", r.Agent.Name)
	}
	if r.Agent.Class != useragent.AgentClassMobileApp {
		t.Fatalf("class = %q", r.Agent.Class)
	}
	if r.Agent.Version != "269.0.0.18.75" {
		t.Fatalf("version = %q, want 269.0.0.18.75", r.Agent.Version)
	}
}

func TestParseWhatsAppAppNotBot(t *testing.T) {
	ua := "Mozilla/5.0 (Linux; Android 13; Pixel 7) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) " +
		"Chrome/120.0.0.0 Mobile Safari/537.36 WhatsApp/2.23.25.76"
	r := useragent.Parse(ua)
	if r.IsBot() {
		t.Fatalf("IsBot = true, want false; agent=%#v device=%#v",
			r.Agent, r.Device)
	}
	if r.Agent.Class != useragent.AgentClassMobileApp {
		t.Fatalf("class = %q, want Mobile App", r.Agent.Class)
	}
	if r.Agent.Name != "WhatsApp" {
		t.Fatalf("agent = %q, want WhatsApp", r.Agent.Name)
	}
}

func TestParseWhatsAppPreviewIsBot(t *testing.T) {
	ua := "WhatsApp/2.23.25.76 A"
	r := useragent.Parse(ua)
	if !r.IsBot() {
		t.Fatalf("IsBot = false, want true; result=%#v", r)
	}
}

func TestParsePinterestAppNotBot(t *testing.T) {
	ua := "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) " +
		"AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 " +
		"Pinterest/11.20"
	r := useragent.Parse(ua)
	if r.IsBot() {
		t.Fatalf("IsBot = true, want false; agent=%#v device=%#v",
			r.Agent, r.Device)
	}
	if r.Agent.Class != useragent.AgentClassMobileApp {
		t.Fatalf("class = %q, want Mobile App", r.Agent.Class)
	}
	if r.Agent.Name != "Pinterest" {
		t.Fatalf("agent = %q, want Pinterest", r.Agent.Name)
	}
}

func TestParsePinterestbotIsBot(t *testing.T) {
	ua := "Mozilla/5.0 (compatible; Pinterestbot/1.0; " +
		"+http://www.pinterest.com/bot.html)"
	r := useragent.Parse(ua)
	if !r.IsBot() {
		t.Fatalf("IsBot = false, want true; result=%#v", r)
	}
	if r.Agent.Name != "Pinterestbot" {
		t.Fatalf("agent = %q, want Pinterestbot", r.Agent.Name)
	}
}

func TestClientHintsFromMap(t *testing.T) {
	h := useragent.ClientHintsFromMap(map[string]string{
		"Sec-CH-UA-Platform": `"Linux"`,
		"Sec-CH-UA-Mobile":   "?0",
	})
	if h.Platform != "Linux" {
		t.Fatalf("platform = %q", h.Platform)
	}
	if h.Mobile == nil || *h.Mobile {
		t.Fatalf("mobile = %v", h.Mobile)
	}
}

func TestParseEmpty(t *testing.T) {
	r := useragent.Parse("")
	if r.Device.Class != useragent.DeviceClassUnknown {
		t.Fatalf("device = %q", r.Device.Class)
	}
}

func TestParseNestedParensNoGarbageDeviceName(t *testing.T) {
	ua := "Mozilla/5.0 ((Linux; Android 13; Pixel 7) AppleWebKit/537.36 " +
		"(KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36"
	r := useragent.Parse(ua)
	if r.Device.Name == "(Linux" {
		t.Fatalf("device name = %q", r.Device.Name)
	}
	if r.Agent.Name != "Chrome" {
		t.Fatalf("agent = %q, want Chrome", r.Agent.Name)
	}
}

func TestParseLegacyFirefoxStrongSecurity(t *testing.T) {
	ua := "Mozilla/5.0 (X11; U; Linux i686; en-US; rv:1.9.0.4) " +
		"Gecko/20100101 Firefox/3.0.4"
	r := useragent.Parse(ua)
	if r.AgentSecurity != useragent.AgentSecurityStrong {
		t.Fatalf("security = %q, want Strong security", r.AgentSecurity)
	}
	if r.Agent.Name != "Firefox" {
		t.Fatalf("agent = %q, want Firefox", r.Agent.Name)
	}
}

func boolPtr(v bool) *bool { return &v }
