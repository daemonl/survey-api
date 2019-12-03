package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Deps struct{}

func BuildRouter(deps Deps) http.Handler {
	r := mux.NewRouter()
	r.Methods("GET").Path("/up").Handler(JSONWrap(upHandler))
	return r
}

func upHandler(req *http.Request) (interface{}, error) {
	return map[string]interface{}{
		"status": "OK",
	}, nil
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
