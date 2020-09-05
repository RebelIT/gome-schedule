package main

import (
	"github.com/rebelit/gome-schedule/common/config"
	"github.com/rebelit/gome-schedule/schedule"
	"github.com/rebelit/gome-schedule/schedule/web"
	"log"
	"net/http"
)

func main() {
	log.Printf("INFO: I'm starting")
	config.Runtime()
	schedule.InitializeDatabases()

	go schedule.Runner()

	start(config.App.ListenPort)
	return
}

func start(port string) {
	log.Printf("INFO: listening %s on http:%s", config.App.Name, port)
	router := web.NewRouter()
	log.Fatal(http.ListenAndServe(":"+port, router))
}
