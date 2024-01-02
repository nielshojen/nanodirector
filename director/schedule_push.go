package director

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/nielshojen/nanodirector/db"
	"github.com/nielshojen/nanodirector/types"
	"github.com/nielshojen/nanodirector/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/vmihailenco/taskq/v3"
)

func ScheduledCheckin(pushQueue taskq.Queue, onceIn time.Duration) {

	var task = taskq.RegisterTask(&taskq.TaskOptions{
		Name: "push",
		Handler: func(uuid string) error {
			err := PushDevice(uuid)
			if err != nil {
				ErrorLogger(LogHolder{Message: err.Error()})
			}
			return nil
		},
	})

	counter := 0
	for {
		if !DevicesFetchedFromMDM {
			time.Sleep(30 * time.Second)
			log.Info("Devices are still being fetched from MicroMDM")
			counter++
			if counter > 10 {
				break
			}
		} else {
			break
		}
	}

	var wg sync.WaitGroup
	sem := make(chan int, 1)

	fn := func(sem chan int, wg *sync.WaitGroup) {
		defer wg.Done()
		log.Info("Running scheduled checkin")
		err := processScheduledCheckin(pushQueue, task, onceIn)
		if err != nil {
			ErrorLogger(LogHolder{Message: err.Error()})
		}
		<-sem
	}

	for {
		sem <- 1
		wg.Add(1)
		go fn(sem, &wg)
	}
}

func ProcessScheduledCheckinQueue(pushQueue taskq.Queue) {
	ctx := context.Background()
	p := pushQueue.Consumer()
	DebugLogger(LogHolder{Message: "Processing items from scheduled checkin Queue"})
	err := p.Start(ctx)
	if err != nil {
		msg := fmt.Errorf("Starting consumer: %v", err.Error())
		ErrorLogger(LogHolder{Message: msg.Error()})
	}
}

func processScheduledCheckin(pushQueue taskq.Queue, task *taskq.Task, onceIn time.Duration) error {
	if utils.DebugMode() {
		DebugLogger(LogHolder{Message: "Processing scheduledCheckin in debug mode"})
	}

	err := pushAll(pushQueue, task, onceIn)
	if err != nil {
		return errors.Wrap(err, "processScheduledCheckin::pushAll")
	}

	var certificates []types.Certificate

	err = db.DB.Unscoped().Model(&certificates).Where("device_ud_id is NULL").Delete(&types.Certificate{}).Error
	if err != nil {
		return errors.Wrap(err, "processScheduledCheckin::CleanupNullCertificates")
	}

	var profileLists []types.ProfileList

	err = db.DB.Unscoped().Model(&profileLists).Where("device_ud_id is NULL").Delete(&types.ProfileList{}).Error
	if err != nil {
		return errors.Wrap(err, "processScheduledCheckin::CleanupNullProfileLists")
	}

	thirtyMinsAgo := time.Now().Add(-30 * time.Minute)
	err = db.DB.Where("unlock_pins.pin_set < ?", thirtyMinsAgo).Delete(&types.UnlockPin{}).Error
	if err != nil {
		return errors.Wrap(err, "processScheduledCheckin::DeleteRandomUnlockPins")
	}

	var device types.Device
	err = db.DB.Model(&device).Where("`erase` = ? AND `lock` = ? AND `unlock_pin` != ?", false, false, "").Update("unlock_pin", "").Error
	if err != nil {
		return errors.Wrap(err, "processScheduledCheckin::ResetFixedPin")
	}

	return nil
}

func deviceChunkSlice(slice []types.Device, chunkSize int) [][]types.Device {
	var chunks [][]types.Device
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		// necessary check to avoid slicing beyond
		// slice capacity
		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

func pushAll(pushQueue taskq.Queue, task *taskq.Task, onceIn time.Duration) error {
	var devices []types.Device
	var dbDevices []types.Device

	DelaySeconds := getDelay()

	err := db.DB.Find(&dbDevices).Scan(&dbDevices).Error
	if err != nil {
		return errors.Wrap(err, "PushAll: Scan devices")
	}

	for i := range dbDevices {
		device := dbDevices[i]
		needsPush := deviceNeedsPush(device)

		if needsPush {
			InfoLogger(LogHolder{DeviceUDID: device.UDID, DeviceSerial: device.SerialNumber, Message: "Adding Device to push list"})
			devices = append(devices, device)
		}

	}

	DebugLogger(LogHolder{
		Message: "Pushing to all in debug mode",
	})

	counter := 0
	total := 0
	devicesPerSecond := float64(len(devices)) / float64((DelaySeconds - 1))
	DebugLogger(LogHolder{Message: "Processed devices per 0.5 seconds", Metric: strconv.Itoa(int(devicesPerSecond))})

	devicesPerMinute := int(math.Ceil(float64(len(devices)) / 60))
	InfoLogger(LogHolder{Message: fmt.Sprintf("%d will be processed each minute", devicesPerMinute)})
	deviceChunks := deviceChunkSlice(devices, devicesPerMinute)
	InfoLogger(LogHolder{Message: fmt.Sprintf("%d chunks of %d devices each will be processed", len(deviceChunks), devicesPerMinute)})
	ctx := context.Background()
	msgTxt := fmt.Sprintf("commands will only be queued for an individual device every %s at maximum", onceIn)
	InfoLogger(LogHolder{Message: msgTxt})
	for i := range deviceChunks {
		for j := range deviceChunks[i] {
			device := deviceChunks[i][j]
			if float64(counter) >= devicesPerSecond {
				DebugLogger(LogHolder{Message: "Sleeping due to having processed devices", Metric: strconv.Itoa(total)})
				time.Sleep(500 * time.Millisecond)
				counter = 0
			}
			DebugLogger(LogHolder{Message: "pushAll processed", Metric: strconv.Itoa(counter)})

			msg := task.WithArgs(ctx, device.UDID)

			msg.OnceInPeriod(onceIn)
			err := pushQueue.Add(msg)
			switch {
			case errors.Is(msg.Err, taskq.ErrDuplicate):
				// handle duplicate task
				DebugLogger(LogHolder{DeviceUDID: device.UDID, DeviceSerial: device.SerialNumber, Message: msg.Err.Error()})
			case err != nil:
				ErrorLogger(LogHolder{DeviceUDID: device.UDID, DeviceSerial: device.SerialNumber, Message: err.Error()})
			case msg.Err != nil:
				// handle duplicate task
				ErrorLogger(LogHolder{DeviceUDID: device.UDID, DeviceSerial: device.SerialNumber, Message: msg.Err.Error()})
			}

			counter++
			total++
		}
		// Wait 1 minute before processing the next chunk of devices
		msg := fmt.Sprintf("%d/%d devices processed", (i+1)*devicesPerMinute, len(devices))
		InfoLogger(LogHolder{Message: msg})
		InfoLogger(LogHolder{Message: "Sleeping 1 minute before processing next chunk of devices"})
		time.Sleep(time.Minute * 1)
	}
	InfoLogger(LogHolder{Message: "Completed scheduling pushes", Metric: strconv.Itoa(len(devices))})
	return nil
}

func deviceNeedsPush(device types.Device) bool {
	now := time.Now()
	oneDayAgo := time.Now().Add(-24 * time.Hour)

	InfoLogger(LogHolder{DeviceUDID: device.UDID, DeviceSerial: device.SerialNumber, Message: "Considering device for scheduled push"})

	if now.Before(device.NextPush) && !device.NextPush.IsZero() {
		InfoLogger(LogHolder{DeviceUDID: device.UDID, DeviceSerial: device.SerialNumber, Message: "Not Pushing. Next push is in metric", Metric: device.NextPush.String()})
		return false
	}

	if device.LastCertificateList.IsZero() || device.LastProfileList.IsZero() || device.LastSecurityInfo.IsZero() || device.LastDeviceInfo.IsZero() {
		InfoLogger(LogHolder{DeviceUDID: device.UDID, DeviceSerial: device.SerialNumber, Message: "One or more of the info commands hasn't ever been received"})
		return true
	}

	// We've not had all of the info payloads within the last day
	if (device.LastCertificateList.Before(oneDayAgo) || device.LastProfileList.Before(oneDayAgo) || device.LastSecurityInfo.Before(oneDayAgo) || device.LastDeviceInfo.Before(oneDayAgo)) && (!device.LastCertificateList.IsZero() && !device.LastProfileList.IsZero() && !device.LastSecurityInfo.IsZero() && !device.LastDeviceInfo.IsZero()) {
		InfoLogger(LogHolder{DeviceUDID: device.UDID, DeviceSerial: device.SerialNumber, Message: "Have not received all of the info commands within the last six hours."})
		return true
	}

	// If it's been updated within the last three hours, try to push again as it might still be online
	// if device.LastCheckedIn.After(threeHoursAgo) {
	// 	InfoLogger(LogHolder{DeviceUDID: device.UDID, DeviceSerial: device.SerialNumber, Message: "Checked in more than three hours ago"})
	// 	if now.Before(device.NextPush) {
	// 		InfoLogger(LogHolder{DeviceUDID: device.UDID, DeviceSerial: device.SerialNumber, Message: "Not Pushing. Next push is in metric", Metric: device.NextPush.String()})
	// 		return false
	// 	}
	// }

	return true
}

func PushDevice(udid string) error {
	device := types.Device{UDID: udid}
	InfoLogger(LogHolder{DeviceUDID: device.UDID, Message: "Sending push to device"})
	now := time.Now()
	InfoLogger(LogHolder{DeviceUDID: device.UDID, DeviceSerial: device.SerialNumber, Message: "Performing device push"})

	// This will ensure the apns push is expired when we push again - this helps prevent a flood of responses from the device whilst ensuring there is always a push waiting for the device when it comes online.
	retry := now.Add(time.Minute * time.Duration(int64(utils.OnceIn())))
	retryUnix := retry.Unix()

	endpoint, err := url.Parse(utils.NanoURL())
	if err != nil {
		return errors.Wrap(err, "PushDevice")
	}

	endpoint.Path = path.Join(endpoint.Path, "v1/push", device.UDID)
	queryString := endpoint.Query()
	queryString.Set("expiration", strconv.FormatInt(retryUnix, 10))
	endpoint.RawQuery = queryString.Encode()
	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		return errors.Wrap(err, "PushDevice")
	}
	req.SetBasicAuth("nanomdm", utils.NanoAPIKey())

	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "PushDevice")
	}

	err = resp.Body.Close()
	if err != nil {
		return errors.Wrap(err, "PushDevice")
	}

	InfoLogger(LogHolder{DeviceUDID: device.UDID, Message: "Sent push to device"})

	return nil
}
