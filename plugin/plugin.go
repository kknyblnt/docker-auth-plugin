package plugin

import (
	"log"

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
	tokenResponse, err := kcmHandleGetAccessToken(*p.keycloakConfig, *kcm.NewKeycloakCredentials("demouser", "demouser")) // Placeholder credentials
	if err != nil {
		log.Printf("Authorization failed (probably a KC failure while getting access token): %v", err)
		return authorization.Response{Allow: false, Msg: "Access denied by kknyblnt/docker-auth-plugin"}
	}
	introspectResponse, err := kcmHandleTokenIntrospect(*p.keycloakConfig, tokenResponse)
	if err != nil || !introspectResponse.Active {
		return authorization.Response{Allow: false, Msg: "Access denied by kknyblnt/docker-auth-plugin"}
	}

	// Check user permissions specific to the request...
	// This is a simplified example; implement detailed checks as per your requirements.

	return authorization.Response{Allow: true}
}

func (p *DockerAuthPlugin) AuthZRes(req authorization.Request) authorization.Response {
	return authorization.Response{Allow: true} // Implement as needed
}
