package main

import (
	"log"
	"net/http"

	"github.com/daemonl/registerapi/api"

	"gopkg.daemonl.com/envconf"
)

var config struct {
	Bind string `env:"BIND" default:":80"`
}

func main() {
	if err := envconf.Parse(&config); err != nil {
		log.Fatal(err.Error())
	}

	// Wrapper so that defer calls in serve will run
	if err := serve(); err != nil {
		log.Fatal(err.Error())
	}
}

func serve() error {
	router := api.BuildRouter(api.Deps{})
	return http.ListenAndServe(config.Bind, router)
}
