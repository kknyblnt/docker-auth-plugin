package config

import (
	kcm "docker-auth-plugin/auth/kc"
	"encoding/json"
	"errors"
	"io/ioutil"
)

func loadConfig(filename string) (map[string]interface{}, error) {
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

func parseKCMConfig(c map[string]interface{}) (*kcm.KeycloakConfig, error) {

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
		kcConf.Protocol = "openid-connect" // Set default protocol if not specified
	}

	return &kcConf, nil
}
