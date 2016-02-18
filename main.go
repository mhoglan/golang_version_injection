package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
        "os"
)

//go:generate go run scripts/includetxt.go version_info

type JsonVersionInfo struct {
	Version_info map[string]string    //see textfile_constants.go for the fields that get put here
}

var BUILD_LABEL string

func main() {

	command := os.Args[1]
	switch command {
	case "version":
		buf := new(bytes.Buffer)
		json.Indent(buf, []byte(version_info), "", "  ")
		fmt.Println(buf)
	default:
		log.Fatal("Unknown subcommand: ", command)
	}
}
