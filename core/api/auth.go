package api

import (
	kcm "docker-auth-plugin/auth/kc"
	pluginconfig "docker-auth-plugin/core/config"
	"encoding/json"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
		return
	}

	// Process Keycloak login data
	var loginData kcm.KeycloakConfig
	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		http.Error(w, "Error parsing Keycloak data", http.StatusBadRequest)
		return
	}

	kcData := map[string]interface{}{
		"url":       loginData.URL,
		"realm":     loginData.Realm,
		"client_id": loginData.ClientID,
		"secret":    loginData.Secret,
	}

	kcConfig, err := pluginconfig.ParseKCMConfig(kcData)
	if err != nil {
		http.Error(w, "Failed to parse Keycloak configuration: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Include username and password in the response for demonstration purposes (not recommended in production)
	response := struct {
		*kcm.KeycloakConfig
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		KeycloakConfig: kcConfig,
		Username:       username,
		Password:       password,
	}

	responseData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error generating response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseData)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Placeholder for logout logic
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged out successfully"))
}
