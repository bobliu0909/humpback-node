package main

import "github.com/humpback/gounits/system"
import "github.com/bobliu0909/humpback-node/server"

import (
	"log"
	"os"
)

func main() {

	service, err := server.NewNodeService()
	if err != nil {
		log.Printf("service error:%s\n", err.Error())
		os.Exit(system.PorcessExitCode(err))
	}

	defer func() {
		service.Stop()
		os.Exit(0)
	}()

	go startRouter(service.Configuration.API.Host)
	if err := service.Startup(); err != nil {
		log.Printf("service start error:%s\n", err.Error())
		os.Exit(system.PorcessExitCode(err))
	}
	system.InitSignal(nil)
}
