package plugin

import (
	"log"
	"slices"
	"time"

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

func validateIntrospect(introspectResponse kcm.TokenIntrospectionResponse, keycloakConfig kcm.KeycloakConfig) bool {

	if slices.Contains(introspectResponse.RealmAccess.Roles, keycloakConfig.RealmDockerAdminRole) {
		log.Println("Introspect response contains ADMIN realm role")
		return true
	}

	if slices.Contains(introspectResponse.RealmAccess.Roles, keycloakConfig.RealmDockerRole) {
		log.Println("Introspect response contains realm role")
		return true
	}

	return false
}

var AuthZFailureResponse = authorization.Response{Allow: false, Msg: "Access denied by kknyblnt/docker-auth-plugin"}

func (p *DockerAuthPlugin) AuthZReq(req authorization.Request) authorization.Response {
	req_allowed := false
	if p.keycloakConfig.CurrentKCToken == "" || time.Now().After(p.keycloakConfig.TokenExpiration) {
		tokenResponse, err := kcmGetAccessToken(*p.keycloakConfig)
		if err != nil {
			log.Printf("Authorization failed (probably a KC failure while getting access token): %v", err)
			return AuthZFailureResponse
		}
		log.Println("Token granted")
		p.keycloakConfig.CurrentKCToken = tokenResponse.AccessToken
	}
	introspectResponse, err := kcmtokenIntrospect(*p.keycloakConfig, p.keycloakConfig.CurrentKCToken)
	if err != nil || !introspectResponse.Active {
		return AuthZFailureResponse
	}
	p.keycloakConfig.TokenExpiration = time.Unix(int64(introspectResponse.Exp), 0)
	req_allowed = validateIntrospect(*introspectResponse, *p.keycloakConfig)
	return authorization.Response{Allow: req_allowed}
}

func (p *DockerAuthPlugin) AuthZRes(req authorization.Request) authorization.Response {
	return authorization.Response{Allow: true}
}
