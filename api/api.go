package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/daemonl/registerapi/surveys"
	"github.com/gorilla/mux"
)

type Deps struct {
	SurveyStore interface {
		AddSurveyResponse(context.Context, surveys.Response) (*surveys.StoredResponse, error)
		GetSurveyResponse(context.Context, string) (*surveys.StoredResponse, error)
		GetStats(context.Context) (*surveys.Stats, error)
	}
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
		return responseStore.GetSurveyResponse(req.Context(), responseID)
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
