package plugin

import (
	"log"
	"slices"

	kcm "docker-auth-plugin/auth/kc"

	"github.com/docker/go-plugins-helpers/authorization"
)

type DockerAuthPlugin struct {
	keycloakConfig *kcm.KeycloakConfig
}

func NewDockerAuthPlugin(cfg *kcm.KeycloakConfig) *DockerAuthPlugin {
	return &DockerAuthPlugin{
		keycloakConfig: cfg,
	}
}

func (p *DockerAuthPlugin) AuthZReq(req authorization.Request) authorization.Response {

	tokenResponse, err := kcmHandleGetAccessToken(*p.keycloakConfig, *kcm.NewKeycloakCredentials(p.keycloakConfig.Username, p.keycloakConfig.Password))
	if err != nil {
		log.Printf("Authorization failed (probably a KC failure while getting access token): %v", err)
		return authorization.Response{Allow: false, Msg: "Access denied by kknyblnt/docker-auth-plugin"}
	}
	log.Println("Token granted")

	introspectResponse, err := kcmHandleTokenIntrospect(*p.keycloakConfig, tokenResponse)
	if err != nil || !introspectResponse.Active {
		return authorization.Response{Allow: false, Msg: "Access denied by kknyblnt/docker-auth-plugin"}
	}
	log.Println("Access granted")

	if slices.Contains(introspectResponse.RealmAccess.Roles, p.keycloakConfig.RealmDockerAdminRole) {
		log.Println("Introspect response contains ADMIN realm role")
		return authorization.Response{Allow: true}
	}

	if slices.Contains(introspectResponse.RealmAccess.Roles, p.keycloakConfig.RealmDockerRole) {
		log.Println("Introspect response contains realm role")
		return authorization.Response{Allow: true}
	}

	kcmHandleLogout(p.keycloakConfig, tokenResponse.RefreshToken)
	return authorization.Response{Allow: false}

}

func (p *DockerAuthPlugin) AuthZRes(req authorization.Request) authorization.Response {
	return authorization.Response{Allow: true} // Implement as needed
}
