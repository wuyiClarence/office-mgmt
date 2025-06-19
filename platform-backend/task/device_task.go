package task

import (
	"context"
	"fmt"
	db "platform-backend/db"
	"platform-backend/dto/enum"
	"platform-backend/models"
	"platform-backend/repository"
	"platform-backend/service/deviceoprator"
	"platform-backend/utils/log"
	"time"
)

type BaseTaskImpl struct {
	Policy *models.Policy
}

func (b *BaseTaskImpl) GetPolicy() *models.Policy {
	return b.Policy
}

type PowerOnTask struct {
	BaseTaskImpl
}

func CheckTimeRange(p models.Policy) bool {
	now := time.Now().UTC()

	if p.StartDate != nil {
		diff := now.Sub(*p.StartDate)

		if diff < 0 {
			_, _ = fmt.Fprintf(log.SystemLogger, "policy:%s not in timerange,starttime: %v\n", p.PolicyName, p.StartDate)
			return false
		}
	}

	if p.EndDate != nil {
		diff := now.Sub(*p.EndDate)

		if diff > 5*time.Millisecond {
			_, _ = fmt.Fprintf(log.SystemLogger, "policy:%s not in timerange,starttime: %v\n", p.PolicyName, p.EndDate)
			return false
		}
	}
	return true
}

func GetPolicyDeviceIds(p models.Policy) ([]*models.Device, error) {
	var deviceIds []int64
	var err error
	devices := make([]*models.Device, 0)

	deviceRep := repository.NewDeviceRepository(db.MysqlDB.DB())

	ctx := context.Background()

	if p.AssociateType == enum.AssociateTypeDeviceGroup {
		policyDeviceGroupRep := repository.NewPolicyDeviceGroupRelRepository(db.MysqlDB.DB())
		deviceGroupRelRep := repository.NewDeviceDeviceGroupRelRepository(db.MysqlDB.DB())

		deviceGroupIds, err := policyDeviceGroupRep.GetDeviceGroupIDs(ctx, p.ID)
		if err != nil {
			_, _ = fmt.Fprintf(log.SystemLogger, "%s policy:%s get deviceGroupIds error %v\n", time.Now().Format("2006-01-02 15:04:05"), p.PolicyName, err)
			return nil, err
		}
		deviceIds, err = deviceGroupRelRep.GetDeviceIdsByDeviceGroupIds(ctx, deviceGroupIds)
		if err != nil {
			_, _ = fmt.Fprintf(log.SystemLogger, "%s policy:%s get deviceIds error %v\n", time.Now().Format("2006-01-02 15:04:05"), p.PolicyName, err)
			return nil, err
		}
	}
	if p.AssociateType == enum.AssociateTypeDevice {
		policyDeviceRep := repository.NewPolicyDeviceRelRepository(db.MysqlDB.DB())
		deviceIds, err = policyDeviceRep.GetDeviceIDs(ctx, p.ID)
		if err != nil {
			_, _ = fmt.Fprintf(log.SystemLogger, "%s policy:%s get deviceid error %v\n", time.Now().Format("2006-01-02 15:04:05"), p.PolicyName, err)
			return nil, err
		}
	}

	if len(deviceIds) == 0 {
		_, _ = fmt.Fprintf(log.SystemLogger, "%s policy:%s no devices\n", time.Now().Format("2006-01-02 15:04:05"), p.PolicyName)
		return devices, nil
	}

	devices, err = deviceRep.GetDevicesByIds(ctx, deviceIds)
	if err != nil {
		_, _ = fmt.Fprintf(log.SystemLogger, "%s policy:%s get device error %v\n", time.Now().Format("2006-01-02 15:04:05"), p.PolicyName, err)
		return nil, err
	}

	return devices, nil
}

func (p *PowerOnTask) Exec() {
	intime := CheckTimeRange(*p.Policy)
	if !intime {
		return
	}

	devices, err := GetPolicyDeviceIds(*p.Policy)
	if err != nil {
		return
	}
	for _, device := range devices {
		err := deviceoprator.PowerOn(context.Background(), device)
		if err != nil {
			_, _ = fmt.Fprintf(log.SystemLogger, "%s send MQTT to %s,err %v\n", time.Now().Format("2006-01-02 15:04:05"), device.Mac, err)
		} else {
			_, _ = fmt.Fprintf(log.SystemLogger, "%s send MQTT to %s\n", time.Now().Format("2006-01-02 15:04:05"), device.Mac)
		}
	}
	_, _ = fmt.Fprintf(log.SystemLogger, "======on=======:%s %s\n", p.Policy.PolicyName, time.Now().Format("2006-01-02 15:04:05"))
}

type PowerOffTask struct {
	BaseTaskImpl
}

func (p *PowerOffTask) Exec() {

	intime := CheckTimeRange(*p.Policy)
	if !intime {
		return
	}
	devices, err := GetPolicyDeviceIds(*p.Policy)
	if err != nil {
		return
	}

	for _, device := range devices {
		err := deviceoprator.PowerOff(context.Background(), device)

		if err != nil {
			_, _ = fmt.Fprintf(log.SystemLogger, "%s send MQTT to %s, err:%v\n", time.Now().Format("2006-01-02 15:04:05"), device.Mac, err)
		} else {
			_, _ = fmt.Fprintf(log.SystemLogger, "%s send MQTT to %s\n", time.Now().Format("2006-01-02 15:04:05"), device.Mac)
		}
	}
	_, _ = fmt.Fprintf(log.SystemLogger, "======off=======:%s %s\n", p.Policy.PolicyName, time.Now().Format("2006-01-02 15:04:05"))
}

type TimeOnlineTask struct {
}

func (p *TimeOnlineTask) Exec() {
	now := time.Now().UTC()

	ctx := context.Background()

	fmt.Printf("timeonline %s", time.Now().Format("2006-01-02 15:04:05"))
	deviceRep := repository.NewDeviceRepository(db.MysqlDB.DB())

	allDevices, err := deviceRep.FindAll(ctx, map[string]interface{}{
		"status": 1,
	}, false)
	if err != nil {
		return
	}

	for _, device := range allDevices {
		if device.OnLineTime != nil {
			diff := now.Sub(*device.OnLineTime)
			if diff > 1*time.Minute {
				if device.OpStatus == enum.DeviceOpStatusPowerOff {
					device.OpStatus = enum.DeviceOpStatusNormal
				}
				device.Status = enum.DeviceStatusOffline
				err := deviceRep.Update(ctx, &device)
				if err != nil {
					fmt.Printf("Error repo Update: %v", err)
					return
				}
			} else {
				// 在线
				if device.OpStatus == enum.DeviceOpStatusPowerOn {
					device.OpStatus = enum.DeviceOpStatusNormal
					err := deviceRep.Update(ctx, &device)
					if err != nil {
						fmt.Printf("Error repo Update: %v", err)
						return
					}
				}
			}
		}

	}

	powerOffDevices, err := deviceRep.FindAll(ctx, map[string]interface{}{
		"op_status": enum.DeviceOpStatusPowerOff,
	}, false)
	if err != nil {
		return
	}
	for _, powerOffDevice := range powerOffDevices {
		deviceoprator.PowerOff(context.Background(), &powerOffDevice)
	}

	powerOnDevices, err := deviceRep.FindAll(ctx, map[string]interface{}{
		"op_status": enum.DeviceOpStatusPowerOn,
	}, false)
	if err != nil {
		return
	}
	for _, powerOnDevice := range powerOnDevices {
		deviceoprator.PowerOn(context.Background(), &powerOnDevice)
	}
}
