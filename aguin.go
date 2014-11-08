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
		host string
		port int
	)
	log := utils.GetLogger("system")
	flag.StringVar(&host, "h", "", "Host to listen on")
	flag.IntVar(&port, "p", 0, "Port number to listen on")
	
	flag.Parse()
	config.ReadConfig("./config")
	api := martini.Classic()
	api.Use(render.Renderer())
	api.Use(aguin_api.VerifyRequest())
	api.Get("/", aguin_api.IndexGet)
	api.Post("/", aguin_api.IndexPost)
	api.Get("/status", aguin_api.IndexStatus)
	api.NotFound(aguin_api.NotFound)
	serverConfig := config.ServerConf()
	if port > 0 {
		serverConfig.Port = port
	}
	if host != "" {
		serverConfig.Host = host
	}
	log.Print(fmt.Sprintf("listening to address %s:%d", serverConfig.Host, serverConfig.Port))
	http.ListenAndServe(fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port), api)
}
