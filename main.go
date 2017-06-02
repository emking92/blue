package main

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	//	"encoding/json"
	//	"blue/entities"
	//	"blue/entities/controls"
	"factorio-assembly/blueprinter"
	"fmt"
	"io/ioutil"
)

func main() {
	blueprinter.CreateBlueprint()

	readTestBlueprint()
}

func readTestBlueprint() {
	inBytes, err := ioutil.ReadFile("in.bp")
	if err != nil {
		fmt.Println(err)
		return
	}

	blueprintStr := string(inBytes[1:])
	fmt.Println(blueprintStr)

	decodedBytes, err := base64.StdEncoding.DecodeString(blueprintStr)
	if err != nil {
		fmt.Println(err)
		return
	}

	gunzipper, err := zlib.NewReader(bytes.NewReader(decodedBytes))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer gunzipper.Close()

	jsonBytes, err := ioutil.ReadAll(gunzipper)
	if err != nil {
		fmt.Println(err)
		return
	}
	gunzipper.Close()

	fmt.Println(string(jsonBytes))
}
