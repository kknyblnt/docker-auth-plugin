package main

import (
	"fmt"
	"os"
)

func kcmStart(filepath string) {

	// Load Keycloak configuration from file
	configData, err := loadConfig("hello.json")
}

func main() {
	fmt.Println("Started.....")

	mode := "kc"

	if mode == "kc" {
		kcmStart("config.json")
	} else {
		fmt.Println("Invalid mode.....")
		os.Exit(1)
	}

}
