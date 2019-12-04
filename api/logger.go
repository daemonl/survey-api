package api

import (
	"log"
	"net/http"
	"time"
)

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		begin := time.Now()
		recorder := &responseWriterStatusRecorder{
			status:         200,
			ResponseWriter: rw,
		}
		next.ServeHTTP(recorder, req)
		log.Printf("%s %s => %d in %f seconds", req.Method, req.URL.Path, recorder.status, time.Since(begin).Seconds())

	})
}

type responseWriterStatusRecorder struct {
	status int
	http.ResponseWriter
}
