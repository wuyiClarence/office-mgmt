package mqtt

import (
	"fmt"
	"sync/atomic"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"platform-backend/config"
	"platform-backend/utils"
	mylog "platform-backend/utils/log"

	"github.com/google/uuid"
)

var (
	ReconnectFlag int32

	instance *MqttClient
	once     utils.OnceV2
)

type MqttClient struct {
	client mqtt.Client
}

func GetMqttClient() (*MqttClient, error) {
	err := once.Do(func() error {
		var err error

		instance, err = newMqttClient()
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return instance, err
}

func newMqttClient() (*MqttClient, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.MyConfig.MqttConfig.Broker)
	ClientID := fmt.Sprintf("%s-%s", config.MyConfig.MqttConfig.ClientID, uuid.New().String()[:8])
	_, _ = fmt.Fprintf(mylog.SystemLogger, "MQTT ClientID %s", ClientID)
	opts.SetClientID(ClientID)
	opts.SetUsername(config.MyConfig.MqttConfig.UserName)
	opts.SetPassword(config.MyConfig.MqttConfig.Password)
	opts.SetKeepAlive(time.Duration(config.MyConfig.MqttConfig.KeepAlive) * time.Second)
	opts.SetConnectTimeout(time.Duration(config.MyConfig.MqttConfig.ConnectTimeOut) * time.Second)

	//重连
	opts.SetMaxReconnectInterval(time.Duration(config.MyConfig.MqttConfig.MaxReconnectInterval) * time.Second)
	opts.SetAutoReconnect(true)

	opts.OnConnect = func(c mqtt.Client) {
		_, _ = fmt.Fprintf(mylog.SystemLogger, "MQTT连接成功！！")

		atomic.StoreInt32(&ReconnectFlag, 0)
	}
	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		_, _ = fmt.Fprintf(mylog.SystemLogger, "MQTT连接丢失: %v\n", err)

		//主动重连
		_ = c.Connect()
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return &MqttClient{client: client}, nil
}

func (m *MqttClient) Publish(topic string, qos byte, retained bool, message string) error {
	token := m.client.Publish(topic, qos, retained, message)
	token.Wait()
	return token.Error()
}

func (m *MqttClient) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) error {
	token := m.client.Subscribe(topic, qos, callback)
	token.Wait()
	return token.Error()
}

func (m *MqttClient) Unsubscribe(topic string) error {
	token := m.client.Unsubscribe(topic)
	token.Wait()
	return token.Error()
}

func (m *MqttClient) Disconnect() {
	m.client.Disconnect(250)
}
