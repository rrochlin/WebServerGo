package main

import (
	"fmt"
	"net/http"
)

func HandlerHealthz(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) HandlerHits(w http.ResponseWriter, req *http.Request) {
	content := fmt.Sprintf(`<html>
		<body style="background-color:#121212; color:white">
		<h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %v times!</p>
  </body>
</html>`, cfg.api.fileserverHits.Load())
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(content))
}

func (cfg *apiConfig) HandlerReset(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Reset Handler Called")
	if cfg.api.platform != "dev" {
		ErrorForbidden("Cannot reset in Production", w)
		return
	}
	cfg.api.fileserverHits.Store(0)
	cfg.db.query.TruncateUsers(req.Context())
	w.WriteHeader(200)
}
