package surveys

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pborman/uuid"
)

type S3Store struct {
	service *s3.S3
	bucket  string
}

func NewS3Store(service *s3.S3, bucket string) *S3Store {
	return &S3Store{
		service: service,
		bucket:  bucket,
	}
}

func (s *S3Store) AddSurveyResponse(ctx context.Context, entry Response) (*StoredResponse, error) {
	id := uuid.New()
	stored := &StoredResponse{
		Response: entry,
		ID:       id,
	}

	data, err := json.Marshal(stored)
	if err != nil {
		return nil, err
	}

	_, err = s.service.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fmt.Sprintf("%s/response.json", id)),
		Body:   bytes.NewReader(data),
	})

	if err != nil {
		return nil, err
	}

	return stored, nil
}

func (s *S3Store) GetSurveyResponse(ctx context.Context, id string) (*StoredResponse, error) {
	s3Response, err := s.service.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fmt.Sprintf("%s/response.json", id)),
	})
	if err != nil {
		if awserr, ok := err.(awserr.Error); ok && awserr.Code() == s3.ErrCodeNoSuchKey {
			return nil, NotFoundError
		}
		return nil, err
	}

	resp := &StoredResponse{}
	return resp, json.NewDecoder(s3Response.Body).Decode(resp)
}

func (s *S3Store) GetStats(ctx context.Context) (*Stats, error) {
	return nil, fmt.Errorf("Stats are not supported for S3")
}
