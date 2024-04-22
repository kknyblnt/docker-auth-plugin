package api

import (
	kcm "docker-auth-plugin/auth/kc"
	pluginconfig "docker-auth-plugin/core/config"
	plugin "docker-auth-plugin/plugin"
	"encoding/json"
	"log"
	"net/http"
)

func LoginHandler(pluginData *plugin.AuthPlugin, w http.ResponseWriter, r *http.Request) {
	log.Println("Login attempted")
	username, password, ok := r.BasicAuth()
	if !ok {
		http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
		return
	}

	var loginData kcm.KeycloakConfig
	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
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

	plugin.LoginFlow()

	var responseData []byte

	responseData, _ = json.Marshal(map[string]string{"message": "Logged in successfully"})
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseData)
}
func LogoutHandler(plugin *plugin.DockerAuthPlugin, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged out successfully"))
}
