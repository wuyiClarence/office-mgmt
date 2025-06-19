package message

import (
	"sync/atomic"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"platform-backend/dto/enum"
	mymqtt "platform-backend/mqtt"
)

type Handler struct {
	MsgHandler mqtt.MessageHandler
	Qos        enum.MqttQos
}

func Start() error {
	err := doSubscribe()
	if err != nil {
		return err
	}

	//监控并重新订阅
	go monitorSubscribe()

	return nil
}

func monitorSubscribe() {
	ticker := time.NewTicker(time.Duration(10) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if atomic.LoadInt32(&mymqtt.ReconnectFlag) != 0 {
			continue
		}

		_ = doSubscribe()
	}

}

func doSubscribe() error {
	mqttClient, err := mymqtt.GetMqttClient()
	if err != nil {
		return err
	}

	var TopicsHandlerMap = map[string]Handler{
		"/keepalive": {
			MsgHandler: DeviceKeepaliveMessageHandler,
			Qos:        enum.QoS_0,
		},
	}

	for topic, handler := range TopicsHandlerMap {
		if err = mqttClient.Subscribe(topic, byte(handler.Qos), handler.MsgHandler); err != nil {
			return err
		}
	}

	atomic.StoreInt32(&mymqtt.ReconnectFlag, 1)

	return nil
}
