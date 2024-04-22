package plugin

import (
	"log"
	"slices"

	kcm "docker-auth-plugin/auth/kc"

	"github.com/docker/go-plugins-helpers/authorization"
)

type AuthPlugin struct {
	dockerAuthPlugin DockerAuthPlugin
	tokenResponse    kcm.TokenResponse
}

type DockerAuthPlugin struct {
	keycloakConfig *kcm.KeycloakConfig
}

func NewAuthPlugin(cfg *kcm.KeycloakConfig) *AuthPlugin {
	return &AuthPlugin{
		dockerAuthPlugin: *NewDockerAuthPlugin(cfg),
	}
}

func NewDockerAuthPlugin(cfg *kcm.KeycloakConfig) *DockerAuthPlugin {
	return &DockerAuthPlugin{
		keycloakConfig: cfg,
	}
}

func LoginFlow(keycloakConfig *kcm.KeycloakConfig, username string, password string) (*kcm.TokenResponse, error) {
	tokenResponse, err := kcmHandleGetAccessToken(*keycloakConfig, *kcm.NewKeycloakCredentials(username, password))
	if err != nil {
		log.Printf("Authorization failed (probably a KC failure while getting access token): %v", err)
		return nil, err
	}

	introspectResponse, err := kcmHandleTokenIntrospect(*keycloakConfig, tokenResponse)
	if err != nil || !introspectResponse.Active {
		log.Printf("Access introspect failed")
		return nil, err
	}
	log.Println("Access granted")
	return tokenResponse, nil
}

func Introspect(keycloakConfig *kcm.KeycloakConfig, tokenResponse *kcm.TokenResponse) (*kcm.TokenIntrospectionResponse, error) {
	introspectResponse, err := kcmHandleTokenIntrospect(*keycloakConfig, tokenResponse)
	if err != nil || !introspectResponse.Active {
		log.Printf("Access introspect failed")
		return nil, err
	}
	log.Println("Access granted")
	return introspectResponse, nil
}

func (p *DockerAuthPlugin) AuthZReq(req authorization.Request) authorization.Response {

	req_allowed := false
	tokenResponse, err := LoginFlow(p.keycloakConfig, p.keycloakConfig.Username, p.keycloakConfig.Password)
	if err != nil {
		log.Fatal("Failure with the login flow")
		req_allowed = false
	}

	introspectResponse, err := Introspect(p.keycloakConfig, tokenResponse)
	if err != nil {
		log.Fatal("Failure with the introspect flow")
		req_allowed = false
	}

	if slices.Contains(introspectResponse.RealmAccess.Roles, p.keycloakConfig.RealmDockerAdminRole) {
		log.Println("Introspect response contains ADMIN realm role")
		req_allowed = true
	}

	if slices.Contains(introspectResponse.RealmAccess.Roles, p.keycloakConfig.RealmDockerRole) {
		log.Println("Introspect response contains realm role")
		req_allowed = true
	}

	kcmHandleLogout(p.keycloakConfig, tokenResponse.RefreshToken)
	return authorization.Response{Allow: req_allowed}

}

func (p *DockerAuthPlugin) AuthZRes(req authorization.Request) authorization.Response {
	return authorization.Response{Allow: true}
}
