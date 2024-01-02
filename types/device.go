package types

import "time"

type Device struct {
	DeviceName                       string
	Username                         string
	Fullname                         string
	BuildVersion                     string
	ModelName                        string
	Model                            string
	OSVersion                        string
	ProductName                      string
	SerialNumber                     string
	DeviceCapacity                   float32
	AvailableDeviceCapacity          float32
	BatteryLevel                     float32
	CellularTechnology               int
	IMEI                             string
	MEID                             string
	ModemFirmwareVersion             string
	IsSupervised                     bool
	IsDeviceLocatorServiceEnabled    bool
	IsActivationLockEnabled          bool
	IsDoNotDisturbInEffect           bool
	IsMDMLostModeEnabled             bool
	DeviceID                         string
	EASDeviceIdentifier              string
	IsCloudBackupEnabled             bool
	OSUpdateSettings                 OSUpdateSettings `gorm:"foreignKey:DeviceUDID" json:",omitempty"`
	LocalHostName                    string
	HostName                         string
	SystemIntegrityProtectionEnabled bool
	// ActiveManagedUsers               []string
	AppAnalyticsEnabled bool
	// AutoSetupAdminAccounts interface
	AwaitingConfiguration       bool `gorm:"default:false"`
	MaximumResidentUsers        int
	BluetoothMAC                string
	CarrierSettingsVersion      string
	CurrentCarrierNetwork       string
	CurrentMCC                  string
	CurrentMNC                  string
	DataRoamingEnabled          string
	DiagnosticSubmissionEnabled bool
	ICCID                       string
	IsMultiUser                 bool
	IsNetworkTethered           bool
	IsRoaming                   bool
	iTunesStoreAccountHash      string
	iTunesStoreAccountIsActive  bool
	LastCloudBackupDate         time.Time `gorm:"type:DATETIME;default:NULL"`
	//MDMOptions                  string
	// EthernetMACs []string
	// OrganizationInfo interface{}
	PersonalHotspotEnabled bool
	PhoneNumber            string
	PushToken              string
	// ServiceSubscriptions interface{}
	SIMCarrierNetwork          string
	SIMMCC                     string
	SIMMNC                     string
	SubscriberCarrierNetwork   string
	SubscriberMCC              string
	SubscriberMNC              string
	SupplementalBuildVersion   string
	SupplementalOSVersionExtra string
	// SoftwareUpdateSettings     SoftwareUpdateSettings `gorm:"foreignKey:DeviceUDID" json:",omitempty"`
	VoiceRoamingEnabled  bool
	WiFiMAC              string
	EthernetMAC          string
	UDID                 string `gorm:"primaryKey"`
	Active               bool
	Profiles             []DeviceProfile            `gorm:"foreignKey:DeviceUDID" json:",omitempty"`
	Commands             []Command                  `gorm:"foreignKey:DeviceUDID" json:",omitempty"`
	Certificates         []Certificate              `gorm:"foreignKey:DeviceUDID" json:",omitempty"`
	InstallApplications  []DeviceInstallApplication `gorm:"foreignKey:DeviceUDID" json:",omitempty"`
	SecurityInfo         SecurityInfo               `gorm:"foreignKey:DeviceUDID" json:",omitempty"`
	ProfileList          []ProfileList              `gorm:"foreignKey:DeviceUDID" json:",omitempty"`
	UpdatedAt            time.Time                  `gorm:"type:DATETIME;default:NULL"`
	AuthenticateRecieved bool                       `gorm:"default:false"`
	TokenUpdateRecieved  bool                       `gorm:"default:false"`
	InitialTasksRun      bool                       `gorm:"default:false"`
	Erase                bool                       `gorm:"default:false"`
	Lock                 bool                       `gorm:"default:false"`
	UnlockPin            string
	TempUnlockPin        UnlockPin `gorm:"foreignKey:DeviceUDID"`
	LastInfoRequested    time.Time `gorm:"type:DATETIME;default:NULL"`
	NextPush             time.Time `gorm:"type:DATETIME;default:NULL"`
	// LastScheduledPush        time.Time  `gorm:"type:DATETIME;default:NULL"`
	LastCheckedIn       time.Time `gorm:"type:DATETIME;default:NULL"`
	LastCertificateList time.Time `gorm:"type:DATETIME;default:NULL"`
	LastProfileList     time.Time `gorm:"type:DATETIME;default:NULL"`
	LastDeviceInfo      time.Time `gorm:"type:DATETIME;default:NULL"`
	LastSecurityInfo    time.Time `gorm:"type:DATETIME;default:NULL"`
}

var DeviceInformationQueries = []string{
	"ActiveManagedUsers",
	"AppAnalyticsEnabled",
	"AutoSetupAdminAccounts",
	"AvailableDeviceCapacity",
	"AwaitingConfiguration",
	"BatteryLevel",
	"BluetoothMAC",
	"BuildVersion",
	"CarrierSettingsVersion",
	"CellularTechnology",
	"CurrentMCC",
	"CurrentMNC",
	"DataRoamingEnabled",
	"DeviceCapacity",
	"DeviceID",
	"DeviceName",
	"DiagnosticSubmissionEnabled",
	"EASDeviceIdentifier",
	"ICCID",
	"IMEI",
	"IsActivationLockEnabled",
	"IsCloudBackupEnabled",
	"IsDeviceLocatorServiceEnabled",
	"IsDoNotDisturbInEffect",
	"IsMDMLostModeEnabled",
	"IsMultiUser",
	"IsNetworkTethered",
	"IsRoaming",
	"IsSupervised",
	"iTunesStoreAccountHash",
	"iTunesStoreAccountIsActive",
	"LastCloudBackupDate",
	"MaximumResidentUsers",
	"MDMOptions",
	"MEID",
	"Model",
	"ModelName",
	"ModemFirmwareVersion",
	"OrganizationInfo",
	"OSUpdateSettings",
	"OSVersion",
	"PersonalHotspotEnabled",
	"PhoneNumber",
	"ProductName",
	"PushToken",
	"SerialNumber",
	"ServiceSubscriptions",
	"SIMCarrierNetwork",
	"SIMMCC",
	"SIMMNC",
	"SubscriberCarrierNetwork",
	"SubscriberMCC",
	"SubscriberMNC",
	"SystemIntegrityProtectionEnabled",
	"UDID",
	"VoiceRoamingEnabled",
	"WiFiMAC",
	"EthernetMAC",
	"SupplementalBuildVersion",
	"SupplementalOSVersionExtra",
	"SoftwareUpdateSettings",
}

type OSUpdateSettings struct {
	DeviceUDID                      string `gorm:"primaryKey"`
	CatalogURL                      string
	IsDefaultCatalog                bool
	PreviousScanDate                time.Time `gorm:"type:DATETIME;default:NULL"`
	PreviousScanResult              int
	PerformPeriodicCheck            bool
	AutomaticCheckEnabled           bool
	BackgroundDownloadEnabled       bool
	AutomaticAppInstallationEnabled bool
	AutomaticOSInstallationEnabled  bool
	AutomaticSecurityUpdatesEnabled bool
}

type DeviceInformationQueryResponses struct {
	QueryResponses Device `plist:"QueryResponses"`
}

type DeviceCommandPayload struct {
	SerialNumbers []string `json:"serial_numbers,omitempty"`
	DeviceUDIDs   []string `json:"udids,omitempty"`
	Value         bool     `json:"value"`
	PushNow       bool     `json:"push_now"`
	Metadata      bool     `json:"metadata"`
	Pin           string   `json:"pin,omitempty"`
}

// type SoftwareUpdateSettings struct {
// 	RecommendationsCadence int
// 	DeviceUDID             string `gorm:"primaryKey"`
// }
