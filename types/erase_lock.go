package types

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EscrowPayload struct {
	Serial     string `form:"serial"`
	Pin        string `form:"recovery_password"`
	Username   string `form:"username"`
	SecretType string `form:"secret_type"`
}

type UnlockPin struct {
	ID         string `gorm:"primaryKey;type:char(36)"`
	UnlockPin  string
	PinSet     time.Time `gorm:"type:DATETIME;default:NULL"`
	DeviceUDID string
}

func (unlockpin *UnlockPin) BeforeCreate(scope *gorm.DB) error {
	unlockpin.ID = uuid.NewString()
	return nil
}
