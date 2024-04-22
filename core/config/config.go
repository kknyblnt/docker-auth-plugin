package pluginconfig

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	kcm "docker-auth-plugin/auth/kc"
)

func LoadConfigFromFile(filename string) (map[string]interface{}, error) {
	var config map[string]interface{}

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func ParseKCMConfig(c map[string]interface{}) (*kcm.KeycloakConfig, error) {

	jsonConfig, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	var kcConf kcm.KeycloakConfig
	err = json.Unmarshal(jsonConfig, &kcConf)
	if err != nil {
		return nil, errors.New("error parsing Keycloak configuration: " + err.Error())
	}

	if kcConf.Protocol == "" {
		kcConf.Protocol = "openid-connect"
	}
	if kcConf.RealmDockerRole == "" {
		kcConf.RealmDockerRole = "docker-auth-plugin-roles"
	}
	if kcConf.RealmDockerAdminRole == "" {
		kcConf.RealmDockerAdminRole = "docker-auth-plugin-admin-roles"
	}

	return &kcConf, nil
}
