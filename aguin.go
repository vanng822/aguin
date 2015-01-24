package main

import (
	aguin_api "aguin/api"
	"aguin/config"
	"aguin/model"
	"aguin/utils"
	"flag"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/vanng822/gopid"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func newMartini() *martini.ClassicMartini {
	r := martini.NewRouter()
	m := martini.New()
	m.Use(martini.Logger())
	m.Use(martini.Recovery())
	m.MapTo(r, (*martini.Routes)(nil))
	m.Action(r.Handle)
	return &martini.ClassicMartini{m, r}
}

func main() {
	var (
		configPath string
		host       string
		port       int
		pidFile    string
		force      bool
	)
	log := utils.GetLogger("system")
	flag.StringVar(&host, "h", "", "Host to listen on")
	flag.IntVar(&port, "p", 0, "Port number to listen on")
	flag.StringVar(&configPath, "c", "", "Path to configurations")
	flag.StringVar(&pidFile, "pid", "", "Pid file")
	flag.BoolVar(&force, "f", false, "Force and remove pid file")
	flag.Parse()
	if configPath != "" {
		config.SetConfigPath(configPath)
	}
	
	config.ReadConfig()
	model.EnsureIndex(true)
	api := newMartini()
	api.Use(render.Renderer())
	api.Use(aguin_api.VerifyRequest())
	api.Get("/", aguin_api.IndexGet)
	api.Post("/", aguin_api.IndexPost)
	api.Options("/", func(res http.ResponseWriter) {
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
	if pidFile != "" {
		serverConfig.PidFile = pidFile
	}

	if serverConfig.PidFile != "" {
		gopid.CheckPid(serverConfig.PidFile, force)
		gopid.CreatePid(serverConfig.PidFile)
		defer gopid.CleanPid(serverConfig.PidFile)	
	}
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Kill, os.Interrupt, syscall.SIGTERM)
	log.Info("listening to address %s:%d", serverConfig.Host, serverConfig.Port)
	go http.ListenAndServe(fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port), api)
	sig := <-sigc
	log.Info("Got signal: %s", sig)
}
