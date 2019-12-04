package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func simpleError(code int, text string) HTTPErrorResponse {
	return simpleErrorResponse{
		code: code,
		body: map[string]interface{}{
			"status": text,
		},
	}
}

type simpleErrorResponse struct {
	code int
	body interface{}
}

func (sr simpleErrorResponse) HTTPStatus() int {
	return sr.code
}
func (sr simpleErrorResponse) HTTPBody() interface{} {
	return sr.body
}

func (sr simpleErrorResponse) Error() string {
	return fmt.Sprintf("HTTP %d", sr.code)
}

type HTTPErrorResponse interface {
	error
	HTTPResponse
}

type HTTPResponse interface {
	HTTPStatus() int
	HTTPBody() interface{}
}

func JSONWrap(handler func(req *http.Request) (interface{}, error)) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		responseObject, err := handler(req)
		if err != nil {
			if httpResp, ok := err.(HTTPResponse); ok {
				doJSONResponse(rw, httpResp.HTTPStatus(), httpResp.HTTPBody())
				return
			}
			log.Printf("Unhandled: %s", err.Error())
			doJSONResponse(rw, 500, map[string]interface{}{
				"error": "Internal Server Error",
			})
			return
		}
		if httpResp, ok := responseObject.(HTTPResponse); ok {
			doJSONResponse(rw, httpResp.HTTPStatus(), httpResp.HTTPBody())
			return
		}
		doJSONResponse(rw, 200, responseObject)
	})
}

func doJSONResponse(rw http.ResponseWriter, status int, body interface{}) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	json.NewEncoder(rw).Encode(body)
	return
}
