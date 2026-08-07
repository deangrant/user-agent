package useragent

// DeviceClass classifies the hardware that made the request.
type DeviceClass string

// Device class values.
const (
	DeviceClassDesktop             DeviceClass = "Desktop"
	DeviceClassAnonymized          DeviceClass = "Anonymized"
	DeviceClassUnknown             DeviceClass = "Unknown"
	DeviceClassMobile              DeviceClass = "Mobile"
	DeviceClassTablet              DeviceClass = "Tablet"
	DeviceClassPhone               DeviceClass = "Phone"
	DeviceClassWatch               DeviceClass = "Watch"
	DeviceClassAugmentedReality    DeviceClass = "Augmented Reality"
	DeviceClassVirtualReality      DeviceClass = "Virtual Reality"
	DeviceClassEReader             DeviceClass = "eReader"
	DeviceClassSetTopBox           DeviceClass = "Set-top box"
	DeviceClassTV                  DeviceClass = "TV"
	DeviceClassHomeAppliance       DeviceClass = "Home Appliance"
	DeviceClassGameConsole         DeviceClass = "Game Console"
	DeviceClassHandheldGameConsole DeviceClass = "Handheld Game Console"
	DeviceClassVoice               DeviceClass = "Voice"
	DeviceClassSmartDisplay        DeviceClass = "Smart Display"
	DeviceClassCar                 DeviceClass = "Car"
	DeviceClassRobot               DeviceClass = "Robot"
	DeviceClassRobotMobile         DeviceClass = "Robot Mobile"
	DeviceClassRobotImitator       DeviceClass = "Robot Imitator"
	DeviceClassCloud               DeviceClass = "Cloud"
	DeviceClassHacker              DeviceClass = "Hacker"
)

// OSClass classifies the operating system.
type OSClass string

// Operating system class values.
const (
	OSClassDesktop     OSClass = "Desktop"
	OSClassMobile      OSClass = "Mobile"
	OSClassCloud       OSClass = "Cloud"
	OSClassEmbedded    OSClass = "Embedded"
	OSClassGameConsole OSClass = "Game Console"
	OSClassHacker      OSClass = "Hacker"
	OSClassAnonymized  OSClass = "Anonymized"
	OSClassUnknown     OSClass = "Unknown"
)

// EngineClass classifies the layout engine.
type EngineClass string

// Layout engine class values.
const (
	EngineClassBrowser    EngineClass = "Browser"
	EngineClassDesktopApp EngineClass = "Desktop App"
	EngineClassMobileApp  EngineClass = "Mobile App"
	EngineClassHacker     EngineClass = "Hacker"
	EngineClassRobot      EngineClass = "Robot"
	EngineClassCloud      EngineClass = "Cloud"
	EngineClassSpecial    EngineClass = "Special"
	EngineClassUnknown    EngineClass = "Unknown"
)

// AgentClass classifies the user agent application.
type AgentClass string

// Agent class values.
const (
	AgentClassBrowser        AgentClass = "Browser"
	AgentClassBrowserWebview AgentClass = "Browser Webview"
	AgentClassDesktopApp     AgentClass = "Desktop App"
	AgentClassMobileApp      AgentClass = "Mobile App"
	AgentClassRobot          AgentClass = "Robot"
	AgentClassRobotMobile    AgentClass = "Robot Mobile"
	AgentClassCloudApp       AgentClass = "Cloud Application"
	AgentClassServer         AgentClass = "Server"
	AgentClassEmailClient    AgentClass = "Email Client"
	AgentClassVoice          AgentClass = "Voice"
	AgentClassSpecial        AgentClass = "Special"
	AgentClassTestClient     AgentClass = "Testclient"
	AgentClassHacker         AgentClass = "Hacker"
	AgentClassUnknown        AgentClass = "Unknown"
)

// AgentSecurity classifies indicated transport security.
type AgentSecurity string

// Agent security values.
const (
	AgentSecurityNone    AgentSecurity = "No security"
	AgentSecurityWeak    AgentSecurity = "Weak security"
	AgentSecurityStrong  AgentSecurity = "Strong security"
	AgentSecurityUnknown AgentSecurity = "Unknown"
	AgentSecurityHacker  AgentSecurity = "Hacker"
)

// Result holds analyzed fields extracted from a User-Agent and
// optional Client Hints.
type Result struct {
	UserAgent       string
	Device          Device
	OperatingSystem OperatingSystem
	LayoutEngine    LayoutEngine
	Agent           Agent
	AgentSecurity   AgentSecurity
}

// Device describes the client hardware.
type Device struct {
	Class DeviceClass
	Name  string
	Brand string
	CPU   string
}

// OperatingSystem describes the client OS.
type OperatingSystem struct {
	Class        OSClass
	Name         string
	Version      string
	VersionBuild string
	NameVersion  string
}

// LayoutEngine describes the rendering engine.
type LayoutEngine struct {
	Class            EngineClass
	Name             string
	Version          string
	VersionMajor     string
	NameVersion      string
	NameVersionMajor string
}

// Agent describes the browser or application.
type Agent struct {
	Class            AgentClass
	Name             string
	Version          string
	VersionMajor     string
	NameVersion      string
	NameVersionMajor string
}

// IsMobile reports whether the device is a phone, tablet, watch, or
// generic mobile class.
func (r Result) IsMobile() bool {
	switch r.Device.Class {
	case DeviceClassPhone, DeviceClassTablet, DeviceClassWatch,
		DeviceClassMobile, DeviceClassEReader,
		DeviceClassHandheldGameConsole:
		return true
	default:
		return false
	}
}

// IsTablet reports whether the device is classified as a tablet.
func (r Result) IsTablet() bool {
	return r.Device.Class == DeviceClassTablet
}

// IsDesktop reports whether the device is a desktop or laptop.
func (r Result) IsDesktop() bool {
	return r.Device.Class == DeviceClassDesktop
}

// IsPhone reports whether the device is classified as a phone.
func (r Result) IsPhone() bool {
	return r.Device.Class == DeviceClassPhone
}

// IsBot reports whether the agent or device looks like non-browser
// automation: robots/crawlers, HTTP libraries and other server clients,
// cloud applications, hacker tools, or test clients. Intended for
// abuse and security filtering as well as crawler detection.
func (r Result) IsBot() bool {
	switch r.Device.Class {
	case DeviceClassRobot, DeviceClassRobotMobile,
		DeviceClassRobotImitator, DeviceClassCloud,
		DeviceClassHacker:
		return true
	}
	switch r.Agent.Class {
	case AgentClassRobot, AgentClassRobotMobile,
		AgentClassServer, AgentClassCloudApp,
		AgentClassHacker, AgentClassTestClient:
		return true
	}
	return false
}

// IsComputer reports whether the device is a general-purpose computer
// (desktop class).
func (r Result) IsComputer() bool {
	return r.IsDesktop()
}
