package surveys

import (
	"context"
	"errors"

	"github.com/pborman/uuid"

	"go.mongodb.org/mongo-driver/bson"
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

func (s *Store) db() *mongo.Database {
	return s.client.Database(s.dbName)
}

func (s *Store) AddSurveyResponse(ctx context.Context, entry Response) (*StoredResponse, error) {
	id := uuid.New()
	stored := &StoredResponse{
		Response: entry,
		ID:       id,
	}

	_, err := s.db().Collection("surveys").InsertOne(ctx, stored)
	if err != nil {
		return nil, err
	}

	return stored, nil
}

func (s *Store) GetSurveyResponse(ctx context.Context, id string) (*StoredResponse, error) {
	row := s.db().Collection("surveys").FindOne(ctx, bson.M{
		"_id": id,
	})
	if err := row.Err(); err != nil {
		return nil, err
	}
	resp := &StoredResponse{}
	return resp, row.Decode(resp)
}

type Stats struct {
	Count int64 `json:"count"`
}

func (s *Store) GetStats(ctx context.Context) (*Stats, error) {
	// Could be more interesting results like how many people like dogs
	stats := &Stats{}
	var err error
	stats.Count, err = s.db().Collection("surveys").CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	return stats, nil
}
