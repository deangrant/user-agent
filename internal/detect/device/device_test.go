package device

import (
	"testing"

	"github.com/deangrant/user-agent/internal/detect"
	"github.com/deangrant/user-agent/internal/tokenize"
)

func TestContainsToken(t *testing.T) {
	tests := []struct {
		s, key string
		want   bool
	}{
		{"draftphone", "aft", false},
		{"aftmm", "aft", true},
		{"; aft;", "aft", true},
		{"amazon aftmm build/x", "aft", true},
		{"something", "aft", false},
		{"kf", "kf", true},
		{"draftkf", "kf", false},
		{"apple watch", "watch", true},
		{"stopwatch", "watch", false},
	}
	for _, tt := range tests {
		got := containsToken(tt.s, tt.key)
		if got != tt.want {
			t.Fatalf("containsToken(%q, %q) = %v, want %v",
				tt.s, tt.key, got, tt.want)
		}
	}
}

func TestDraftPhoneNotFireTV(t *testing.T) {
	ua := "Mozilla/5.0 (Linux; Android 13; DraftPhone) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) " +
		"Chrome/120.0.0.0 Mobile Safari/537.36"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
		OSName: "Android",
	}
	New().Detect(state)
	if state.DeviceBrand == "Amazon" {
		t.Fatalf("brand = Amazon, want not Fire TV false positive; state=%#v",
			state)
	}
	if state.DeviceName == "Fire TV" {
		t.Fatalf("name = Fire TV, state=%#v", state)
	}
	if state.DeviceClass == "Set-top box" {
		t.Fatalf("class = Set-top box, state=%#v", state)
	}
}

func TestFireTVAFTMM(t *testing.T) {
	ua := "Mozilla/5.0 (Linux; Android 9; AFTMM Build/PS7233) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) " +
		"Silk/44.1.54 like Chrome/44.0.2403.63 Safari/537.36"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
		OSName: "Android",
	}
	New().Detect(state)
	if state.DeviceBrand != "Amazon" {
		t.Fatalf("brand = %q, want Amazon", state.DeviceBrand)
	}
	if state.DeviceName != "Fire TV" {
		t.Fatalf("name = %q, want Fire TV", state.DeviceName)
	}
	if state.DeviceClass != "Set-top box" {
		t.Fatalf("class = %q, want Set-top box", state.DeviceClass)
	}
}

func TestPixelWithoutMobileIsPhone(t *testing.T) {
	ua := "Mozilla/5.0 (Linux; Android 13; Pixel 7) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) " +
		"Chrome/120.0.0.0 Safari/537.36"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
		OSName: "Android",
	}
	New().Detect(state)
	if state.DeviceClass != "Phone" {
		t.Fatalf("class = %q, want Phone", state.DeviceClass)
	}
	if state.DeviceBrand != "Google" {
		t.Fatalf("brand = %q, want Google", state.DeviceBrand)
	}
	if state.DeviceName != "Pixel 7" {
		t.Fatalf("name = %q, want Pixel 7", state.DeviceName)
	}
}

func TestAndroidTabletTokenIsTablet(t *testing.T) {
	ua := "Mozilla/5.0 (Linux; Android 12; Lenovo TB-X606F) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) " +
		"Chrome/120.0.0.0 Safari/537.36"
	// Explicit tablet token (and no Mobile) should classify as Tablet.
	ua += " Tablet"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
		OSName: "Android",
	}
	New().Detect(state)
	if state.DeviceClass != "Tablet" {
		t.Fatalf("class = %q, want Tablet", state.DeviceClass)
	}
}

func TestAndroidMobileIsPhone(t *testing.T) {
	ua := "Mozilla/5.0 (Linux; Android 13; Pixel 7) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) " +
		"Chrome/120.0.0.0 Mobile Safari/537.36"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
		OSName: "Android",
	}
	New().Detect(state)
	if state.DeviceClass != "Phone" {
		t.Fatalf("class = %q, want Phone", state.DeviceClass)
	}
}

func TestGalaxyWatchSMRIsWatch(t *testing.T) {
	ua := "Mozilla/5.0 (Linux; Android 11; SAMSUNG SM-R860) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) " +
		"SamsungBrowser/1.0 Chrome/111.0.0.0 Mobile Safari/537.36"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
		OSName: "Android",
	}
	New().Detect(state)
	if state.DeviceClass != "Watch" {
		t.Fatalf("class = %q, want Watch", state.DeviceClass)
	}
	if state.DeviceBrand != "Samsung" {
		t.Fatalf("brand = %q, want Samsung", state.DeviceBrand)
	}
}

func TestAndroidTVIsTV(t *testing.T) {
	ua := "Mozilla/5.0 (Linux; Android 11; Android TV; " +
		"Build/RTM6.230109.121) AppleWebKit/537.36 " +
		"(KHTML, like Gecko) Chrome/94.0.4606.61 Safari/537.36"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
		OSName: "Android",
	}
	New().Detect(state)
	if state.DeviceClass != "TV" {
		t.Fatalf("class = %q, want TV", state.DeviceClass)
	}
}
