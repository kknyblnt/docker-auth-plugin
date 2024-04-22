package plugin

import (
	kcm "docker-auth-plugin/auth/kc"
	"log"
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
		log.Println("Logout failed:", err)
	} else {
		log.Println("Successfully logged out.")
	}
}
