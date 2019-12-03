package api

import (
	"io"
	"net/http/httptest"
	"testing"
)

func TestRouter(t *testing.T) {
	router := BuildRouter(Deps{})

	do := func(method, path string, body io.Reader) *httptest.ResponseRecorder {
		rw := httptest.NewRecorder()
		req := httptest.NewRequest(method, path, body)
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
