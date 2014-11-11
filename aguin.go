package main

import (
	aguin_api "aguin/api"
	"aguin/config"
	"flag"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"aguin/utils"
)

func main() {
	var (
		configPath string
		host string
		port int
	)
	log := utils.GetLogger("system")
	flag.StringVar(&host, "h", "", "Host to listen on")
	flag.IntVar(&port, "p", 0, "Port number to listen on")
	flag.StringVar(&configPath, "c", "", "Path to configurations")
	
	flag.Parse()
	if configPath != "" {
		config.SetConfigPath(configPath)
	}
	config.ReadConfig()
	
	api := martini.Classic()
	api.Use(render.Renderer())
	api.Use(aguin_api.VerifyRequest())
	api.Get("/", aguin_api.IndexGet)
	api.Post("/", aguin_api.IndexPost)
	api.Get("/status", aguin_api.IndexStatus)
	api.Options("/", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Allow", "POST, GET")
	})
	api.NotFound(aguin_api.NotFound)
	serverConfig := config.ServerConf()
	if port > 0 {
		serverConfig.Port = port
	}
	if host != "" {
		serverConfig.Host = host
	}
	log.Info("listening to address %s:%d", serverConfig.Host, serverConfig.Port)
	http.ListenAndServe(fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port), api)
}
