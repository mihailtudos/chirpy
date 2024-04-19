package main

import (
	"log"
	"net/http"

	"github.com/mihailtudos/chirpy/internal/database"
	"github.com/mihailtudos/chirpy/middleware"
)

type apiConfig struct {
	fileserverHits int
	db             *database.DB
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	db, err := database.NewDB("./database.json")
	if err != nil {
		log.Fatal("db cannot be created", err)
	}

	apiCfg := apiConfig{
		fileserverHits: 0,
		db:             db,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /admin/metrics/", apiCfg.handlerMetrics)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirps)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)

	corsMux := middleware.LogRequest(middlewareCors(mux))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
