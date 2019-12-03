package api

import (
	"net/http/httptest"
	"testing"
)

func TestRouter(t *testing.T) {
	router := BuildRouter(Deps{})
	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/up", nil)
	router.ServeHTTP(rw, req)

	if rw.Code != 200 {
		t.Fatalf("/up status %d", rw.Code)
	}
}
