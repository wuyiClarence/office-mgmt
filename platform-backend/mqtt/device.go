package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	db "platform-backend/db"
	"platform-backend/dto"
	"platform-backend/dto/enum"
	"platform-backend/models"
	"platform-backend/repository"
)

func (m *MqttClient) SendVirHostPowerOff(device *models.Device) error {
	deviceRep := repository.NewDeviceRepository(db.MysqlDB.DB())

	ctx := context.Background()

	host, err := deviceRep.FindOne(ctx, map[string]interface{}{"id": device.HostDeviceId})
	if err != nil {
		return err
	}
	data := dto.VirHostShutdownMsg{
		Mac:      host.Mac,
		VirHosts: make([]dto.VirHosts, 0),
	}
	virType := "kvm"
	if device.DeviceType == enum.DeviceTypeKvm {
		virType = "kvm"
	}
	data.VirHosts = append(data.VirHosts, dto.VirHosts{
		Name:    device.DeviceName,
		VirType: virType,
	})

	topic := fmt.Sprintf("/%s/shutdownvirhost", host.UniqueId)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = m.Publish(topic, byte(enum.QoS_2), false, string(jsonData))

	if err != nil {
		return err
	}
	return nil
}

func (m *MqttClient) SendPowerOff(device *models.Device) error {
	uniqueId := device.UniqueId

	data := dto.ShutdownMsg{
		Mac: device.Mac,
	}

	topic := fmt.Sprintf("/%s/shutdownhost", uniqueId)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = m.Publish(topic, byte(enum.QoS_2), false, string(jsonData))
	if err != nil {
		return err
	}

	return nil
}

func (m *MqttClient) SendVirHostPowerOn(device *models.Device) error {
	deviceRep := repository.NewDeviceRepository(db.MysqlDB.DB())

	ctx := context.Background()

	host, err := deviceRep.FindOne(ctx, map[string]interface{}{"id": device.HostDeviceId})
	if err != nil {
		return err
	}
	data := dto.VirHostWakeOnMsg{
		Mac:      host.Mac,
		VirHosts: make([]dto.VirHosts, 0),
	}
	virType := "kvm"
	if device.DeviceType == enum.DeviceTypeKvm {
		virType = "kvm"
	}
	data.VirHosts = append(data.VirHosts, dto.VirHosts{
		Name:    device.DeviceName,
		VirType: virType,
	})

	topic := fmt.Sprintf("/%s/wakevirhost", host.UniqueId)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = m.Publish(topic, byte(enum.QoS_2), false, string(jsonData))

	if err != nil {
		return err
	}
	return nil
}
func (m *MqttClient) SendPowerOn(device *models.Device, wakemac string) error {
	uniqueId := device.UniqueId

	if device.DeviceType == enum.DeviceTypeKvm {
		err := m.SendVirHostPowerOn(device)
		if err != nil {
			return err
		}
	} else {
		data := dto.WakeOnMsg{
			Mac:     device.Mac,
			WakeMac: wakemac,
		}

		topic := fmt.Sprintf("/%s/wakehost", uniqueId)
		jsonData, err := json.Marshal(data)
		if err != nil {
			return err
		}
		err = m.Publish(topic, byte(enum.QoS_2), false, string(jsonData))
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *MqttClient) PlatformSendPowerOn(wakemac string) error {

	data := dto.WakeOnMsg{
		WakeMac: wakemac,
	}

	topic := fmt.Sprintf("/wakehost")
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = m.Publish(topic, byte(enum.QoS_2), false, string(jsonData))
	if err != nil {
		return err
	}

	return nil
}
