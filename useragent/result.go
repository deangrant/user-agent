package useragent

import "github.com/deangrant/user-agent/internal/uaclass"

// DeviceClass classifies the hardware that made the request.
type DeviceClass string

// Device class values.
const (
	DeviceClassDesktop             DeviceClass = uaclass.Desktop
	DeviceClassAnonymized          DeviceClass = uaclass.Anonymized
	DeviceClassUnknown             DeviceClass = uaclass.Unknown
	DeviceClassMobile              DeviceClass = uaclass.Mobile
	DeviceClassTablet              DeviceClass = uaclass.Tablet
	DeviceClassPhone               DeviceClass = uaclass.Phone
	DeviceClassWatch               DeviceClass = uaclass.Watch
	DeviceClassAugmentedReality    DeviceClass = uaclass.AugmentedReality
	DeviceClassVirtualReality      DeviceClass = uaclass.VirtualReality
	DeviceClassEReader             DeviceClass = uaclass.EReader
	DeviceClassSetTopBox           DeviceClass = uaclass.SetTopBox
	DeviceClassTV                  DeviceClass = uaclass.TV
	DeviceClassHomeAppliance       DeviceClass = uaclass.HomeAppliance
	DeviceClassGameConsole         DeviceClass = uaclass.GameConsole
	DeviceClassHandheldGameConsole DeviceClass = uaclass.HandheldGameConsole
	DeviceClassVoice               DeviceClass = uaclass.Voice
	DeviceClassSmartDisplay        DeviceClass = uaclass.SmartDisplay
	DeviceClassCar                 DeviceClass = uaclass.Car
	DeviceClassRobot               DeviceClass = uaclass.Robot
	DeviceClassRobotMobile         DeviceClass = uaclass.RobotMobile
	DeviceClassRobotImitator       DeviceClass = uaclass.RobotImitator
	DeviceClassCloud               DeviceClass = uaclass.Cloud
	DeviceClassHacker              DeviceClass = uaclass.Hacker
)

// OSClass classifies the operating system.
type OSClass string

// Operating system class values.
const (
	OSClassDesktop     OSClass = uaclass.Desktop
	OSClassMobile      OSClass = uaclass.Mobile
	OSClassCloud       OSClass = uaclass.Cloud
	OSClassEmbedded    OSClass = uaclass.Embedded
	OSClassGameConsole OSClass = uaclass.GameConsole
	OSClassHacker      OSClass = uaclass.Hacker
	OSClassAnonymized  OSClass = uaclass.Anonymized
	OSClassUnknown     OSClass = uaclass.Unknown
)

// EngineClass classifies the layout engine.
type EngineClass string

// Layout engine class values.
const (
	EngineClassBrowser    EngineClass = uaclass.Browser
	EngineClassDesktopApp EngineClass = uaclass.DesktopApp
	EngineClassMobileApp  EngineClass = uaclass.MobileApp
	EngineClassHacker     EngineClass = uaclass.Hacker
	EngineClassRobot      EngineClass = uaclass.Robot
	EngineClassCloud      EngineClass = uaclass.Cloud
	EngineClassSpecial    EngineClass = uaclass.Special
	EngineClassUnknown    EngineClass = uaclass.Unknown
)

// AgentClass classifies the user agent application.
type AgentClass string

// Agent class values.
const (
	AgentClassBrowser        AgentClass = uaclass.Browser
	AgentClassBrowserWebview AgentClass = uaclass.BrowserWebview
	AgentClassDesktopApp     AgentClass = uaclass.DesktopApp
	AgentClassMobileApp      AgentClass = uaclass.MobileApp
	AgentClassRobot          AgentClass = uaclass.Robot
	AgentClassRobotMobile    AgentClass = uaclass.RobotMobile
	AgentClassCloudApp       AgentClass = uaclass.CloudApp
	AgentClassServer         AgentClass = uaclass.Server
	AgentClassEmailClient    AgentClass = uaclass.EmailClient
	AgentClassVoice          AgentClass = uaclass.Voice
	AgentClassSpecial        AgentClass = uaclass.Special
	AgentClassTestClient     AgentClass = uaclass.TestClient
	AgentClassHacker         AgentClass = uaclass.Hacker
	AgentClassUnknown        AgentClass = uaclass.Unknown
)

// AgentSecurity classifies indicated transport security.
type AgentSecurity string

// Agent security values.
const (
	AgentSecurityNone    AgentSecurity = uaclass.SecurityNone
	AgentSecurityWeak    AgentSecurity = uaclass.SecurityWeak
	AgentSecurityStrong  AgentSecurity = uaclass.SecurityStrong
	AgentSecurityUnknown AgentSecurity = uaclass.Unknown
	AgentSecurityHacker  AgentSecurity = uaclass.Hacker
)

// Result holds analyzed fields extracted from a User-Agent and
// optional Client Hints.
type Result struct {
	// UserAgent is the input User-Agent string.
	UserAgent string
	// Device describes the client hardware.
	Device Device
	// OperatingSystem describes the client OS.
	OperatingSystem OperatingSystem
	// LayoutEngine describes the rendering engine when detectable.
	LayoutEngine LayoutEngine
	// Agent describes the browser or application.
	Agent Agent
	// AgentSecurity is the indicated transport security class.
	AgentSecurity AgentSecurity
	// ClientHintsMismatch is true when Client Hints conflict with the
	// UA-derived OS family and the UA OS was kept.
	ClientHintsMismatch bool
}

// Device describes the client hardware.
type Device struct {
	// Class is the device category (phone, desktop, robot, …).
	Class DeviceClass
	// Name is the model or product name when known.
	Name string
	// Brand is the manufacturer when known.
	Brand string
	// CPU is the architecture hint when known (e.g. ARM, x64).
	CPU string
}

// OperatingSystem describes the client OS.
type OperatingSystem struct {
	// Class is the OS category (desktop, mobile, …).
	Class OSClass
	// Name is the OS product name (e.g. Windows, Android).
	Name string
	// Version is the OS version string.
	Version string
	// VersionBuild is a build identifier when present (e.g. Android).
	VersionBuild string
	// NameVersion joins Name and Version when both exist.
	NameVersion string
}

// LayoutEngine describes the rendering engine.
type LayoutEngine struct {
	// Class is the engine category (browser, robot, mobile app, …).
	Class EngineClass
	// Name is the engine product (Blink, WebKit, Gecko, …).
	Name string
	// Version is the full engine version when known.
	Version string
	// VersionMajor is the major version component of Version.
	VersionMajor string
	// NameVersion joins Name and Version when both exist.
	NameVersion string
	// NameVersionMajor joins Name and VersionMajor when both exist.
	NameVersionMajor string
}

// Agent describes the browser or application.
type Agent struct {
	// Class is the agent category (browser, robot, mobile app, …).
	Class AgentClass
	// Name is the agent product name.
	Name string
	// Version is the full agent version when known.
	Version string
	// VersionMajor is the major version component of Version.
	VersionMajor string
	// NameVersion joins Name and Version when both exist.
	NameVersion string
	// NameVersionMajor joins Name and VersionMajor when both exist.
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

// IsComputer is equivalent to IsDesktop: both report whether
// Device.Class is desktop (including laptops).
func (r Result) IsComputer() bool {
	return r.IsDesktop()
}
