package types

import (
	"time"
)

type Command struct {
	UpdatedAt    time.Time
	CommandUUID  string `gorm:"primaryKey"`
	Status       string
	DeviceUDID   string   `json:"udid"`
	RequestType  string   `json:"request_type"`
	Payload      string   `json:"payload,omitempty"`
	Queries      []string `json:"Queries,omitempty" gorm:"type:TEXT"`
	Identifier   string   `json:"identifier,omitempty"`
	ManifestURL  string   `json:"manifest_url,omitempty"`
	ErrorString  string
	AttemptCount int
}

type Payload struct {
	UDID           string     `plist:"UDID"`
	CommandPayload RawCommand `plist:"CommandPayload"`
}

type RawCommand struct {
	CommandUUID string         `plist:"CommandUUID"`
	Command     CommandPayload `plist:"Command"`
}

type CommandPayload struct {
	RequestType                         string                 `plist:"RequestType"`
	Payload                             []byte                 `plist:"Payload,omitempty"`
	Queries                             []string               `plist:"Queries,omitempty"`
	Identifier                          string                 `plist:"Identifier,omitempty"`
	ManifestURL                         string                 `plist:"ManifestURL,omitempty"`
	Pin                                 string                 `plist:"PIN,omitempty"`
	SkipPrimarySetupAccountCreation     bool                   `plist:"SkipPrimarySetupAccountCreation,omitempty"`
	SetPrimarySetupAccountAsRegularUser bool                   `plist:"SetPrimarySetupAccountAsRegularUser,omitempty"`
	DontAutoPopulatePrimaryAccountInfo  bool                   `plist:"DontAutoPopulatePrimaryAccountInfo,omitempty"`
	LockPrimaryAccountInfo              bool                   `plist:"LockPrimaryAccountInfo,omitempty"`
	PrimaryAccountFullName              string                 `plist:"PrimaryAccountFullName,omitempty"`
	PrimaryAccountUserName              string                 `plist:"PrimaryAccountUserName,omitempty"`
	RequestRequiresNetworkTether        bool                   `plist:"RequestRequiresNetworkTether,omitempty"`
	ManagedLocalUserShortName           string                 `plist:"ManagedLocalUserShortName,omitempty"`
	AutoSetupAdminAccounts              AutoSetupAdminAccounts `plist:"AutoSetupAdminAccounts,omitempty" gorm:"foreignKey:UDID;references:UDID"`
}

type AutoSetupAdminAccounts struct {
	ShortName    string `plist:"ShortName,omitempty"`
	FullName     string `plist:"FullName,omitempty"`
	PasswordHash []byte `plist:"PasswordHash,omitempty"`
	Hidden       bool   `plist:"Hidden,omitempty"`
}

type CommandResponse struct {
	CommandUUID string        `json:"command_uuid"`
	RequestType string        `json:"request_type"`
	Status      CommandStatus `json:"status"`
}

type CommandStatus struct {
	DeviceUDID DeviceResult
}

type DeviceResult struct {
	CommandUUID string `json:"push_result"`
}
