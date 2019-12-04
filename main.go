package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/daemonl/survey-api/api"
	"github.com/daemonl/survey-api/awssecret"
	"github.com/daemonl/survey-api/surveys"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"gopkg.daemonl.com/envconf"
)

var config struct {
	Bind string `env:"BIND" default:":80"`

	DataStore string `env:"DATA_STORE_URL" default:"mongodb://localhost:27017/surveys"`
}

func main() {
	awssecret.Default()
	if err := envconf.Parse(&config); err != nil {
		log.Fatal(err.Error())
	}

	// Wrapper so that defer calls in serve will run
	if err := serve(); err != nil {
		log.Fatal(err.Error())
	}
}

func serve() error {
	dataStoreURL, err := url.Parse(strings.Split(config.DataStore, ",")[0])
	if err != nil {
		return fmt.Errorf("Invalid DATA_STORE_URL: %s", config.DataStore)
	}

	var surveyStore api.SurveyStore

	switch dataStoreURL.Scheme {
	case "mongodb":
		client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.DataStore))
		if err != nil {
			return err
		}
		defer client.Disconnect(context.TODO())

		surveyStore = surveys.NewMongoStore(client, dataStoreURL.Path[1:]) // Drops the `/`

	case "s3":
		if dataStoreURL.Path != "/" {
			return fmt.Errorf("No path prefix is supported for S3")
		}
		sess, err := session.NewSession()
		if err != nil {
			return err
		}
		service := s3.New(sess)
		log.Printf("S3 Store in %s", dataStoreURL.Host)
		surveyStore = surveys.NewS3Store(service, dataStoreURL.Host)

	default:
		return fmt.Errorf("Must set a DATA_STORE_URL of either mongodb or s3")
	}

	router := api.BuildRouter(&api.Deps{
		SurveyStore: surveyStore,
	})
	log.Printf("Listening, Bound to %s", config.Bind)
	return http.ListenAndServe(config.Bind, router)
}
