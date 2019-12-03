package surveys

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	mgo mongo.Client
}

var todo = errors.New("TODO")

func (s *Store) AddSurveyResponse(entry Response) error {
	return todo
}

func (s *Store) GetSurveyResponse(id string) (*Response, error) {
	return nil, todo

}

type Stats struct{}

func (s *Store) GetStats() (*Stats, error) {
	return nil, todo
}
