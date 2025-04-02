package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rrochlin/WebServerGo/internal/database"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("failed to connect to database %v", err)
		return
	}

	dbQueries := database.New(db)
	mux := http.NewServeMux()
	cfg := apiConfig{
		api: apiSettings{
			fileserverHits: atomic.Int32{},
			platform:       os.Getenv("PLATFORM"),
			host:           os.Getenv("HOST_URL"),
			secret:         os.Getenv("SECRET"),
			polkaKey:       os.Getenv("POLKA_API_KEY"),
		},
		db: dbSettings{
			query: dbQueries,
		},
	}

	mux.Handle("/app/", cfg.middlewareMetricsInc(handler))
	mux.HandleFunc("GET /api/healthz", HandlerHealthz)
	mux.HandleFunc("GET /admin/metrics", cfg.HandlerHits)
	mux.HandleFunc("POST /admin/reset", cfg.HandlerReset)
	mux.HandleFunc("POST /api/chirps", cfg.HandlerChirps)
	mux.HandleFunc("POST /api/users", cfg.HandlerUsers)
	mux.HandleFunc("GET /api/chirps", cfg.HandlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.HandlerGetChirp)
	mux.HandleFunc("POST /api/login", cfg.HandlerLogin)
	mux.HandleFunc("POST /api/refresh", cfg.HandlerRefresh)
	mux.HandleFunc("POST /api/revoke", cfg.HandlerRevoke)
	mux.HandleFunc("PUT /api/users", cfg.HandlerUpdateUser)

	var server = http.Server{
		Addr:    fmt.Sprintf("%v:8080", cfg.api.host),
		Handler: mux,
	}
	server.ListenAndServe()
}

var handler = http.StripPrefix("/app/",
	http.FileServer(http.Dir(".")),
)

// API related configuration
type apiSettings struct {
	fileserverHits atomic.Int32
	platform       string
	host           string
	secret         string
	polkaKey       string
}

// Database related configuration
type dbSettings struct {
	query *database.Queries
}

// Main configuration struct
type apiConfig struct {
	api apiSettings
	db  dbSettings
}
