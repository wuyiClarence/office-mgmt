package config

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	mylog "platform-mdns/utils/log"
)

type Config struct {
	AppConfig  *appConfig  `json:"app_config" mapstructure:"App"`
	MdnsConfig *MdnsConfig `json:"mdns_config" mapstructure:"MdnsConfig"`
	MqttConfig *MqttConfig `json:"mqtt_config" mapstructure:"MqttConfig"`
}

type appConfig struct {
	Name      string `json:"name" mapstructure:"Name"`
	Address   string `json:"address" mapstructure:"Address"`
	Mod       string `json:"mod" mapstructure:"Mod"`
	LogExpire int    `json:"log_expire" mapstructure:"LogExpire"`
}

type MdnsConfig struct {
	Domain      string `json:"domain" mapstructure:"Domain"`
	Host        string `json:"host" mapstructure:"Host"`
	Key         string `json:"key" mapstructure:"Key"`
	ServiceName string `json:"service_name" mapstructure:"ServiceName"`
	Port        int    `json:"port" mapstructure:"Port"`
}

type MqttConfig struct {
	Broker   string `json:"broker" mapstructure:"Broker"`
	UserName string `json:"user_name" mapstructure:"UserName"`
	Password string `json:"password" mapstructure:"Password"`
	Port     int    `json:"port" mapstructure:"Port"`
}

var myViper *viper.Viper

var MyConfig *Config

func InitConfig() error {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config"
	}

	myViper = viper.New()
	myViper.AddConfigPath(configPath)
	myViper.SetConfigName("config")
	myViper.SetConfigType("yaml")

	err := myViper.ReadInConfig()
	if err != nil {
		return err
	}

	if err = myViper.Unmarshal(&MyConfig); err != nil {
		return err
	}

	return nil
}

func DynamicReloadConfig() {
	myViper.WatchConfig()
	myViper.OnConfigChange(func(event fsnotify.Event) {
		if err := myViper.Unmarshal(&MyConfig); err != nil {
			_, _ = fmt.Fprintf(mylog.SystemLogger, "更新配置失败：%s", err.Error())
		} else {
			_, _ = fmt.Fprintf(mylog.SystemLogger, "更新配置")
		}
	})
}
