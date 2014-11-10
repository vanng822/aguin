package config
import (
	"fmt"
	"strings"
	"os"
	"encoding/json"
)

type AppConfig struct {
	Mongodb string
	EncryptionEnabled bool
}

type ServerConfig struct {
	Port int
	Host string
}

var (
	app AppConfig
	server ServerConfig
	configPath string
)

func init() {
	configPath = "./config"
}
func ServerConf() ServerConfig {
	return server
}

func AppConf() AppConfig {
	return app
}

func ReadConfig(path string) {
	configPath = strings.TrimRight(path, "/")
	serr := ReadServerConfig(fmt.Sprintf("%s/conf.json", configPath))
	merr := ReadAppConfig(fmt.Sprintf("%s/app.json", configPath))
	if merr != nil || serr != nil {
		panic(fmt.Sprintf("Can not load configuration:%s/app.json or %s/conf.json", configPath, configPath))
	}
}

func ReadAppConfig(filename string) error {
	file, _ := os.Open(filename)
	decoder := json.NewDecoder(file)
	return decoder.Decode(&app)
}

func ReadServerConfig(filename string) error {
	file, _ := os.Open(filename)
	decoder := json.NewDecoder(file)
	return decoder.Decode(&server)
}
