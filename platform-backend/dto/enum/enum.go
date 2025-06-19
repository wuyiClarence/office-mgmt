package enum

const SuperAdminUserAccName = "SuperAdmin"
const SuperAdminDefaultPassword = "SuperAdmin!123456"
const UserDefaultPassword = "Password123456!"

type ActionType string

const (
	ActionTypeOn  ActionType = "power_on"
	ActionTypeOff ActionType = "power_off"
)

type TimeConditionType int

const (
	TypeOnce      TimeConditionType = 1
	TypeEveryDay  TimeConditionType = 2
	TypeDayOfWeek TimeConditionType = 3
)

func (x ActionType) String() string {
	return string(x)
}

type DeviceType int8

const (
	DeviceTypePhysical DeviceType = 1
	DeviceTypeKvm      DeviceType = 2
)

type DeviceStatus int8

const (
	DeviceStatusOffline DeviceStatus = 0
	DeviceStatusOnLine  DeviceStatus = 1
)

type DeviceOpStatus int8

const (
	DeviceOpStatusNormal   DeviceOpStatus = 0
	DeviceOpStatusPowerOff DeviceOpStatus = 1
	DeviceOpStatusPowerOn  DeviceOpStatus = 2
	DeviceOpStatusReboot   DeviceOpStatus = 3
)

type AssociateType int8

const (
	AssociateTypeDevice      AssociateType = 1
	AssociateTypeDeviceGroup AssociateType = 2
)

type PolicyStatus int8

const (
	PolicyStatusDisable PolicyStatus = 0
	PolicyStatusEnable  PolicyStatus = 1
)

type UserStatus int8

const (
	UserStatusOk      UserStatus = 0
	UserStatusDeleted UserStatus = 1
)

type MqttQos byte

const (
	QoS_0 MqttQos = 0 // 最多一次
	QoS_1 MqttQos = 1 // 至少一次
	QoS_2 MqttQos = 2 // 只有一次
)

type ResourceType string

const (
	ResourceTypeRole        ResourceType = "role"
	ResourceTypeDevice      ResourceType = "device"
	ResourceTypeDeviceGroup ResourceType = "device_group"
	ResourceTypePolicy      ResourceType = "policy"
	ResourceTypeUser        ResourceType = "user"
)

type PermissionKey string

const (
	PermissionKeyView           PermissionKey = "view"
	PermissionKeyEdit           PermissionKey = "edit"
	PermissionKeyDelete         PermissionKey = "delete"
	PermissionKeyPowerOn        PermissionKey = "poweron"
	PermissionKeyPowerOff       PermissionKey = "poweroff"
	PermissionKeyPermissionMgmt PermissionKey = "permissionmgmt"
)
