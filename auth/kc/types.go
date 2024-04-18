package kcm

// KeycloakConfig holds the necessary configuration for connecting to Keycloak.
type KeycloakConfig struct {
	URL                  string `json:"kc_url"`
	Realm                string `json:"kc_realm"`
	ClientID             string `json:"kc_client_id"`
	Secret               string `json:"kc_secret"`
	Protocol             string `json:"kc_protocol"`
	RealmDockerRole      string `json:"kc_realm_docker_role"`
	RealmDockerAdminRole string `json:"kc_realm_docker_admin_role"`
	Username             string `json:"username"`
	Password             string `json:"password"`
	CurrentKCToken       string
}

type KeycloakCredentials struct {
	Username string
	Password string
}

// TokenIntrospectionResponse represents the JSON structure returned from the token introspection endpoint.
type TokenIntrospectionResponse struct {
	Active            bool     `json:"active"`
	TokenType         string   `json:"token_type"`
	ClientID          string   `json:"client_id"`
	Username          string   `json:"username"`
	TokenUseType      string   `json:"typ"`
	Azp               string   `json:"azp"`
	SessionState      string   `json:"session_state"`
	Acr               string   `json:"acr"`
	AllowedOrigins    []string `json:"allowed-origins"`
	Scope             string   `json:"scope"`
	SID               string   `json:"sid"`
	EmailVerified     bool     `json:"email_verified"`
	PreferredUsername string   `json:"preferred_username"`
	Exp               int      `json:"exp"`
	Iat               int      `json:"iat"`
	Jti               string   `json:"jti"`
	Iss               string   `json:"iss"`
	Aud               string   `json:"aud"`
	Sub               string   `json:"sub"`
	RealmAccess       struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
	ResourceAccess map[string]struct {
		Roles []string `json:"roles"`
	} `json:"resource_access"`
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
