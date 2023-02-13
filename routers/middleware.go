package routers

import (
	"log"
	"net/http"
	"time"
)

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		observer := &responseObserver{ResponseWriter: w}
		next.ServeHTTP(observer, r)
		elapsed := time.Since(start)
		log.Printf("%d %s: %s in %s", observer.status, r.Method, r.URL.Path, elapsed)
	})
}

type responseObserver struct {
	http.ResponseWriter
	status int
}

func (o *responseObserver) WriteHeader(code int) {
	o.ResponseWriter.WriteHeader(code)
	o.status = code
}
