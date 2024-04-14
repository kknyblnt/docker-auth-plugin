package kcm

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Print request and response details for debugging
	fmt.Println("Request URL:", req.URL)
	fmt.Println("Request Body:", data.Encode())
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(body))

	// Unmarshal the response into the TokenIntrospectionResponse struct
	var introspectionResponse TokenIntrospectionResponse
	if err := json.Unmarshal(body, &introspectionResponse); err != nil {
		return nil, err
	}

	return &introspectionResponse, nil
}
