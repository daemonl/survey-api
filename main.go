package main

import (
	"context"
	"log"
	"net/http"

	"github.com/daemonl/registerapi/api"
	"github.com/daemonl/registerapi/surveys"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"gopkg.daemonl.com/envconf"
)

var config struct {
	Bind      string `env:"BIND" default:":80"`
	MongoURL  string `env:"MONGO_DB_URL" default:"mongodb://localhost:27017"`
	MongoName string `env:"MONGO_DB_NAME" default:"surveys"`
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
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.MongoURL))
	if err != nil {
		return err
	}
	defer client.Disconnect(context.TODO())

	surveyStore := surveys.NewStore(client, config.MongoName)

	router := api.BuildRouter(&api.Deps{
		SurveyStore: surveyStore,
	})
	return http.ListenAndServe(config.Bind, router)
}
