package director

import (
	"crypto/rand"
	intErrors "errors"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gopkg.in/ajg/form.v1"

	"github.com/mdmdirector/mdmdirector/db"
	"github.com/mdmdirector/mdmdirector/types"
	"github.com/mdmdirector/mdmdirector/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

func EraseLockDevice(udid string) error {

	device, err := GetDevice(udid)
	if err != nil {
		ErrorLogger(LogHolder{Message: err.Error()})
	}

	var requestType string

	pin, err := generatePin(device)

	if err != nil {
		return errors.Wrap(err, "EraseLockDevice::generatePin")
	}

	if device.Lock {
		requestType = "DeviceLock"
	}

	// Erase wins if both are set
	if device.Erase {
		requestType = "EraseDevice"
	}

	if requestType == "" {
		log.Info("Neither lock or erase are set")
		return nil
	}

	err = escrowPin(device, pin)
	if err != nil {
		return errors.Wrap(err, "EraseLockDevice:escrowPin")
	}
	log.Infof("Sending %v to %v", requestType, device.UDID)
	var payload types.Payload
	payload.UDID = device.UDID
	payload.CommandPayload.Command.RequestType = requestType
	payload.CommandPayload.Command.Pin = pin
	_, err = SendCommand(payload)
	if err != nil {
		return errors.Wrap(err, "EraseLockDevice:SendCommand")
	}

	return nil
}

func escrowPin(device types.Device, pin string) error {
	escrowURL := utils.EscrowURL()
	if escrowURL == "" {
		DebugLogger(
			LogHolder{
				DeviceUDID:   device.UDID,
				DeviceSerial: device.SerialNumber,
				Message:      "No Escrow URL set, returning early",
			},
		)
		return nil
	}
	endpoint, err := url.Parse(escrowURL)
	if err != nil {
		return errors.Wrap(err, "escrowPin:URL")
	}

	urlString := strings.Trim(endpoint.String(), "/")
	urlString += "/"

	var payload types.EscrowPayload

	payload.Serial = device.SerialNumber
	payload.Pin = pin
	payload.Username = "mdmdirector"
	payload.SecretType = "unlock_pin"

	// log.Debug(payload)

	encoded, err := form.EncodeToValues(payload)
	if err != nil {
		return errors.Wrap(err, "escrowPin")
	}

	log.Debugf("Escrowing %v to %v", device.UDID, urlString)
	response, err := http.PostForm(urlString, encoded)

	if err != nil {
		return errors.Wrap(err, "escrowPin")
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)

	if err != nil {
		return errors.Wrap(err, "escrowPin:"+string(body))
	}

	return nil
}

func generatePin(device types.Device) (string, error) {
	// Look for an existing unlock pin generated within the last 30 mins
	var unlockPinModel types.UnlockPin
	var savedUnlockPin types.UnlockPin

	if device.UnlockPin != "" {
		return device.UnlockPin, nil
	}
	thirtyMinsAgo := time.Now().Add(-30 * time.Minute)

	if utils.DebugMode() {
		thirtyMinsAgo = time.Now().Add(-5 * time.Minute)
	}

	if err := db.DB.Model(&unlockPinModel).Where("unlock_pins.pin_set > ? AND unlock_pins.device_ud_id = ?", thirtyMinsAgo, device.UDID).Order("pin_set DESC").First(&unlockPinModel).Scan(&savedUnlockPin).Error; err != nil {
		if intErrors.Is(err, gorm.ErrRecordNotFound) {
			log.Debug("Pin was created more than 30 mins ago")
			out := ""
			for i := 0; i < 6; i++ {
				result, err := rand.Int(rand.Reader, big.NewInt(10))
				if err != nil {
					return "", errors.Wrap(err, "generatePin")
				}
				out += result.String()
			}
			var newUnlockPin types.UnlockPin
			newUnlockPin.DeviceUDID = device.UDID
			newUnlockPin.PinSet = time.Now()
			newUnlockPin.UnlockPin = out
			err = db.DB.Create(&newUnlockPin).Error
			if err != nil {
				return "", errors.Wrap(err, "generatePin:SavePin")
			}

			return out, nil
		}
	}
	// Found a saved one
	return savedUnlockPin.UnlockPin, nil
}
