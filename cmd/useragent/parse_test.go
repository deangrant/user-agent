package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/deangrant/user-agent/useragent"
)

func TestHeadersFromReader(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantUA   string
		wantPlat string
		wantErr  bool
	}{
		{
			name: "with blank line",
			input: "User-Agent: Mozilla/5.0\r\n" +
				"Sec-CH-UA-Platform: \"Android\"\r\n\r\n",
			wantUA:   "Mozilla/5.0",
			wantPlat: "\"Android\"",
		},
		{
			name: "without trailing blank line",
			input: "User-Agent: Mozilla/5.0\n" +
				"Sec-CH-UA-Platform: \"Linux\"\n",
			wantUA:   "Mozilla/5.0",
			wantPlat: "\"Linux\"",
		},
		{
			name:    "empty",
			input:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := headersFromReader(strings.NewReader(tt.input))
			if tt.wantErr {
				if err == nil {
					t.Fatal("error = nil, want error")
				}
				return
			}
			if err != nil {
				t.Fatalf("error = %v", err)
			}
			if got := h.Get("User-Agent"); got != tt.wantUA {
				t.Fatalf("User-Agent = %q, want %q", got, tt.wantUA)
			}
			if got := h.Get("Sec-CH-UA-Platform"); got != tt.wantPlat {
				t.Fatalf("Platform = %q, want %q", got, tt.wantPlat)
			}
		})
	}
}

func TestNormalizeHeaderBlock(t *testing.T) {
	got := normalizeHeaderBlock("User-Agent: x")
	if !strings.HasSuffix(got, "\n\n") {
		t.Fatalf("got %q, want trailing blank line", got)
	}
}

func TestWriteResultLines(t *testing.T) {
	r := useragent.Result{
		UserAgent: "ua",
		Device: useragent.Device{
			Class: useragent.DeviceClassDesktop,
		},
		Agent: useragent.Agent{
			Class: useragent.AgentClassBrowser,
			Name:  "Chrome",
		},
		AgentSecurity: useragent.AgentSecurityUnknown,
	}
	var buf bytes.Buffer
	if err := writeResult(&buf, r, false, false); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	for _, want := range []string{
		"[+] UserAgent: ua\n",
		"[+] Device.Class: Desktop\n",
		"[+] Agent.Name: Chrome\n",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q in:\n%s", want, out)
		}
	}
	for _, skip := range []string{
		"AgentSecurity:",
		"Device.Name:",
		"Unknown",
	} {
		if strings.Contains(out, skip) {
			t.Fatalf("unexpected %q in:\n%s", skip, out)
		}
	}
	if strings.Contains(out, "{") {
		t.Fatalf("unexpected JSON: %s", out)
	}
	if strings.Contains(out, "\033[") {
		t.Fatalf("unexpected ANSI color in buffer output: %q", out)
	}
}

func TestWriteResultJSON(t *testing.T) {
	r := useragent.Result{
		UserAgent: "ua",
		Agent: useragent.Agent{
			Class: useragent.AgentClassBrowser,
			Name:  "Chrome",
		},
	}
	var pretty, compact bytes.Buffer
	if err := writeResult(&pretty, r, true, false); err != nil {
		t.Fatal(err)
	}
	if err := writeResult(&compact, r, false, true); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(pretty.String(), "\n") {
		t.Fatalf("pretty missing newlines: %q", pretty.String())
	}
	if strings.Count(compact.String(), "\n") != 1 {
		t.Fatalf("compact = %q", compact.String())
	}
	var decoded useragent.Result
	if err := json.Unmarshal(pretty.Bytes(), &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Agent.Name != "Chrome" {
		t.Fatalf("decoded agent = %#v", decoded.Agent)
	}
}

func TestRunParseUA(t *testing.T) {
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) " +
		"Chrome/120.0.0.0 Safari/537.36"
	var stdout, stderr bytes.Buffer
	err := run([]string{"parse", ua},
		strings.NewReader(""), &stdout, &stderr)
	if err != nil {
		t.Fatalf("run: %v (stderr=%s)", err, stderr.String())
	}
	out := stdout.String()
	if !strings.Contains(out, "[+] Agent.Name: Chrome\n") {
		t.Fatalf("missing Agent.Name in:\n%s", out)
	}
	if !strings.Contains(out, "[+] OperatingSystem.Name: Windows\n") {
		t.Fatalf("missing OS.Name in:\n%s", out)
	}
	if !strings.Contains(out, "[+] Device.Class: Desktop\n") {
		t.Fatalf("missing Device.Class in:\n%s", out)
	}
}

func TestRunParseStdinHeaders(t *testing.T) {
	input := "User-Agent: Mozilla/5.0 (Linux; Android 13; Pixel 7) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) " +
		"Chrome/120.0.0.0 Mobile Safari/537.36\n" +
		"Sec-CH-UA-Platform: \"Android\"\n" +
		"Sec-CH-UA-Platform-Version: \"13.0.0\"\n" +
		"Sec-CH-UA-Model: \"Pixel 7\"\n" +
		"Sec-CH-UA-Mobile: ?1\n"
	var stdout, stderr bytes.Buffer
	err := run([]string{"parse", "-"},
		strings.NewReader(input), &stdout, &stderr)
	if err != nil {
		t.Fatalf("run: %v (stderr=%s)", err, stderr.String())
	}
	out := stdout.String()
	if !strings.Contains(out, "[+] Device.Name: Pixel 7\n") {
		t.Fatalf("missing Device.Name in:\n%s", out)
	}
	if !strings.Contains(out, "[+] OperatingSystem.Name: Android\n") {
		t.Fatalf("missing OS.Name in:\n%s", out)
	}
}

func TestRunParseWithHintFlags(t *testing.T) {
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) " +
		"Chrome/100.0.0.0 Safari/537.36"
	args := []string{
		"parse",
		"-sec-ch-ua-platform", `"Windows"`,
		"-sec-ch-ua-platform-version", `"0.1.0"`,
		"-sec-ch-ua-full-version-list",
		`"Google Chrome";v="100.0.4896.75", "Chromium";v="100.0.4896.75"`,
		ua,
	}
	var stdout, stderr bytes.Buffer
	err := run(args, strings.NewReader(""), &stdout, &stderr)
	if err != nil {
		t.Fatalf("run: %v (stderr=%s)", err, stderr.String())
	}
	out := stdout.String()
	if !strings.Contains(out, "[+] OperatingSystem.Version: 7\n") {
		t.Fatalf("want Windows 7 in:\n%s", out)
	}
	if !strings.Contains(out, "[+] Agent.Version: 100.0.4896.75\n") {
		t.Fatalf("want agent version in:\n%s", out)
	}
}

func TestRunParseJSONFlag(t *testing.T) {
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) " +
		"Chrome/120.0.0.0 Safari/537.36"
	var stdout, stderr bytes.Buffer
	err := run([]string{"parse", "-compact", ua},
		strings.NewReader(""), &stdout, &stderr)
	if err != nil {
		t.Fatalf("run: %v (stderr=%s)", err, stderr.String())
	}
	var got useragent.Result
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("json: %v body=%s", err, stdout.String())
	}
	if got.Agent.Name != "Chrome" {
		t.Fatalf("agent = %q", got.Agent.Name)
	}
}

func TestRunMissingSubcommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := run(nil, strings.NewReader(""), &stdout, &stderr)
	if err == nil {
		t.Fatal("want usage error")
	}
	var ue *usageError
	if !errors.As(err, &ue) {
		t.Fatalf("err = %T %v", err, err)
	}
}
