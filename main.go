package main

import (
	"fmt"
	"log"
	"os"

	pluginconfig "docker-auth-plugin/core/config"
)

func kcmStart(filepath string) {

	// Load Keycloak configuration from file
	configData, err := pluginconfig.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	keycloakConfig, err := pluginconfig.ParseKCMConfig(configData)
	if err != nil {
		log.Fatalf("Error parsing Keycloak config: %v", err)
	}

	tokenResponse, err := keycloakConfig.GetAccessToken()
	if err != nil {
		log.Fatalf("Error getting access token: %v", err)
	}

	// Correctly print the access token
	fmt.Println("Access Token:", tokenResponse.AccessToken)
	fmt.Println("Token Type:", tokenResponse.TokenType)
	fmt.Println("Expires In:", tokenResponse.ExpiresIn)
	fmt.Println("Refresh Token:", tokenResponse.RefreshToken)
	fmt.Println("Refresh Expires In:", tokenResponse.RefreshExpiresIn)
	fmt.Println("Not Before Policy:", tokenResponse.NotBeforePolicy)
	fmt.Println("Session State:", tokenResponse.SessionState)
	fmt.Println("Scope:", tokenResponse.Scope)

}

func main() {
	fmt.Println("Started.....")

	mode := "debug"

	if mode == "kc" {
		kcmStart("config.json")
	}
	if mode == "debug" {
		fmt.Println("debug")
	} else {
		fmt.Println("Invalid mode.....")
		os.Exit(1)
	}

}
