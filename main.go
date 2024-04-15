package main

import (
	pluginconfig "docker-auth-plugin/core/config"
	plugin "docker-auth-plugin/plugin"
	"flag"
	"fmt"
	"log"
	"os/user"
	"strconv"

	"github.com/docker/go-plugins-helpers/authorization"
)

const (
	pluginSocket = "/run/docker/plugins/kknyblnt-docker-auth-plugin.sock"
)

func main() {
	fmt.Println("Started.....")

	mode := "kc"

	if mode == "kc" {

		flag.Parse()
		u, _ := user.Lookup("root")
		gid, _ := strconv.Atoi(u.Gid)

		configData, err := pluginconfig.LoadConfig("config.json")
		if err != nil {
			log.Fatalf("Error loading config: %v", err)
		}
		keycloakConfig, err := pluginconfig.ParseKCMConfig(configData)
		if err != nil {
			log.Fatalf("Error parsing Keycloak config: %v", err)
		}

		plugin := plugin.NewDockerAuthPlugin(keycloakConfig)
		handler := authorization.NewHandler(plugin)
		err = handler.ServeUnix(pluginSocket, gid)
		if err != nil {
			log.Fatalf("Error serving plugin: %v", err)
		}
	}
	if mode == "debug" {
		fmt.Println("debug")
	}

}
