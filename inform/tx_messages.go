package inform

// Messages we send to devices

import (
	"strconv"
	"time"
)

type CommandMessage struct {
	Id         string `json:"_id,omitempty"`
	Type       string `json:"_type"`
	Command    string `json:"cmd"`
	DateTime   string `json:"datetime"`
	DeviceId   string `json:"device_id,omitempty"`
	MacAddress string `json:"mac,omitempty"`
	Model      string `json:"model,omitempty"`
	OffVoltage int    `json:"off_volt,omitempty"`
	Port       int    `json:"port"`
	SensorId   string `json:"sId,omitempty"`
	ServerTime string `json:"server_time_in_utc"`
	Time       int64  `json:"time"`
	Timer      int    `json:"timer"`
	Value      int    `json:"val"`
	Voltage    int    `json:"volt,omitempty"`
}

// Freshen timestamps
func (m *CommandMessage) Freshen() {
	m.DateTime = time.Now().Format(time.RFC3339)
	m.ServerTime = unixMicroPSTString()
	m.Time = unixMicroPST()
}

func NewOutputCommand(port int, val bool, timer int) *CommandMessage {
	m := &CommandMessage{
		Type:       "cmd",
		Command:    "mfi-output",
		DateTime:   time.Now().Format(time.RFC3339),
		Port:       port,
		ServerTime: unixMicroPSTString(),
		Time:       unixMicroPST(),
		Timer:      timer,
	}

	if val {
		m.Value = 1
	} else {
		m.Value = 0
	}

	return m
}

type NoopMessage struct {
	Type          string `json:"_type"`
	Interval      int    `json:"interval"`
	ServerTimeUTC string `json:"server_time_in_utc"`
}

func unixMicroPST() int64 {
	l, _ := time.LoadLocation("America/Los_Angeles")
	tnano := time.Now().In(l).UnixNano()
	return tnano / int64(time.Millisecond)
}

func unixMicroPSTString() string {
	return strconv.FormatInt(unixMicroPST(), 10)
}

func unixMicroUTCString() string {
	tnano := time.Now().UTC().UnixNano()
	t := tnano / int64(time.Millisecond)
	return strconv.FormatInt(t, 10)
}

func NewNoop(interval int) *NoopMessage {
	return &NoopMessage{
		Type:          "noop",
		Interval:      interval,
		ServerTimeUTC: unixMicroUTCString(),
	}
}
