package main

import (
	"encoding/json"
	"fmt"
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/mcrute/go-inform/inform"
	"io/ioutil"
	"log"
	"os"
)

// Load devices into state
// Gather current initial state from devices
// Track state transitions
// Inputs:
// - Homekit
// - Devices

type Port struct {
	Label string `json:"label"`
	Port  int    `json:"port"`
}

type Device struct {
	Key    string  `json:"key"`
	Name   string  `json:"name"`
	Model  string  `json:"model"`
	Serial string  `json:"serial"`
	Ports  []*Port `json:"ports"`
}

type DeviceMap map[string]*Device

func LoadKeys(file string) (DeviceMap, error) {
	var keys DeviceMap

	kp, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer kp.Close()

	kd, err := ioutil.ReadAll(kp)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(kd, &keys)
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func main() {
	devs, err := LoadKeys("data/device_keys.json")
	if err != nil {
		log.Println("Error loading key file")
		log.Println(err.Error())
		return
	}

	keys := make(map[string]string, len(devs))
	for i, d := range devs {
		keys[i] = d.Key
	}

	h := inform.NewInformHandler(&inform.Codec{keys})
	s, _ := inform.NewServer(h)
	as := make([]*accessory.Accessory, 0, len(devs)*3)

	for i, d := range devs {
		for _, p := range d.Ports {
			a := accessory.NewSwitch(accessory.Info{
				Name:         p.Label,
				SerialNumber: fmt.Sprintf("%s-%d", d.Serial, p.Port),
				Manufacturer: "Ubiquiti",
				Model:        d.Model,
			})

			// Capture these for the closure, otherwise they're bound to the
			// single loop variable and will only see the final value of that
			// variable
			dev, port := i, p.Port

			a.Switch.On.OnValueRemoteUpdate(func(on bool) {
				h.SetState(dev, port, on)
			})

			h.AddPort(dev, port)
			as = append(as, a.Accessory)
		}
	}

	// The root accessory is what gets used to name the bridge so let's make it
	// an actual bridge
	br := accessory.New(accessory.Info{
		Name:         "UnifiBridge",
		Manufacturer: "Mike Crute",
		Model:        "0.1",
	}, accessory.TypeBridge)

	config := hc.Config{
		Pin:         "12344321",
		Port:        "12345",
		StoragePath: "./db",
	}

	t, err := hc.NewIPTransport(config, br, as...)
	if err != nil {
		log.Fatal(err)
		return
	}

	hc.OnTermination(func() {
		t.Stop()
		os.Exit(0) // Otherwise homekit doesn't actually stop
	})

	go t.Start()
	s.ListenAndServe()
}
