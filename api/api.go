package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/daemonl/survey-api/surveys"
	"github.com/gorilla/mux"
)

type SurveyStore interface {
	AddSurveyResponse(context.Context, surveys.Response) (*surveys.StoredResponse, error)
	GetSurveyResponse(context.Context, string) (*surveys.StoredResponse, error)
	GetStats(context.Context) (*surveys.Stats, error)
}

type Deps struct {
	SurveyStore SurveyStore
}

func BuildRouter(deps *Deps) http.Handler {
	r := mux.NewRouter()

	r.Methods("GET").Path("/up").Handler(JSONWrap(upHandler))

	r.Methods("POST").Path("/responses").Handler(
		JSONWrap(
			buildAddResponseHandler(deps.SurveyStore),
		),
	)

	r.Methods("GET").Path("/responses/{id}").Handler(
		JSONWrap(
			buildGetResponseHandler(deps.SurveyStore),
		),
	)

	r.Methods("GET").Path("/stats").Handler(
		JSONWrap(
			buildGetStatsHandler(deps.SurveyStore),
		),
	)

	r.NotFoundHandler = JSONWrap(func(req *http.Request) (interface{}, error) {
		return nil, simpleError(404, "Not Found")
	})

	r.Use(requestLogger)

	return r
}

func parseRequest(req *http.Request, into interface{}) error {
	err := json.NewDecoder(req.Body).Decode(into)
	if err == nil {
		return nil
	}
	if err == io.EOF {
		return simpleError(400, "Body is required")
	}
	return err
}

func buildAddResponseHandler(responseStore interface {
	AddSurveyResponse(context.Context, surveys.Response) (*surveys.StoredResponse, error)
}) func(req *http.Request) (interface{}, error) {
	return func(req *http.Request) (interface{}, error) {
		surveyResponse := surveys.Response{}
		if err := parseRequest(req, &surveyResponse); err != nil {
			return nil, err
		}

		if validationResponse := surveyResponse.Validate(); validationResponse != nil {
			return nil, simpleErrorResponse{
				code: 400,
				body: validationResponse,
			}
		}

		return responseStore.AddSurveyResponse(req.Context(), surveyResponse)
	}
}

func buildGetResponseHandler(responseStore interface {
	GetSurveyResponse(context.Context, string) (*surveys.StoredResponse, error)
}) func(req *http.Request) (interface{}, error) {

	return func(req *http.Request) (interface{}, error) {
		responseID := mux.Vars(req)["id"]
		res, err := responseStore.GetSurveyResponse(req.Context(), responseID)
		if err == surveys.NotFoundError {
			return nil, simpleError(404, "Response Not Found")
		} else if err != nil {
			return nil, err

		}
		return res, nil
	}
}

func buildGetStatsHandler(responseStore interface {
	GetStats(context.Context) (*surveys.Stats, error)
}) func(req *http.Request) (interface{}, error) {

	return func(req *http.Request) (interface{}, error) {
		return responseStore.GetStats(req.Context())
	}
}

func upHandler(req *http.Request) (interface{}, error) {
	return map[string]interface{}{
		"status": "OK",
	}, nil
}
