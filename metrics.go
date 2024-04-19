package main

import (
	"fmt"
	"net/http"
	"os"
)


func (a *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (a *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	file, err := os.ReadFile("./admin/index.html")

	if err != nil {
		fmt.Println(err)
		return
	}
	tmlp := fmt.Sprintf(string(file), a.fileserverHits)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(tmlp))
}
