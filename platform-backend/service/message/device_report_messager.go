package message

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	db "platform-backend/db"
	"platform-backend/dto"
	"platform-backend/dto/enum"
	"platform-backend/models"
	"platform-backend/repository"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gorm.io/gorm"
)

// DeviceKeepaliveMessageHandler 消息处理失败 => 1.重试 2.移交到专门的失败处理去重试 3.丢弃
func DeviceKeepaliveMessageHandler(client mqtt.Client, msg mqtt.Message) {
	var data dto.KeepAliveData
	repo := repository.NewDeviceRepository(db.MysqlDB.DB())
	ctx := context.Background()
	// fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	err := json.Unmarshal(msg.Payload(), &data)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v", err)
		return
	}

	if len(data.Mac) == 0 || len(data.UniqueId) == 0 {
		return
	}

	device, err := repo.FindOne(ctx, map[string]interface{}{"unique_id": data.UniqueId})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Printf("Error repo FindOne: %v", err)
		return
	}
	now := time.Now().UTC()

	if errors.Is(err, gorm.ErrRecordNotFound) {
		device = &models.Device{
			UniqueId:   data.UniqueId,
			DeviceName: data.Name,
			Mac:        data.Mac,
			Ip:         data.IP,
			OsType:     data.OsType,
			OnLineTime: &now,
			DeviceType: enum.DeviceTypePhysical,
			Status:     enum.DeviceStatusOnLine,
		}
		err := repo.Create(ctx, device)
		if err != nil {
			fmt.Printf("Error repo Create: %v", err)
			return
		}

	} else {
		if device.OpStatus == enum.DeviceOpStatusPowerOn {
			device.OpStatus = enum.DeviceOpStatusNormal
		}
		device.DeviceName = data.Name
		device.Mac = data.Mac
		device.Ip = data.IP
		device.OsType = data.OsType
		device.DeviceType = enum.DeviceTypePhysical
		device.Status = enum.DeviceStatusOnLine
		device.OnLineTime = &now
		err := repo.Update(ctx, device)
		if err != nil {
			fmt.Printf("Error repo Update: %v", err)
			return
		}
	}

	device, err = repo.FindOne(ctx, map[string]interface{}{"unique_id": data.UniqueId})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Printf("Error repo FindOne: %v", err)
		return
	}

	for _, item := range data.VirHosts {
		deviceType := enum.DeviceTypeKvm
		deviceStatus := enum.DeviceStatusOnLine
		if item.VirType == "kvm" {
			deviceType = enum.DeviceTypeKvm
		}
		if item.State == "running" {
			deviceStatus = enum.DeviceStatusOnLine
		} else {
			deviceStatus = enum.DeviceStatusOffline
		}

		host, err := repo.FindOne(ctx, map[string]interface{}{
			"device_name":    item.Name,
			"device_type":    deviceType,
			"host_device_id": device.ID})
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Printf("Error repo FindOne: %v", err)
			return
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			virhost := &models.Device{
				UniqueId:     data.UniqueId + "_" + item.Name,
				DeviceName:   item.Name,
				Mac:          "",
				Ip:           "",
				OsType:       "",
				OnLineTime:   &now,
				DeviceType:   deviceType,
				Status:       deviceStatus,
				HostDeviceId: device.ID,
			}
			err := repo.Create(ctx, virhost)
			if err != nil {
				fmt.Printf("Error repo Create: %v", err)
				return
			}
		} else {
			if deviceStatus == enum.DeviceStatusOnLine {
				if host.OpStatus == enum.DeviceOpStatusPowerOn {
					host.OpStatus = enum.DeviceOpStatusNormal
				}
			}
			if deviceStatus == enum.DeviceStatusOffline {
				if host.OpStatus == enum.DeviceOpStatusPowerOff {
					host.OpStatus = enum.DeviceOpStatusNormal
				}
			}
			host.OnLineTime = &now
			host.Status = deviceStatus
			err := repo.Update(ctx, host)
			if err != nil {
				fmt.Printf("Error repo Update: %v", err)
				return
			}
		}
	}
}
