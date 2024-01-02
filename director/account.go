package director

import (
	"fmt"

	"github.com/groob/plist"
	"github.com/nielshojen/nanodirector/db"
	"github.com/nielshojen/nanodirector/types"
	"github.com/nielshojen/nanodirector/utils"
	"github.com/pkg/errors"
)

func GetAccount(serialnumber string) (types.Account, error) {
	var account types.Account

	if serialnumber == "" {
		err := fmt.Errorf("No serialnumber set")
		return account, errors.Wrap(err, "GetAccount")
	}

	err := db.DB.Model(account).Where("serial_number = ?", serialnumber).First(&account).Scan(&account).Error
	if err != nil {
		return account, errors.Wrapf(
			err,
			"Couldn't scan to Account model from GetAccount %v",
			serialnumber,
		)
	}
	return account, nil
}

func SetAccount(account types.Account) (types.Account, error) {
	var accountModel types.Account
	DebugLogger(
		LogHolder{
			Message:    "Account Create Received",
			DeviceUDID: account.UDID,
		},
	)

	err := db.DB.Model(&accountModel).Where("ud_id = ?", account.UDID).Assign(&account).FirstOrCreate(&account).Error
	if err != nil {
		return account, errors.Wrap(err, "Updating Account")
	}
	return account, nil
}

func AccountConfiguration(device types.Device) (types.Device, error) {

	account, err := GetAccount(device.SerialNumber)
	if err != nil {
		return device, errors.Wrap(err, "No Account found for Device")
	}

	var payload types.Payload
	payload.UDID = device.UDID

	AdminUserFullname := "Administrator"
	AdminUserUsername := utils.AdminUserUsername()
	AdminUserPassword := utils.AdminUserPassword()

	payload.CommandPayload.Command.RequestType = "AccountConfiguration"

	if AdminUserPassword != "" {
		salted, err := utils.SaltedSHA512PBKDF2(AdminUserPassword)
		if err != nil {
			return device, errors.Wrap(err, "salting plaintext password")
		}
		hashDict := struct {
			SaltedSHA512PBKDF2 utils.SaltedSHA512PBKDF2Dictionary `plist:"SALTED-SHA512-PBKDF2"`
		}{
			SaltedSHA512PBKDF2: salted,
		}
		hashPlist, err := plist.Marshal(hashDict)
		if err != nil {
			return device, errors.Wrap(err, "marshal salted password to plist")
		}
		payload.CommandPayload.Command.AutoSetupAdminAccounts.PasswordHash = hashPlist

		payload.CommandPayload.Command.AutoSetupAdminAccounts.FullName = AdminUserFullname
		payload.CommandPayload.Command.AutoSetupAdminAccounts.ShortName = AdminUserUsername
		payload.CommandPayload.Command.AutoSetupAdminAccounts.Hidden = false
	}

	payload.CommandPayload.Command.DontAutoPopulatePrimaryAccountInfo = false
	payload.CommandPayload.Command.LockPrimaryAccountInfo = true
	payload.CommandPayload.Command.PrimaryAccountFullName = account.PrimaryAccountFullName
	payload.CommandPayload.Command.PrimaryAccountUserName = account.PrimaryAccountUserName
	payload.CommandPayload.Command.RequestRequiresNetworkTether = false
	payload.CommandPayload.Command.SetPrimarySetupAccountAsRegularUser = true
	payload.CommandPayload.Command.SkipPrimarySetupAccountCreation = false
	payload.CommandPayload.Command.ManagedLocalUserShortName = AdminUserUsername

	_, err = SendCommand(payload)
	if err != nil {
		return device, errors.Wrap(err, "Could not set Account for Device")
	}
	return device, nil
}
