package plugin

import (
	kcm "docker-auth-plugin/auth/kc"
	"fmt"
	"log"
)

func kcmGetAccessToken(keycloakConfig kcm.KeycloakConfig) (*kcm.TokenResponse, error) {
	tokenResponse, err := keycloakConfig.GetAccessToken()
	if err != nil {
		log.Fatalf("Error getting access token: %v", err)
		return nil, err
	}
	return tokenResponse, nil
}

func kcmtokenIntrospect(keycloakConfig kcm.KeycloakConfig, accessToken string) (*kcm.TokenIntrospectionResponse, error) {
	introspectResponse, err := keycloakConfig.IntrospectToken(accessToken)
	if err != nil {
		log.Fatalf("Error introspecting: %v", err)
		return nil, err
	}
	return introspectResponse, nil
}

func kcmLogout(keycloakConfig *kcm.KeycloakConfig, refreshToken string) {
	err := keycloakConfig.Logout(refreshToken)
	if err != nil {
		fmt.Println("Logout failed:", err)
	} else {
		fmt.Println("Successfully logged out.")
	}
}
