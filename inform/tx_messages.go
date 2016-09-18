package inform

type AdminMetadata struct {
	Id       string `json:"_id"`
	Language string `json:"lang"`
	Username string `json:"name"`
	Password string `json:"x_password"`
}

// TODO: Convert string time to time.Time
type CommandMessage struct {
	Metadata   *AdminMetadata `json:"_admin"`
	Id         string         `json:"_id"`
	Type       string         `json:"_type"`    // cmd
	Command    string         `json:"cmd"`      // mfi-output
	DateTime   string         `json:"datetime"` // 2016-07-28T01:17:55Z
	DeviceId   string         `json:"device_id"`
	MacAddress string         `json:"mac"`
	Model      string         `json:"model"`
	OffVoltage int            `json:"off_volt"`
	Port       int            `json:"port"`
	MessageId  string         `json:"sId"` // ??
	ServerTime string         `json:"server_time_in_utc"`
	Time       string         `json:"time"`
	Timer      int            `json:"timer"`
	Value      int            `json:"val"`
	Voltage    int            `json:"volt"`
}
