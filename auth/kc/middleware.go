package kcm

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// KeycloakConfig holds the necessary configuration for connecting to Keycloak.
type KeycloakConfig struct {
	URL      string `json:"kc_url"`
	Realm    string `json:"kc_realm"`
	ClientID string `json:"kc_client_id"`
	Secret   string `json:"kc_secret"`
	Protocol string `json:"kc_protocol"` // Optional, defaults to "openid-connect"
}

type KeycloakCredentials struct {
	Username string
	Password string
}

// TokenResponse represents the response from Keycloak token endpoint.
type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

// NewKeycloakConfig creates a new KeycloakConfig with default protocol.
func NewKeycloakConfig(url, realm, clientID, secret string) *KeycloakConfig {
	return &KeycloakConfig{
		URL:      url,
		Realm:    realm,
		ClientID: clientID,
		Secret:   secret,
		Protocol: "openid-connect",
	}
}

func NewKeycloakCredentials(username, password string) *KeycloakCredentials {
	return &KeycloakCredentials{
		Username: username,
		Password: password,
	}
}

// GetAccessToken fetches the access token from Keycloak.
func (kc *KeycloakConfig) GetAccessToken(creds KeycloakCredentials) (*TokenResponse, error) {
	url := fmt.Sprintf("%s/realms/%s/protocol/%s/token", kc.URL, kc.Realm, kc.Protocol)

	data := strings.NewReader(fmt.Sprintf("client_id=%s&client_secret=%s&grant_type=password&username=%s&password=%s", kc.ClientID, kc.Secret, creds.Username, creds.Password))

	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tokenResponse TokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return nil, err
	}

	return &tokenResponse, nil
}
