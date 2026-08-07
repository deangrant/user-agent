// Package uaclass holds shared UA analysis class string constants.
// Public typed aliases in package useragent re-export these values.
package uaclass

// Shared and cross-cutting class values.
const (
	Unknown = "Unknown"
	Desktop = "Desktop"
	Mobile  = "Mobile"
	Browser = "Browser"
	Robot   = "Robot"
	Hacker  = "Hacker"
	Cloud   = "Cloud"
	Special = "Special"
)

// Device class values.
const (
	Anonymized          = "Anonymized"
	Tablet              = "Tablet"
	Phone               = "Phone"
	Watch               = "Watch"
	AugmentedReality    = "Augmented Reality"
	VirtualReality      = "Virtual Reality"
	EReader             = "eReader"
	SetTopBox           = "Set-top box"
	TV                  = "TV"
	HomeAppliance       = "Home Appliance"
	GameConsole         = "Game Console"
	HandheldGameConsole = "Handheld Game Console"
	Voice               = "Voice"
	SmartDisplay        = "Smart Display"
	Car                 = "Car"
	RobotMobile         = "Robot Mobile"
	RobotImitator       = "Robot Imitator"
)

// OS-only class values (Desktop/Mobile/Cloud/Hacker/Unknown shared above).
const (
	Embedded = "Embedded"
)

// Engine/agent application class values.
const (
	DesktopApp     = "Desktop App"
	MobileApp      = "Mobile App"
	BrowserWebview = "Browser Webview"
	CloudApp       = "Cloud Application"
	Server         = "Server"
	EmailClient    = "Email Client"
	TestClient     = "Testclient"
)

// Agent security values.
const (
	SecurityNone   = "No security"
	SecurityWeak   = "Weak security"
	SecurityStrong = "Strong security"
)
