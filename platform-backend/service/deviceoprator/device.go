package deviceoprator

import (
	"context"
	"fmt"
	db "platform-backend/db"
	"platform-backend/dto/enum"
	"platform-backend/models"
	mymqtt "platform-backend/mqtt"
	"platform-backend/repository"
	"platform-backend/utils/log"
	"time"
)

func PowerOn(ctx context.Context, device *models.Device) error {
	mqttClient, err := mymqtt.GetMqttClient()
	if err != nil {
		return err
	}

	deviceRepo := repository.NewDeviceRepository(db.MysqlDB.DB())
	device.OpStatus = enum.DeviceOpStatusPowerOn

	err = deviceRepo.Update(ctx, device)
	if err != nil {
		return err
	}

	if device.DeviceType == enum.DeviceTypePhysical {
		// err := sendWOLPacket(device.Mac)
		// if err != nil {
		// }
		err = mqttClient.PlatformSendPowerOn(device.Mac)
		if err != nil {
			return err
		}
		wakehosts, err := deviceRepo.GetPhysicalDevices(ctx)
		if err != nil {
			return err
		}
		for _, wakehost := range wakehosts {
			if wakehost.Mac != device.Mac {
				err = mqttClient.SendPowerOn(wakehost, device.Mac)
				if err != nil {
					continue
				}
			}
		}
	}
	if device.DeviceType == enum.DeviceTypeKvm {
		err = mqttClient.SendVirHostPowerOn(device)
	}
	if err != nil {
		_, _ = fmt.Fprintf(log.SystemLogger, "%s send MQTT to %s,err %v\n", time.Now().Format("2006-01-02 15:04:05"), device.Mac, err)
	} else {
		_, _ = fmt.Fprintf(log.SystemLogger, "%s send MQTT to %s\n", time.Now().Format("2006-01-02 15:04:05"), device.Mac)
	}
	return err
}

func PowerOff(ctx context.Context, device *models.Device) error {
	mqttClient, err := mymqtt.GetMqttClient()
	if err != nil {
		return err
	}

	deviceRepo := repository.NewDeviceRepository(db.MysqlDB.DB())
	device.OpStatus = enum.DeviceOpStatusPowerOff

	err = deviceRepo.Update(ctx, device)
	if err != nil {
		return err
	}

	if device.DeviceType == enum.DeviceTypePhysical {
		err = mqttClient.SendPowerOff(device)
	}
	if device.DeviceType == enum.DeviceTypeKvm {
		err = mqttClient.SendVirHostPowerOff(device)
	}
	if err != nil {
		_, _ = fmt.Fprintf(log.SystemLogger, "%s send MQTT to %s,err %v\n", time.Now().Format("2006-01-02 15:04:05"), device.Mac, err)
	} else {
		_, _ = fmt.Fprintf(log.SystemLogger, "%s send MQTT to %s\n", time.Now().Format("2006-01-02 15:04:05"), device.Mac)
	}

	return err
}

// // sendWOLPacket 发送 Wake-on-LAN 魔术包
// func sendWOLPacket(macAddr string) error {
// 	mac, err := net.ParseMAC(macAddr)
// 	if err != nil {
// 		return errors.New("无效的 MAC 地址")
// 	}

// 	// WOL 魔术包由 6 个 0xFF 和 16 组 MAC 地址组成
// 	packet := make([]byte, 102)
// 	for i := 0; i < 6; i++ {
// 		packet[i] = 0xFF
// 	}
// 	for i := 6; i < len(packet); i += len(mac) {
// 		copy(packet[i:i+len(mac)], mac)
// 	}

// 	// 发送 UDP 数据包到广播地址
// 	conn, err := net.Dial("udp", "255.255.255.255:9")
// 	if err != nil {
// 		return err
// 	}
// 	defer conn.Close()

// 	_, err = conn.Write(packet)
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Println("WOL 魔术包已发送至:", macAddr)
// 	return nil
// }
