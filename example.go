package main

import (
	"encoding/json"
	"fmt"
	"github.com/mcrute/go-inform/inform"
	"io/ioutil"
	"os"
)

func main() {
	fp, err := os.Open("data/test_files/1.bin")
	if err != nil {
		fmt.Println("Error loading file")
		return
	}
	defer fp.Close()

	kp, err := os.Open("data/device_keys.json")
	if err != nil {
		fmt.Println("Error loading key file")
		return
	}
	defer kp.Close()

	var keys map[string]string
	kd, _ := ioutil.ReadAll(kp)
	json.Unmarshal(kd, &keys)

	codec := &inform.Codec{keys}

	msg, err := codec.Unmarshal(fp)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("%s", msg)

	out, _ := os.Create("test.out")
	defer out.Close()

	pkt, err := codec.Marshal(msg)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	out.Write(pkt)
}
