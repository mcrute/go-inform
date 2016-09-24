package inform

// Messages we receive from devices

import (
	"encoding/json"
)

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
