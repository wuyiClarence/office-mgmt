package dto

type MdnsInfo struct {
	Broker   string `json:"mqtt_broker"`
	UserName string `json:"mqtt_username"`
	Password string `json:"mqtt_password"`
	Port     int    `json:"mqtt_port"`
}
