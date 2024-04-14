package main

import (
	"fmt"
	"log"

	kcm "docker-auth-plugin/auth/kc"
	pluginconfig "docker-auth-plugin/core/config"
)

func kcmHandleGetAccessToken(keycloakConfig kcm.KeycloakConfig, creds kcm.KeycloakCredentials) (*kcm.TokenResponse, error) {
	tokenResponse, err := keycloakConfig.GetAccessToken(*kcm.NewKeycloakCredentials(creds.Username, creds.Password))
	if err != nil {
		log.Fatalf("Error getting access token: %v", err)
		return nil, err
	}
	return tokenResponse, nil
}

func kcmHandleTokenIntrospect(keycloakConfig kcm.KeycloakConfig, getAccessTokenResponse *kcm.TokenResponse) (*kcm.TokenIntrospectionResponse, error) {
	introspectResponse, err := keycloakConfig.IntrospectToken(getAccessTokenResponse)
	if err != nil {
		log.Fatalf("Error getting access token: %v", err)
		return nil, err
	}
	return introspectResponse, nil
}

func kcmHandleLogout(keycloakConfig *kcm.KeycloakConfig, refreshToken string) {
	err := keycloakConfig.Logout(refreshToken)
	if err != nil {
		fmt.Println("Logout failed:", err)
	} else {
		fmt.Println("Successfully logged out.")
	}
}

func kcmStart(filepath string, username string, password string) {

	// Load Keycloak configuration from file
	configData, err := pluginconfig.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	keycloakConfig, err := pluginconfig.ParseKCMConfig(configData)
	if err != nil {
		log.Fatalf("Error parsing Keycloak config: %v", err)
	}

	//GET ACCESS TOKEN
	tokenResponse, err := kcmHandleGetAccessToken(*keycloakConfig, *kcm.NewKeycloakCredentials(username, password))
	if err != nil {
		log.Fatalf("Error getting access token: %v", err)
	}

	fmt.Println("ACCESS TOKEN RESPONSE:")
	// Correctly print the access token
	fmt.Println("Access Token:", tokenResponse.AccessToken)
	fmt.Println("Token Type:", tokenResponse.TokenType)
	fmt.Println("Expires In:", tokenResponse.ExpiresIn)
	fmt.Println("Refresh Token:", tokenResponse.RefreshToken)
	fmt.Println("Refresh Expires In:", tokenResponse.RefreshExpiresIn)
	fmt.Println("Not Before Policy:", tokenResponse.NotBeforePolicy)
	fmt.Println("Session State:", tokenResponse.SessionState)
	fmt.Println("Scope:", tokenResponse.Scope)

	//INTROSPECT
	introspecResponse, err := kcmHandleTokenIntrospect(*keycloakConfig, tokenResponse)
	if err != nil {
		log.Fatalf("Error introspecting access token: %v", err)
	}

	fmt.Println("Active??? ", introspecResponse.Active)

	kcmHandleLogout(keycloakConfig, tokenResponse.RefreshToken)

}

func main() {
	fmt.Println("Started.....")

	mode := "kc"

	if mode == "kc" {
		usernmame := "demouser"
		password := "demouser"
		kcmStart("config.json", usernmame, password)
	}
	if mode == "debug" {
		fmt.Println("debug")
	}

}
