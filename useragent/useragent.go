package useragent

import (
	"net/http"
	"sync"

	"github.com/deangrant/user-agent/internal/detect"
	detectagent "github.com/deangrant/user-agent/internal/detect/agent"
	detectapp "github.com/deangrant/user-agent/internal/detect/app"
	detectbot "github.com/deangrant/user-agent/internal/detect/bot"
	detectdevice "github.com/deangrant/user-agent/internal/detect/device"
	detectengine "github.com/deangrant/user-agent/internal/detect/engine"
	detectos "github.com/deangrant/user-agent/internal/detect/os"
	"github.com/deangrant/user-agent/internal/merge"
	"github.com/deangrant/user-agent/internal/tokenize"
)

// Analyzer parses User-Agent strings and optional Client Hints.
type Analyzer struct {
	detectors []detect.Detector
}

// Option configures an Analyzer.
type Option func(*Analyzer)

// WithDetectors replaces the default detector pipeline.
// Intended for tests and advanced customization.
func WithDetectors(d ...detect.Detector) Option {
	return func(a *Analyzer) {
		a.detectors = append([]detect.Detector(nil), d...)
	}
}

// NewAnalyzer returns an Analyzer with the default detector pipeline.
func NewAnalyzer(opts ...Option) *Analyzer {
	a := &Analyzer{
		detectors: defaultDetectors(),
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func defaultDetectors() []detect.Detector {
	return []detect.Detector{
		detectbot.New(),
		detectapp.New(),
		detectagent.New(),
		detectengine.New(),
		detectos.New(),
		detectdevice.New(),
	}
}

// Parse analyzes a User-Agent string without Client Hints.
func (a *Analyzer) Parse(ua string) Result {
	return a.ParseWithHints(ua, ClientHints{})
}

// ParseWithHints analyzes a User-Agent string with Client Hints.
func (a *Analyzer) ParseWithHints(ua string, hints ClientHints) Result {
	if a == nil {
		a = NewAnalyzer()
	}
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
		Hints:  hintsToInternal(hints),
	}
	for _, d := range a.detectors {
		d.Detect(state)
	}
	merge.Apply(state)
	return resultFromState(state)
}

// ParseHeaders analyzes User-Agent and Client Hints from headers.
func (a *Analyzer) ParseHeaders(h http.Header) Result {
	ua := ""
	if h != nil {
		ua = h.Get("User-Agent")
	}
	return a.ParseWithHints(ua, ClientHintsFromHeader(h))
}

// ParseRequest analyzes User-Agent and Client Hints from a request.
func (a *Analyzer) ParseRequest(r *http.Request) Result {
	if r == nil {
		return a.Parse("")
	}
	return a.ParseHeaders(r.Header)
}

var (
	defaultOnce sync.Once
	defaultA    *Analyzer
)

func defaultAnalyzer() *Analyzer {
	defaultOnce.Do(func() {
		defaultA = NewAnalyzer()
	})
	return defaultA
}

// Parse analyzes a User-Agent string using the shared default Analyzer.
func Parse(ua string) Result {
	return defaultAnalyzer().Parse(ua)
}

// ParseWithHints analyzes a User-Agent with Client Hints using the
// shared default Analyzer.
func ParseWithHints(ua string, hints ClientHints) Result {
	return defaultAnalyzer().ParseWithHints(ua, hints)
}

// ParseHeaders analyzes headers using the shared default Analyzer.
func ParseHeaders(h http.Header) Result {
	return defaultAnalyzer().ParseHeaders(h)
}

// ParseRequest analyzes a request using the shared default Analyzer.
func ParseRequest(r *http.Request) Result {
	return defaultAnalyzer().ParseRequest(r)
}

func resultFromState(s *detect.State) Result {
	devUnknown := string(DeviceClassUnknown)
	osUnknown := string(OSClassUnknown)
	engUnknown := string(EngineClassUnknown)
	agentUnknown := string(AgentClassUnknown)
	secUnknown := string(AgentSecurityUnknown)
	return Result{
		UserAgent: s.UA,
		Device: Device{
			Class: DeviceClass(orUnknown(s.DeviceClass, devUnknown)),
			Name:  s.DeviceName,
			Brand: s.DeviceBrand,
			CPU:   s.DeviceCPU,
		},
		OperatingSystem: OperatingSystem{
			Class:        OSClass(orUnknown(s.OSClass, osUnknown)),
			Name:         s.OSName,
			Version:      s.OSVersion,
			VersionBuild: s.OSVersionBuild,
			NameVersion:  s.OSNameVersion,
		},
		LayoutEngine: LayoutEngine{
			Class:            EngineClass(orUnknown(s.EngineClass, engUnknown)),
			Name:             s.EngineName,
			Version:          s.EngineVersion,
			VersionMajor:     s.EngineVersionMajor,
			NameVersion:      s.EngineNameVersion,
			NameVersionMajor: s.EngineNameVersionMajor,
		},
		Agent: Agent{
			Class:            AgentClass(orUnknown(s.AgentClass, agentUnknown)),
			Name:             s.AgentName,
			Version:          s.AgentVersion,
			VersionMajor:     s.AgentVersionMajor,
			NameVersion:      s.AgentNameVersion,
			NameVersionMajor: s.AgentNameVersionMajor,
		},
		AgentSecurity: AgentSecurity(
			orUnknown(s.AgentSecurity, secUnknown),
		),
		ClientHintsMismatch: s.ClientHintsMismatch,
	}
}

func orUnknown(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}
