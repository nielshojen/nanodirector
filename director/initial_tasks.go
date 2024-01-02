package director

import (
	"time"

	"github.com/nielshojen/nanodirector/db"
	"github.com/nielshojen/nanodirector/types"
	"github.com/pkg/errors"
)

func RunInitialTasks(udid string) error {
	if udid == "" {
		err := errors.New("No Device UDID")
		return errors.Wrap(err, "RunInitialTasks")
	}

	device, err := GetDevice(udid)
	if err != nil {
		return errors.Wrap(err, "RunInitialTasks")
	}
	// if device.InitialTasksRun == true {
	// 	log.Infof("Initial tasks already run for %v", device.UDID)
	// 	return nil
	// }
	InfoLogger(LogHolder{Message: "Running initial tasks", DeviceSerial: device.SerialNumber, DeviceUDID: device.UDID})
	err = ClearCommands(&device)
	if err != nil {
		return err
	}

	// if device.Erase || device.Lock {
	// 	// Got a device checking in that should be wiped or locked. Make it so.
	// 	err = EraseLockDevice(&device)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	return nil
	// }

	_, err = InstallAllProfiles(device)
	if err != nil {
		return errors.Wrap(err, "RunInitialTasks:InstallAllProfiles")
	}

	_, err = InstallBootstrapPackages(device)
	if err != nil {
		return errors.Wrap(err, "RunInitialTasks:InstallBootstrapPackages")
	}
	err = processDeviceConfigured(device)
	if err != nil {
		return errors.Wrap(err, "RunInitialTasks:processDeviceConfigured")
	}

	return nil
}

func processDeviceConfigured(device types.Device) error {
	var deviceModel types.Device
	err := SendDeviceConfigured(device)
	if err != nil {
		return errors.Wrap(err, "RunInitialTasks")
	}
	err = SaveDeviceConfigured(device)
	if err != nil {
		return err
	}
	err = db.DB.Model(&deviceModel).Select("last_info_requested").Where("ud_id = ?", device.UDID).Updates(map[string]interface{}{"last_info_requested": time.Now()}).Error
	if err != nil {
		return err
	}

	return nil
}

func SendDeviceConfigured(device types.Device) error {
	requestType := "DeviceConfigured"
	var payload types.Payload
	payload.UDID = device.UDID
	payload.CommandPayload.Command.RequestType = requestType
	_, err := SendCommand(payload)
	if err != nil {
		return errors.Wrap(err, "SendDeviceConfigured")
	}
	// Twice for luck
	_, err = SendCommand(payload)
	if err != nil {
		return errors.Wrap(err, "SendDeviceConfigured")
	}
	return nil
}

func SaveDeviceConfigured(device types.Device) error {
	var deviceModel types.Device
	now := time.Now()
	err := db.DB.Model(&deviceModel).Select("token_update_recieved", "authenticate_recieved", "initial_tasks_run", "last_checked_in", "next_push").Where("ud_id = ?", device.UDID).Updates(map[string]interface{}{"token_update_recieved": true, "authenticate_recieved": true, "initial_tasks_run": true, "last_checked_in": now, "next_push": now}).Error
	if err != nil {
		return err
	}

	return nil
}

func ResetDevice(device types.Device) error {
	var deviceModel types.Device
	err := ClearCommands(&device)
	if err != nil {
		return errors.Wrap(err, "ResetDevice:ClearCommands")
	}
	InfoLogger(LogHolder{DeviceUDID: device.UDID, DeviceSerial: device.SerialNumber, Message: "Resetting device"})
	err = db.DB.Model(&deviceModel).Where("ud_id = ?", device.UDID).Updates(map[string]interface{}{"token_update_recieved": false, "authenticate_recieved": false, "initial_tasks_run": false, "active": false}).Error
	if err != nil {
		return errors.Wrap(err, "reset device")
	}

	// err = db.DB.Unscoped().Where("device_ud_id = ?", device.UDID).Delete(types.Certificate{}).Error
	// if err != nil {
	// 	ErrorLogger(LogHolder{Message: err.Error()})
	// }

	// err = db.DB.Unscoped().Where("device_ud_id = ?", device.UDID).Delete(types.ProfileList{}).Error
	// if err != nil {
	// 	ErrorLogger(LogHolder{Message: err.Error()})
	// }

	return nil
}
