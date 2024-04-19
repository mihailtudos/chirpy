package middleware

import (
	"fmt"
	"net/http"
	"time"
)

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%v: %s %s\n", time.Now().UTC(), r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
