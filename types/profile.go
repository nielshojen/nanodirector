package types

import (
	"github.com/google/uuid"

	"gorm.io/gorm"
)

// DeviceProfile (s) are profiles that are individual to the device.
type DeviceProfile struct {
	// ID                uuid.UUID `gorm:"primaryKey;type:char(36);default:uuid_generate_v4()"`
	PayloadUUID       string
	PayloadIdentifier string `gorm:"primaryKey"`
	HashedPayloadUUID string
	MobileconfigData  []byte
	MobileconfigHash  []byte
	DeviceUDID        string `gorm:"primaryKey"`
	Installed         bool   `gorm:"default:true"`
}

// SharedProfile (s) are profiles that go on every device.
type SharedProfile struct {
	ID                string `gorm:"primaryKey;type:char(36)"`
	PayloadUUID       string
	HashedPayloadUUID string
	PayloadIdentifier string
	MobileconfigData  []byte
	MobileconfigHash  []byte
	Installed         bool `gorm:"default:true"`
}

// ProfilePayload - struct to unpack the payload sent to mdmdirector
type ProfilePayload struct {
	SerialNumbers []string `json:"serial_numbers,omitempty"`
	DeviceUDIDs   []string `json:"udids,omitempty"`
	Mobileconfigs []string `json:"profiles"`
	PushNow       bool     `json:"push_now"`
	Metadata      bool     `json:"metadata"`
}

type DeleteProfilePayload struct {
	SerialNumbers []string                     `json:"serial_numbers,omitempty"`
	DeviceUDIDs   []string                     `json:"udids,omitempty"`
	PushNow       bool                         `json:"push_now"`
	Mobileconfigs []DeletedMobileconfigPayload `json:"profiles"`
	Metadata      bool                         `json:"metadata"`
}

type DeletedMobileconfigPayload struct {
	PayloadIdentifier string `json:"payload_identifier"`
}

// ProfileList - returned data from the ProfileList MDM command
type ProfileListData struct {
	ProfileList []ProfileList
}

type ProfileList struct {
	ID                       string `gorm:"primaryKey;type:char(36)"`
	DeviceUDID               string
	HasRemovalPasscode       bool                 `plist:"HasRemovalPasscode"`
	IsEncrypted              bool                 `plist:"IsEncrypted"`
	IsManaged                bool                 `plist:"IsManaged"`
	PayloadContent           []PayloadContentItem `plist:"PayloadContent" gorm:"-"`
	PayloadDescription       string               `plist:"PayloadDescription"`
	PayloadDisplayName       string               `plist:"PayloadDisplayName"`
	PayloadIdentifier        string               `plist:"PayloadIdentifier"`
	PayloadOrganization      string               `plist:"PayloadOrganization"`
	PayloadRemovalDisallowed bool                 `plist:"PayloadRemovalDisallowed"`
	PayloadUUID              string               `plist:"PayloadUUID" gorm:"not null"`
	PayloadVersion           int                  `plist:"PayloadVersion"`
	FullPayload              bool                 `plist:"FullPayload"`
	SignerCertificates       [][]byte             `plist:"SignerCertificates" gorm:"-"`
}

type MetadataItem struct {
	ProfileMetadata []ProfileMetadata `json:"profile_metadata"`
}

type ProfileMetadata struct {
	Status            string `json:"status"`
	PayloadIdentifier string `json:"payload_identifier"`
	PayloadUUID       string `json:"payload_uuid"`
	HashedPayloadUUID string `json:"hashed_payload_uuid"`
}

// https://developer.apple.com/documentation/devicemanagement/profilelistresponse/profilelistitem/payloadcontentitem
type PayloadContentItem struct {
	PayloadDescription  string `plist:"PayloadDescription"`
	PayloadDisplayName  string `plist:"PayloadDisplayName"`
	PayloadIdentifier   string `plist:"PayloadIdentifier"`
	PayloadOrganization string `plist:"PayloadOrganization"`
	PayloadType         string `plist:"PayloadType"`
	PayloadVersion      int    `plist:"PayloadVersion"`
}

func (profile *SharedProfile) BeforeCreate(scope *gorm.DB) error {
	profile.ID = uuid.NewString()
	return nil
}

func (profilelist *ProfileList) BeforeCreate(scope *gorm.DB) error {
	profilelist.ID = uuid.NewString()
	return nil
}

func (profile *DeviceProfile) AfterCreate(tx *gorm.DB) (err error) {
	BumpDeviceLastUpdated(profile.DeviceUDID)
	return nil
}

func (profile *DeviceProfile) AfterUpdate(tx *gorm.DB) (err error) {
	BumpDeviceLastUpdated(profile.DeviceUDID)
	return nil
}
