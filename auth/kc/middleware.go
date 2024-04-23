package kcm

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

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

// GetAccessToken fetches the access token from Keycloak.
func (kc *KeycloakConfig) GetAccessToken() (*TokenResponse, error) {
	url := fmt.Sprintf("%s/realms/%s/protocol/%s/token", kc.URL, kc.Realm, kc.Protocol)

	data := strings.NewReader(fmt.Sprintf("client_id=%s&client_secret=%s&grant_type=password&username=%s&password=%s", kc.ClientID, kc.Secret, kc.Username, kc.Password))

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

// IntrospectToken introspects the validity and details of a given access token.
func (kc *KeycloakConfig) IntrospectToken(tokenResponse *TokenResponse) (*TokenIntrospectionResponse, error) {
	// Build the token introspection endpoint URL
	introspectURL := fmt.Sprintf("%s/realms/%s/protocol/%s/token/introspect", kc.URL, kc.Realm, kc.Protocol)

	// Data to be sent in the request body
	data := url.Values{}
	data.Set("token", tokenResponse.AccessToken)
	data.Set("client_id", kc.ClientID)
	data.Set("client_secret", kc.Secret)

	// Create the request
	req, err := http.NewRequest("POST", introspectURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Send the request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response into the TokenIntrospectionResponse struct
	var introspectionResponse TokenIntrospectionResponse
	if err := json.Unmarshal(body, &introspectionResponse); err != nil {
		return nil, err
	}

	return &introspectionResponse, nil
}

// Logout invalidates the given refresh token in Keycloak.
func (kc *KeycloakConfig) Logout(refreshToken string) error {
	// Construct the URL for the logout endpoint
	logoutURL := fmt.Sprintf("%s/realms/%s/protocol/%s/logout", kc.URL, kc.Realm, kc.Protocol)

	// Prepare the data for the POST request
	data := url.Values{}
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", kc.ClientID)
	data.Set("client_secret", kc.Secret)

	// Create the HTTP request
	req, err := http.NewRequest("POST", logoutURL, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Initialize the HTTP client and send the request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check if the response status code indicates success
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to logout, server responded with status code: %d", resp.StatusCode)
	}

	return nil
}
