// config/config.go

package config

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	mylog "platform-backend/utils/log"
)

type Config struct {
	AppConfig       *appConfig   `json:"app_config" mapstructure:"App"`
	MysqlConfig     *mysqlConfig `json:"mysql_config" mapstructure:"Mysql"`
	AutoMigration   bool         `json:"auto_migration" mapstructure:"AutoMigration"`
	CronGapTime     int          `json:"cron_gap_time" mapstructure:"CronGapTime"`
	MqttConfig      *mqttConfig  `json:"mqtt_config" mapstructure:"MqttConfig"`
	PasswordAuthKey string       `json:"password_auth_key" mapstructure:"PasswordAuthKey"`
}

type appConfig struct {
	Name      string `json:"name" mapstructure:"Name"`
	Address   string `json:"address" mapstructure:"Address"`
	Mod       string `json:"mod" mapstructure:"Mod"`
	LogExpire int    `json:"log_expire" mapstructure:"LogExpire"`
}

type mysqlConfig struct {
	Dsn     string `json:"dsn" mapstructure:"Dsn"`
	MaxIdle int    `json:"max_idle" mapstructure:"MaxIdle"`
	MaxOpen int    `json:"max_open" mapstructure:"MaxOpen"`
	Name    string `json:"name" mapstructure:"Name"`
	Debug   bool   `json:"debug" mapstructure:"Debug"`
}

type mqttConfig struct {
	Broker               string `json:"broker" mapstructure:"Broker"`
	ClientID             string `json:"client_id" mapstructure:"ClientID"`
	UserName             string `json:"user_name" mapstructure:"UserName"`
	Password             string `json:"password" mapstructure:"Password"`
	KeepAlive            int    `json:"keep_alive" mapstructure:"KeepAlive"`
	ConnectTimeOut       int    `json:"connect_time_out" mapstructure:"ConnectTimeOut"`
	MaxReconnectInterval int    `json:"max_reconnect_interval" mapstructure:"MaxReconnectInterval"`
}

var myViper *viper.Viper

var MyConfig *Config

func InitConfig() error {

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config"
	}
	fmt.Printf("configPath: %s\n", configPath)
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
