package surveys

import (
	"context"
	"errors"
	"github.com/pborman/uuid"

	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	client *mongo.Client
	dbName string
}

func NewStore(client *mongo.Client, dbName string) *Store {
	return &Store{
		client: client,
		dbName: dbName,
	}
}

var todo = errors.New("TODO")

func (s *Store) AddSurveyResponse(ctx context.Context, entry Response) (*StoredResponse, error) {
	id := uuid.New()
	stored := &StoredResponse{
		Response: entry,
		ID:       id,
	}

	_, err := s.client.Database(s.dbName).Collection("surveys").InsertOne(ctx, stored)
	if err != nil {
		return nil, err
	}

	return stored, nil
}

func (s *Store) GetSurveyResponse(ctx context.Context, id string) (*StoredResponse, error) {
	return nil, todo

}

type Stats struct{}

func (s *Store) GetStats(ctx context.Context) (*Stats, error) {
	return nil, todo
}
