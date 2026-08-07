// Package detect defines the analysis state and Detector contract.
package detect

import (
	"github.com/deangrant/user-agent/internal/hintparse"
	"github.com/deangrant/user-agent/internal/tokenize"
)

// Detector fills fields on State. Implementations should only set
// fields they can improve and avoid clobbering stronger prior results
// unless they own that concern.
type Detector interface {
	Detect(state *State)
}

// State is the mutable analysis scratchpad shared across detectors.
type State struct {
	UA     string
	Tokens tokenize.Tokens
	Hints  hintparse.Hints

	DeviceClass string
	DeviceName  string
	DeviceBrand string
	DeviceCPU   string

	OSClass        string
	OSName         string
	OSVersion      string
	OSVersionBuild string
	OSNameVersion  string

	EngineClass            string
	EngineName             string
	EngineVersion          string
	EngineVersionMajor     string
	EngineNameVersion      string
	EngineNameVersionMajor string

	AgentClass            string
	AgentName             string
	AgentVersion          string
	AgentVersionMajor     string
	AgentNameVersion      string
	AgentNameVersionMajor string

	AgentSecurity string

	// ClientHintsMismatch is set when Client Hints conflict with
	// UA-derived OS family (e.g. Android UA + Windows CH).
	ClientHintsMismatch bool

	// BotMatched is set when a known bot pattern matched.
	BotMatched bool
	// AppMatched is set when a known app/webview pattern matched.
	AppMatched bool
}

// SetAgent sets agent identity fields and derived name/version strings.
func (s *State) SetAgent(class, name, ver string) {
	if class != "" {
		s.AgentClass = class
	}
	if name != "" {
		s.AgentName = name
	}
	if ver != "" {
		s.AgentVersion = ver
	}
}

// SetEngine sets layout engine fields.
func (s *State) SetEngine(class, name, ver string) {
	if class != "" {
		s.EngineClass = class
	}
	if name != "" {
		s.EngineName = name
	}
	if ver != "" {
		s.EngineVersion = ver
	}
}

// SetOS sets operating system fields.
func (s *State) SetOS(class, name, ver string) {
	if class != "" {
		s.OSClass = class
	}
	if name != "" {
		s.OSName = name
	}
	if ver != "" {
		s.OSVersion = ver
	}
}

// SetDevice sets device fields.
func (s *State) SetDevice(class, name, brand string) {
	if class != "" {
		s.DeviceClass = class
	}
	if name != "" {
		s.DeviceName = name
	}
	if brand != "" {
		s.DeviceBrand = brand
	}
}
