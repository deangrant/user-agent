package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/textproto"
	"os"
	"strings"

	"github.com/deangrant/user-agent/useragent"
)

// parseFlags holds CLI options for the parse subcommand.
type parseFlags struct {
	json    bool
	compact bool

	secCHUA                string
	secCHUAFullVersionList string
	secCHUAPlatform        string
	secCHUAPlatformVersion string
	secCHUAMobile          string
	secCHUAModel           string
	secCHUAArch            string
	secCHUABitness         string
	secCHUAFormFactors     string
}

func runParse(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("parse", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var f parseFlags
	fs.BoolVar(&f.json, "json", false, "emit indented JSON")
	fs.BoolVar(&f.compact, "compact", false, "emit compact JSON")
	fs.StringVar(&f.secCHUA, "sec-ch-ua", "", "Sec-CH-UA value")
	fs.StringVar(&f.secCHUAFullVersionList, "sec-ch-ua-full-version-list", "",
		"Sec-CH-UA-Full-Version-List value")
	fs.StringVar(&f.secCHUAPlatform, "sec-ch-ua-platform", "",
		"Sec-CH-UA-Platform value")
	fs.StringVar(&f.secCHUAPlatformVersion, "sec-ch-ua-platform-version", "",
		"Sec-CH-UA-Platform-Version value")
	fs.StringVar(&f.secCHUAMobile, "sec-ch-ua-mobile", "",
		"Sec-CH-UA-Mobile value")
	fs.StringVar(&f.secCHUAModel, "sec-ch-ua-model", "",
		"Sec-CH-UA-Model value")
	fs.StringVar(&f.secCHUAArch, "sec-ch-ua-arch", "",
		"Sec-CH-UA-Arch value")
	fs.StringVar(&f.secCHUABitness, "sec-ch-ua-bitness", "",
		"Sec-CH-UA-Bitness value")
	fs.StringVar(&f.secCHUAFormFactors, "sec-ch-ua-form-factors", "",
		"Sec-CH-UA-Form-Factors value")

	fs.Usage = func() {
		fmt.Fprint(stderr, `Usage:
  useragent parse [flags] <user-agent>
  useragent parse [flags] -

`)
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return err
	}
	pos := fs.Args()
	if len(pos) == 0 {
		fs.Usage()
		return &usageError{msg: "missing user-agent or '-'"}
	}
	if len(pos) > 1 {
		fs.Usage()
		return &usageError{msg: "too many arguments"}
	}

	var result useragent.Result
	switch pos[0] {
	case "-":
		h, err := headersFromReader(stdin)
		if err != nil {
			return err
		}
		applyHintFlags(h, f)
		result = useragent.ParseHeaders(h)
	default:
		h := http.Header{}
		h.Set("User-Agent", pos[0])
		applyHintFlags(h, f)
		result = useragent.ParseHeaders(h)
	}
	return writeResult(stdout, result, f.json, f.compact)
}

func applyHintFlags(h http.Header, f parseFlags) {
	set := func(name, value string) {
		if value != "" {
			h.Set(name, value)
		}
	}
	set("Sec-CH-UA", f.secCHUA)
	set("Sec-CH-UA-Full-Version-List", f.secCHUAFullVersionList)
	set("Sec-CH-UA-Platform", f.secCHUAPlatform)
	set("Sec-CH-UA-Platform-Version", f.secCHUAPlatformVersion)
	set("Sec-CH-UA-Mobile", f.secCHUAMobile)
	set("Sec-CH-UA-Model", f.secCHUAModel)
	set("Sec-CH-UA-Arch", f.secCHUAArch)
	set("Sec-CH-UA-Bitness", f.secCHUABitness)
	set("Sec-CH-UA-Form-Factors", f.secCHUAFormFactors)
}

// headersFromReader parses MIME-style HTTP headers from r until a blank
// line or EOF.
func headersFromReader(r io.Reader) (http.Header, error) {
	raw, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read headers: %w", err)
	}
	br := bufio.NewReader(strings.NewReader(normalizeHeaderBlock(string(raw))))
	tp := textproto.NewReader(br)
	mh, err := tp.ReadMIMEHeader()
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("read headers: %w", err)
	}
	h := http.Header(mh)
	if len(h) == 0 {
		return nil, &usageError{msg: "stdin contained no headers"}
	}
	return h, nil
}

func writeResult(
	w io.Writer,
	result useragent.Result,
	asJSON, compact bool,
) error {
	if asJSON || compact {
		return writeResultJSON(w, result, compact)
	}
	return writeResultLines(w, result)
}

func writeResultJSON(
	w io.Writer,
	result useragent.Result,
	compact bool,
) error {
	enc := json.NewEncoder(w)
	if !compact {
		enc.SetIndent("", "  ")
	}
	if err := enc.Encode(result); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func writeResultLines(w io.Writer, result useragent.Result) error {
	lines := []struct {
		path  string
		value string
	}{
		{"UserAgent", result.UserAgent},
		{"Device.Class", string(result.Device.Class)},
		{"Device.Name", result.Device.Name},
		{"Device.Brand", result.Device.Brand},
		{"Device.CPU", result.Device.CPU},
		{"OperatingSystem.Class", string(result.OperatingSystem.Class)},
		{"OperatingSystem.Name", result.OperatingSystem.Name},
		{"OperatingSystem.Version", result.OperatingSystem.Version},
		{"OperatingSystem.VersionBuild", result.OperatingSystem.VersionBuild},
		{"OperatingSystem.NameVersion", result.OperatingSystem.NameVersion},
		{"LayoutEngine.Class", string(result.LayoutEngine.Class)},
		{"LayoutEngine.Name", result.LayoutEngine.Name},
		{"LayoutEngine.Version", result.LayoutEngine.Version},
		{"LayoutEngine.VersionMajor", result.LayoutEngine.VersionMajor},
		{"LayoutEngine.NameVersion", result.LayoutEngine.NameVersion},
		{"LayoutEngine.NameVersionMajor", result.LayoutEngine.NameVersionMajor},
		{"Agent.Class", string(result.Agent.Class)},
		{"Agent.Name", result.Agent.Name},
		{"Agent.Version", result.Agent.Version},
		{"Agent.VersionMajor", result.Agent.VersionMajor},
		{"Agent.NameVersion", result.Agent.NameVersion},
		{"Agent.NameVersionMajor", result.Agent.NameVersionMajor},
		{"AgentSecurity", string(result.AgentSecurity)},
	}
	color := useColor(w)
	for _, line := range lines {
		if skipFieldValue(line.value) {
			continue
		}
		if err := writeFieldLine(w, line.path, line.value, color); err != nil {
			return err
		}
	}
	if result.ClientHintsMismatch {
		err := writeFieldLine(w, "ClientHintsMismatch", "true", color)
		if err != nil {
			return err
		}
	}
	return nil
}

func skipFieldValue(value string) bool {
	v := strings.TrimSpace(value)
	return v == "" || strings.EqualFold(v, "Unknown")
}

const (
	ansiGreen = "\033[32m"
	ansiReset = "\033[0m"
)

func writeFieldLine(w io.Writer, path, value string, color bool) error {
	if color {
		_, err := fmt.Fprintf(w, "%s[+]%s %s%s%s: %s\n",
			ansiGreen, ansiReset,
			ansiGreen, path, ansiReset,
			value)
		return err
	}
	_, err := fmt.Fprintf(w, "[+] %s: %s\n", path, value)
	return err
}

// useColor reports whether ANSI colors should be written to w.
func useColor(w io.Writer) bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

// normalizeHeaderBlock ensures a trailing blank line so ReadMIMEHeader
// succeeds when callers pass header text without a final empty line.
func normalizeHeaderBlock(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.TrimRight(s, "\n")
	if s == "" {
		return "\n\n"
	}
	return s + "\n\n"
}
