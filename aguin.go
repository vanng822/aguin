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
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

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
	api := martini.Classic()
	api.Use(render.Renderer())
	api.Use(aguin_api.VerifyRequest())
	api.Get("/", aguin_api.IndexGet)
	api.Post("/", aguin_api.IndexPost)
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
	if pidFile != "" {
		serverConfig.PidFile = pidFile
	}

	log.Info("listening to address %s:%d", serverConfig.Host, serverConfig.Port)

	if serverConfig.PidFile != "" {
		if !force {
			if _, err := os.Stat(serverConfig.PidFile); err == nil {
				panic(fmt.Sprintf("Pidfile %s exist", serverConfig.PidFile))
			}
		}
		pid := syscall.Getpid()
		pidf, err := os.Create(serverConfig.PidFile)
		if err != nil {
			log.Critical(fmt.Sprintf("Could not create pid file, error: %v", err))
			panic("Could not create pid file")
		}
		pidf.WriteString(fmt.Sprintf("%d", pid))
		log.Error("%v", err)
	}
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Kill, os.Interrupt, syscall.SIGTERM)
	go http.ListenAndServe(fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port), api)
	sig := <-sigc
	log.Info("Got signal: %s", sig)
	if serverConfig.PidFile != "" {
		log.Info("Cleaning up pid file")
		err := os.Remove(serverConfig.PidFile)
		if err != nil {
			log.Info("Fail to clean up pid file %s", serverConfig.PidFile)
		}
	}
}
