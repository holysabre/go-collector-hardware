package disk

type Smartctl struct {
	Devices []Device `json:"devices"`
}

type Device struct {
	Name     string `json:"name"`
	InfoName string `json:"info_name"`
	Type     string `json:"type"`
	Protocol string `json:"protocol"`
}

type DiskInfo struct {
	ModelName    string       `json:"model_name"`
	SmartStatus  SmartStatus  `json:"smart_status"`
	UserCapacity UserCapacity `json:"user_capacity"`
	Temperature  Temperature  `json:"temperature"`
	PowerOnTime  PowerOnTime  `json:"power_on_time"`
	SerialNumber string       `json:"serial_number"`
	RotationRate int64        `json:"rotation_rate"`
	Device       Device       `json:"device"`
	SetaVersion  SetaVersion  `json:"seta_version"`
	ModelType    string
}

type SetaVersion struct {
	String string `json:"string"`
	Value  int64  `json:"value"`
}

type SmartStatus struct {
	Passed bool `json:"passed"`
}

type UserCapacity struct {
	Blocks int64 `json:"blocks"`
	Bytes  int64 `json:"bytes"`
}

type Temperature struct {
	Current int8 `json:"current"`
}

type PowerOnTime struct {
	Hours int64 `json:"hours"`
}
