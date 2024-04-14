package main

import (
	"fmt"
	"log"
	"os"
)

func kcmStart(filepath string) {

	// Load Keycloak configuration from file
	kcmConfigRaw := core.loadConfig(filepath)

	// Parse kcm config
	kcmConfig = core.parseKCMConfig(kcmConfigRaw)

	// Retrieve access token
	token, err := kcmConfig.GetAccessToken()
	if err != nil {
		log.Fatalf("Error retrieving access token: %v", err)
	}

	// Print the access token details
	fmt.Println("Access Token:", token.AccessToken)
	fmt.Println("Token Type:", token.TokenType)
	fmt.Println("Expires In:", token.ExpiresIn)
	fmt.Println("Refresh Token:", token.RefreshToken)
	fmt.Println("Refresh Expires In:", token.RefreshExpiresIn)
	fmt.Println("Not Before Policy:", token.NotBeforePolicy)
	fmt.Println("Session State:", token.SessionState)
	fmt.Println("Scope:", token.Scope)

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
