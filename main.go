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
	"golang.org/x/term"
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

func readSecureInput() string {
	fmt.Println("(input hidden)")
	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Failed to read password")
		return ""
	}
	return string(bytePassword)
}

func main() {
	log.Println("kknyblnt/docker-auth-plugin")

	usernameFlag := flag.String("username", "", "Specifies the username (you can specify the username with the DOCKER_AUTH_PLUGIN_KC_USERNAME env variable)")
	passwordFlag := flag.String("password", "", "Specifies the password (you can specify the password with the DOCKER_AUTH_PLUGIN_KC_PASSWORD env variable)")
	configFilePathFlag := flag.String("config", "", "Specifies the config file path, by default its PWD/config.json (you can specify the username with the DOCKER_AUTH_PLUGIN_KC_CONFIG env variable)")
	readFlag := flag.Bool("read", false, "Reads username and password")
	otpFlag := flag.String("otp", "", "Specifies the one-time-password (you can specify the one-time-password with the DOCKER_AUTH_PLUGIN_KC_OTP env variable)")
	helpFlag := flag.Bool("help", false, "Display this help and exit")

	flag.Parse()

	var username, password, otp string

	if *helpFlag {
		flag.PrintDefaults()
		return
	}

	if *readFlag {
		fmt.Print("Enter Username: ")
		username = readSecureInput()
		fmt.Print("Enter Password: ")
		password = readSecureInput()
		fmt.Printf("Enter 2FA OTP: (LEAVE IT EMPTY IF YOU DON'T USE 2FA)")
		otp = readSecureInput()
	} else {
		username = getEnvOrFlag("DOCKER_AUTH_PLUGIN_KC_USERNAME", usernameFlag)
		password = getEnvOrFlag("DOCKER_AUTH_PLUGIN_KC_PASSWORD", passwordFlag)
		otp = getEnvOrFlag("DOCKER_AUTH_PLUGIN_KC_OTP", otpFlag)
	}

	configFilePath := getEnvOrFlag("DOCKER_AUTH_PLUGIN_KC_CONFIG", configFilePathFlag)

	if configFilePath == "" {
		configFilePath = "config.json"
	}

	u, _ := user.Lookup("root")
	gid, _ := strconv.Atoi(u.Gid)

	configData, err := pluginconfig.LoadConfig(configFilePath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	keycloakConfig, err := pluginconfig.ParseKCMConfig(configData)
	if err != nil {
		log.Fatalf("Error parsing Keycloak config: %v", err)
	}

	keycloakConfig.Username = username
	keycloakConfig.Password = password
	keycloakConfig.Otp = otp

	log.Println("Config loaded successfully")

	plugin := plugin.NewDockerAuthPlugin(keycloakConfig)
	handler := authorization.NewHandler(plugin)

	log.Println("Unix socket serve started...")

	err = handler.ServeUnix(pluginSocket, gid)
	if err != nil {
		log.Fatalf("Error serving plugin: %v", err)
	}

}
