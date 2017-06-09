package main

import (
	"blue/compiler"
	"blue/printer"
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	blueFilePath := flag.Arg(0)
	blueFilePath = "test.blue"
	if blueFilePath == "" {
		fmt.Println("blue file not specified")
		os.Exit(1)
	}

	blueFilePath, err := filepath.Abs(blueFilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	blueFile, err := os.Open(blueFilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer blueFile.Close()

	instructions, err := compiler.BuildSource(blueFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = printer.CreateBlueprint(instructions, os.Stdout)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Stdout.WriteString("\n")

	//readTestBlueprint()
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
