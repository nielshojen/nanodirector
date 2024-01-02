package types

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeviceInstallApplication struct {
	ID          string `gorm:"primaryKey;type:char(36)"`
	ManifestURL string
	DeviceUDID  string
}

type SharedInstallApplication struct {
	ID          string `gorm:"primaryKey;type:char(36)"`
	ManifestURL string
}

type InstallApplicationPayload struct {
	SerialNumbers []string      `json:"serial_numbers,omitempty"`
	DeviceUDIDs   []string      `json:"udids,omitempty"`
	ManifestURLs  []ManifestURL `json:"manifest_urls"`
}

type ManifestURL struct {
	URL           string `json:"url"`
	BootstrapOnly bool   `json:"bootstrap_only"`
}

func (application *DeviceInstallApplication) BeforeCreate(scope *gorm.DB) error {
	application.ID = uuid.NewString()
	return nil
}

func (application *SharedInstallApplication) BeforeCreate(scope *gorm.DB) error {
	application.ID = uuid.NewString()
	return nil
}
