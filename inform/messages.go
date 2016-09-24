package inform

import (
	"encoding/json"
	"strconv"
	"time"
)

type AdminMetadata struct {
	Id       string `json:"_id"`
	Language string `json:"lang"`
	Username string `json:"name"`
	Password string `json:"x_password"`
}

type CommandMessage struct {
	Metadata   *AdminMetadata `json:"_admin,omitempty"`
	Id         string         `json:"_id,omitempty"`
	Type       string         `json:"_type"`
	Command    string         `json:"cmd"`
	DateTime   string         `json:"datetime"`
	DeviceId   string         `json:"device_id,omitempty"`
	MacAddress string         `json:"mac,omitempty"`
	Model      string         `json:"model,omitempty"`
	OffVoltage int            `json:"off_volt,omitempty"`
	Port       int            `json:"port"`
	SensorId   string         `json:"sId,omitempty"`
	ServerTime string         `json:"server_time_in_utc"`
	Time       int64          `json:"time"`
	Timer      int            `json:"timer"`
	Value      int            `json:"val"`
	Voltage    int            `json:"volt,omitempty"`
}

func NewOutputCommand(port, val, timer int) *CommandMessage {
	return &CommandMessage{
		Type:       "cmd",
		Command:    "mfi-output",
		DateTime:   time.Now().Format(time.RFC3339),
		Port:       port,
		ServerTime: unixMicroPSTString(),
		Time:       unixMicroPST(),
		Timer:      timer,
		Value:      val,
	}
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

type DeviceMessage struct {
	IsDefault       bool   `json:"default"`
	IP              string `json:"ip"`
	MacAddr         string `json:"mac"`
	ModelNumber     string `json:"model"`
	ModelName       string `json:"model_display"`
	Serial          string `json:"serial"`
	FirmwareVersion string `json:"version"`
	Outputs         []*OutputInfo
}

func (m *DeviceMessage) UnmarshalJSON(data []byte) error {
	type Alias DeviceMessage
	aux := &struct {
		Alarm []struct {
			Entries []struct {
				Tag  string      `json:"tag"`
				Type string      `json:"type"`
				Val  interface{} `json:"val"`
			} `json:"entries"`
			Sensor string `json:"sId"`
		} `json:"alarm"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	m.Outputs = make([]*OutputInfo, len(aux.Alarm))

	for i, a := range aux.Alarm {
		o := &OutputInfo{
			Id:     a.Sensor,
			Port:   i + 1,
			Dimmer: m.ModelNumber == "IWD1U",
		}
		m.Outputs[i] = o

		for _, e := range a.Entries {
			switch t := e.Val; e.Tag {
			case "output":
				o.OutputState = t.(float64) == 1
			case "pf":
				o.PowerFactor = t.(float64)
			case "energy_sum":
				o.EnergySum = t.(float64)
			case "v_rms":
				o.VoltageRMS = t.(float64)
			case "i_rms":
				o.CurrentRMS = t.(float64)
			case "active_pwr":
				o.Watts = t.(float64)
			case "thismonth":
				o.ThisMonth = t.(float64)
			case "lastmonth":
				o.LastMonth = t.(float64)
			case "dimmer_level":
				o.DimmerLevel = int(t.(float64))
			case "dimmer_lock_setting":
				o.DimmerLockSetting = int(t.(float64))
			}
		}
	}

	return nil
}

type OutputInfo struct {
	Id                string
	Port              int
	OutputState       bool
	EnergySum         float64
	VoltageRMS        float64
	PowerFactor       float64
	CurrentRMS        float64
	Watts             float64
	ThisMonth         float64
	LastMonth         float64
	Dimmer            bool
	DimmerLevel       int
	DimmerLockSetting int
}
