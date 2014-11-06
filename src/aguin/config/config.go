package config
import (
	"fmt"
	"strings"
	"os"
	"encoding/json"
)

type AppConfig struct {
	Mongodb string
	EncryptResponse bool
	CryptedRequest bool
}

type ServerConfig struct {
	Port int
	Host string
}

var (
	app  AppConfig
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
	merr := ReadAppConfig(fmt.Sprintf("%s/app.json", configPath))
	serr := ReadServerConfig(fmt.Sprintf("%s/conf.json", configPath))
	if merr != nil || serr != nil {
		panic(fmt.Sprintf("Can not load configuration"))
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
