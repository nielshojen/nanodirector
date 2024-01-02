package types

type Account struct {
	UDID                   string `json:"udid"`
	SerialNumber           string `gorm:"primaryKey,unique" json:"serialnumber"`
	PrimaryAccountFullName string `json:"primary_account_full_name,omitempty"`
	PrimaryAccountUserName string `json:"primary_account_user_name,omitempty"`
}
