package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	"github.com/daemonl/survey-api/surveys"
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

func (ms *MockStore) GetSurveyResponse(ctx context.Context, id string) (*surveys.StoredResponse, error) {
	args := ms.Called(id)
	if err := args.Error(1); err != nil {
		return nil, err
	}
	return args.Get(0).(*surveys.StoredResponse), nil
}

func (ms *MockStore) GetStats(ctx context.Context) (*surveys.Stats, error) {
	args := ms.Called()
	if err := args.Error(1); err != nil {
		return nil, err
	}
	return args.Get(0).(*surveys.Stats), nil
}

func TestRouter(t *testing.T) {
	deps := &Deps{}
	router := BuildRouter(deps)

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
}

func jsonRead(obj interface{}) io.Reader {
	jsonBytes, _ := json.Marshal(obj)
	return bytes.NewReader(jsonBytes)
}

func jsonConvertType(in interface{}, out interface{}) {
	resBody, err := json.Marshal(in)
	if err != nil {
		panic(err.Error())
	}
	if err := json.Unmarshal(resBody, out); err != nil {
		panic(err.Error())
	}
}

func jsonMapType(in interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	jsonConvertType(in, &out)
	return out
}

func TestAddResponseHappy(t *testing.T) {
	surveyStore := &MockStore{}
	surveyStore.Mock.Test(t)
	defer surveyStore.AssertExpectations(t)

	surveyStore.On("AddSurveyResponse", mock.Anything).
		Once().
		Return(&surveys.StoredResponse{
			ID: "storedID",
			Response: surveys.Response{
				Age: 10,
			},
		}, nil)

	handler := buildAddResponseHandler(surveyStore)

	gotResRaw, err := handler(httptest.NewRequest("POST", "/responses", jsonRead(
		surveys.Response{
			Age: 10,
			Animals: map[string]surveys.AnimalResponse{
				"dog": {Rating: 10, Owned: 50},
			},
		})))
	if err != nil {
		t.Fatal(err.Error())
	}

	gotRes := jsonMapType(gotResRaw)

	if gotRes["id"] != "storedID" {
		t.Errorf("ID: %s", gotRes["id"])
	}
	if gotRes["age"] != 10.0 {
		t.Errorf("Age: %d", gotRes["age"])
	}
}

func TestAddResponseInvalid(t *testing.T) {

	handler := buildAddResponseHandler(nil)

	_, err := handler(httptest.NewRequest("POST", "/responses", jsonRead(
		surveys.Response{Age: -1}),
	))
	if err == nil {
		t.Fatal("Expected an error")
	}

	httpResp, ok := err.(HTTPResponse)
	if !ok {
		t.Fatalf("Bad error type: %T", err)
	}

	if httpResp.HTTPStatus() != 400 {
		t.Fatalf("POST /responses status %d", httpResp.HTTPStatus())
	}

	gotRes := jsonMapType(httpResp.HTTPBody())
	if _, ok := gotRes["age"]; !ok {
		t.Errorf("Expecting an error for age, got none: %#v", gotRes)
	}
}

func TestGetResponseHappy(t *testing.T) {
	surveyStore := &MockStore{}
	surveyStore.Mock.Test(t)
	defer surveyStore.AssertExpectations(t)

	surveyStore.On("GetSurveyResponse", "12345").
		Once().
		Return(&surveys.StoredResponse{
			ID: "12345",
			Response: surveys.Response{
				Age: 10,
			},
		}, nil)

	handler := buildGetResponseHandler(surveyStore)
	req := httptest.NewRequest("GET", "/responses/12345", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "12345"})
	gotResRaw, err := handler(req)
	if err != nil {
		t.Fatal(err.Error())
	}

	gotRes := jsonMapType(gotResRaw)

	if gotRes["id"] != "12345" {
		t.Errorf("ID: %s", gotRes["id"])
	}
	if gotRes["age"] != 10.0 {
		t.Errorf("Age: %d", gotRes["age"])
	}
}

func TestGetResponseNotFound(t *testing.T) {
	surveyStore := &MockStore{}
	surveyStore.Mock.Test(t)
	defer surveyStore.AssertExpectations(t)

	surveyStore.On("GetSurveyResponse", "12345").
		Once().
		Return(nil, surveys.NotFoundError)

	handler := buildGetResponseHandler(surveyStore)
	req := httptest.NewRequest("GET", "/responses/12345", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "12345"})

	_, err := handler(req)
	if err == nil {
		t.Fatal("Expected an error")
	}

	httpResp, ok := err.(HTTPResponse)
	if !ok {
		t.Fatalf("Bad error type: %T, %s", err, err.Error())
	}

	if httpResp.HTTPStatus() != 404 {
		t.Fatalf("POST /responses status %d", httpResp.HTTPStatus())
	}

	gotRes := jsonMapType(httpResp.HTTPBody())
	if gotRes["status"] != "Response Not Found" {
		t.Errorf("Expecting Not Found error, got: %#v", gotRes)
	}
}

func TestGetStatsHappy(t *testing.T) {
	surveyStore := &MockStore{}
	surveyStore.Mock.Test(t)
	defer surveyStore.AssertExpectations(t)

	surveyStore.On("GetStats").
		Once().
		Return(&surveys.Stats{
			Count: 123,
		}, nil)

	handler := buildGetStatsHandler(surveyStore)
	req := httptest.NewRequest("GET", "/stats", nil)
	gotResRaw, err := handler(req)
	if err != nil {
		t.Fatal(err.Error())
	}

	gotRes := jsonMapType(gotResRaw)

	if gotRes["count"] != 123.0 {
		t.Errorf("Count: %s", gotRes["count"])
	}
}
