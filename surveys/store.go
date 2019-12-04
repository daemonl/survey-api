package surveys

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	client *mongo.Client
}

func NewStore(client *mongo.Client) *Store {
	return &Store{
		client: client,
	}
}

var todo = errors.New("TODO")

func (s *Store) AddSurveyResponse(entry Response) (*StoredResponse, error) {
	return nil, todo
}

func (s *Store) GetSurveyResponse(id string) (*StoredResponse, error) {
	return nil, todo

}

type Stats struct{}

func (s *Store) GetStats() (*Stats, error) {
	return nil, todo
}
