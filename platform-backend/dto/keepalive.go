package dto

// 主结构体
type KeepAliveData struct {
	UniqueId string     `json:"uniqueId"`
	Mac      string     `json:"mac"`
	OsType   string     `json:"ostype"`
	Name     string     `json:"name"`
	IP       string     `json:"ip"`
	CurTime  string     `json:"curTime"`
	VirHosts []VirHosts `json:"virhost"`
}

// 虚拟主机信息结构体
type VirHosts struct {
	Name    string `json:"name"`
	VirType string `json:"virtype"`
	State   string `json:"state"`
}

type ShutdownMsg struct {
	Mac string `json:"mac"`
}
type VirHostShutdownMsg struct {
	Mac      string     `json:"mac"`
	VirHosts []VirHosts `json:"virhost"`
}

type VirHostWakeOnMsg struct {
	Mac      string     `json:"mac"`
	VirHosts []VirHosts `json:"virhost"`
}

type WakeOnMsg struct {
	Mac     string `json:"mac"`
	WakeMac string `json:"wakemac"`
}
