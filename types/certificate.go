package types

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Certificate represents a certificate.
type Certificate struct {
	ID         string `gorm:"primaryKey;type:char(36)"`
	CommonName string
	Subject    string
	NotAfter   time.Time `gorm:"type:DATETIME;default:NULL"`
	NotBefore  time.Time `gorm:"type:DATETIME;default:NULL"`
	Data       []byte
	Issuer     string
	DeviceUDID string
}

// CertificateListData - returned data from the CertificateList MDM command
type CertificateListData struct {
	CertificateList []CertificateList
}

// CertificateList Each item from CertificateList
type CertificateList struct {
	CommonName string `plist:"CommonName"`
	Data       []byte `plist:"Data"`
	IsIdentity bool   `plist:"IsIdentity"`
}

func (certificate *Certificate) BeforeCreate(scope *gorm.DB) error {
	certificate.ID = uuid.NewString()
	return nil
}
