package pluginconfig

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"

	kcm "docker-auth-plugin/auth/kc"

	"github.com/denisbrodbeck/machineid"
)

func LoadConfig(filename string) (map[string]interface{}, error) {
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
		kcConf.Protocol = "openid-connect" // Set default protocol if not specified
	}
	if kcConf.RealmDockerRole == "" {
		id, err := machineid.ID()
		if err != nil {
			log.Println("Failed to retrieve hardver id, you must set kc_realm_docker_role manually in your config")
			log.Fatal(err)
		}
		kcConf.RealmDockerRole = id
	}

	return &kcConf, nil
}
