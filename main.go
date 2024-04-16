package main

import (
	pluginconfig "docker-auth-plugin/core/config"
	plugin "docker-auth-plugin/plugin"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"strconv"

	"github.com/docker/go-plugins-helpers/authorization"
)

const (
	pluginSocket = "/run/docker/plugins/kknyblnt-docker-auth-plugin.sock"
)

func getEnvOrFlag(envKey string, flagVal *string) string {
	if value, exists := os.LookupEnv(envKey); exists {
		return value
	}
	return *flagVal
}

func main() {
	log.Println("kknyblnt/docker-auth-plugin")

	usernameFlag := flag.String("username", "", "Specifies the username (you can specify the username with the DOCKER_AUTH_PLUGIN_KC_USERNAME env variable)")
	passwordFlag := flag.String("password", "", "Specifies the password (you can specify the password with the DOCKER_AUTH_PLUGIN_KC_PASSWORD env variable)")
	flag.Parse()

	username := getEnvOrFlag("DOCKER_AUTH_PLUGIN_KC_USERNAME", usernameFlag)
	password := getEnvOrFlag("DOCKER_AUTH_PLUGIN_KC_PASSWORD", passwordFlag)

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [OPTIONS]")
		flag.PrintDefaults()
		return
	}

	switch os.Args[1] {
	case "--help":
		fmt.Println("Usage: go run main.go [OPTIONS]")
		flag.PrintDefaults()
		return
	}

	if username == "" || password == "" {
		fmt.Println("Error: Username or password not provided")
		fmt.Println("Usage: go run main.go [OPTIONS]")
		flag.PrintDefaults()
		return
	}

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

	log.Println("Config loaded successfully")

	plugin := plugin.NewDockerAuthPlugin(keycloakConfig)
	handler := authorization.NewHandler(plugin)

	log.Println("Unix socket serve started...")

	err = handler.ServeUnix(pluginSocket, gid)
	if err != nil {
		log.Fatalf("Error serving plugin: %v", err)
	}

}
