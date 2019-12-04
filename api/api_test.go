package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/daemonl/registerapi/surveys"
	"github.com/stretchr/testify/mock"
)

type MockStore struct {
	mock.Mock
}

func (ms *MockStore) AddSurveyResponse(ctx context.Context, resp surveys.Response) (*surveys.StoredResponse, error) {
	args := ms.Called(resp)
	if err := args.Error(1); err != nil {
		return nil, err
	}
	return args.Get(0).(*surveys.StoredResponse), nil
}

func TestRouter(t *testing.T) {
	deps := &Deps{}
	surveyStore := &MockStore{}
	surveyStore.Mock.Test(t)
	deps.SurveyStore = surveyStore
	router := BuildRouter(deps)

	reset := func() {
		surveyStore.Mock = mock.Mock{}
	}

	do := func(method, path string, body interface{}) *httptest.ResponseRecorder {
		rw := httptest.NewRecorder()
		var bodyReader io.Reader = nil
		if body != nil {
			bodyBytes, err := json.Marshal(body)
			if err != nil {
				t.Fatal(err.Error())
			}
			bodyReader = bytes.NewReader(bodyBytes)
		}

		req := httptest.NewRequest(method, path, bodyReader)
		router.ServeHTTP(rw, req)
		return rw
	}

	if res := do("GET", "/up", nil); res.Code != 200 {
		t.Fatalf("/up status %d", res.Code)
	}

	if res := do("GET", "/frobnork", nil); res.Code != 404 {
		t.Errorf("/frobnork status %d, looking for not found", res.Code)
	}

	t.Run("Add Response", func(t *testing.T) {
		defer surveyStore.AssertExpectations(t)
		defer reset()

		surveyStore.On("AddSurveyResponse", mock.Anything).
			Once().
			Return(&surveys.StoredResponse{
				ID: "storedID",
				Response: surveys.Response{
					Age: 10,
				},
			}, nil)

		res := do("POST", "/responses", surveys.Response{
			Age: 10,
			Animals: map[string]surveys.AnimalResponse{
				"dog": {Rating: 10, Owned: 50},
			},
		})
		if res.Code != 200 {
			t.Fatalf("POST /responses status %d", res.Code)
		}
		gotRes := &surveys.StoredResponse{}
		if err := json.NewDecoder(res.Body).Decode(gotRes); err != nil {
			t.Fatal(err.Error())
		}
		if gotRes.ID != "storedID" {
			t.Errorf("ID: %s", gotRes.ID)
		}
		if gotRes.Age != 10 {
			t.Errorf("Age: %d", gotRes.Age)
		}
	})

	t.Run("Validation Error", func(t *testing.T) {
		defer surveyStore.AssertExpectations(t)
		defer reset()

		surveyStore.On("AddSurveyResponse", mock.Anything).
			Once().
			Return(&surveys.StoredResponse{
				ID: "storedID",
				Response: surveys.Response{
					Age: 10,
				},
			}, nil)

		res := do("POST", "/responses", surveys.Response{Age: -1})
		if res.Code != 400 {
			t.Fatalf("POST /responses status %d", res.Code)
		}

		gotRes := map[string]string{}
		if err := json.NewDecoder(res.Body).Decode(&gotRes); err != nil {
			t.Fatal(err.Error())
		}
		if _, ok := gotRes["age"]; !ok {
			t.Errorf("Expecting an error for age, got none")
		}
	})
}
