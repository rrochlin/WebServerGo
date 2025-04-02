package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/rrochlin/WebServerGo/internal/auth"
	"github.com/rrochlin/WebServerGo/internal/database"
)

func (cfg *apiConfig) HandlerGetChirps(w http.ResponseWriter, req *http.Request) {
	s := req.URL.Query().Get("author_id")
	var chirps []database.Chirp
	var err error

	if s == "" {
		chirps, err = cfg.db.query.GetAllChirps(req.Context())
	} else {
		authorId, err := uuid.Parse(s)
		if err != nil {
			ErrorBadRequest(err.Error(), w)
			return
		}
		chirps, err = cfg.db.query.GetAuthorsChirps(req.Context(), authorId)
	}
	sort := req.URL.Query().Get("sort")
	if sort == "desc" {
		slices.SortFunc(chirps, func(a, b database.Chirp) int {
			return a.CreatedAt.Compare(b.CreatedAt) * -1
		})
	}

	if err != nil {
		ErrorServer(fmt.Sprintf("could not fetch chrips: %v", err), w)
		return
	}
	dat, err := json.Marshal(chirps)
	if err != nil {
		ErrorServer(fmt.Sprintf("Could not convert chirps to response: %v", err), w)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)

}

func (cfg *apiConfig) HandlerGetChirp(w http.ResponseWriter, req *http.Request) {
	chirpID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		ErrorBadRequest(fmt.Sprintf("Chirp ID not valid UUID: %v", err), w)
	}
	chirp, err := cfg.db.query.GetChirp(req.Context(), chirpID)
	if err != nil {
		ErrorNotFound(fmt.Sprintf("could not fetch chrip: %v", err), w)
		return
	}
	dat, err := json.Marshal(chirp)
	if err != nil {
		ErrorServer(fmt.Sprintf("Could not convert chirp to response: %v", err), w)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)

}

func (cfg *apiConfig) HandlerPostChirp(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		ErrorUnauthorized(err.Error(), w)
		return
	}
	uuid, err := auth.ValidateJWT(token, cfg.api.secret)
	if err != nil {
		ErrorUnauthorized(err.Error(), w)
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		ErrorBadRequest("failed to parse request body", w)
		return
	}

	if len(params.Body) > 140 {
		ErrorBadRequest("Chirp is too long", w)
		return
	}
	profane := []string{"kerfuffle", "sharbert", "fornax"}
	cleanedBody := ""
	for _, word := range strings.Split(params.Body, " ") {
		if slices.Contains(profane, strings.ToLower(word)) {
			cleanedBody += "****"
		} else {
			cleanedBody += word
		}
		cleanedBody += " "
	}
	cleanedBody = strings.TrimRight(cleanedBody, " ")
	chirp, err := cfg.db.query.CreateChirp(
		req.Context(),
		database.CreateChirpParams{
			Body:   cleanedBody,
			UserID: uuid,
		},
	)
	if err != nil {
		ErrorBadRequest(fmt.Sprintf("failed to create chirp %v", err), w)
		return
	}

	dat, err := json.Marshal(chirp)
	if err != nil {
		ErrorServer(fmt.Sprintf("failed to enocde chirp for response %v", err), w)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(dat)
}

func (cfg *apiConfig) HandlerDeleteChirp(w http.ResponseWriter, req *http.Request) {
	unverifiedToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		ErrorUnauthorized(err.Error(), w)
		return
	}
	userid, err := auth.ValidateJWT(unverifiedToken, cfg.api.secret)
	if err != nil {
		ErrorForbidden(err.Error(), w)
		return
	}

	chirpID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		ErrorBadRequest(fmt.Sprintf("Chirp ID not valid UUID: %v", err), w)
	}

	chirp, err := cfg.db.query.GetChirp(req.Context(), chirpID)
	if err != nil {
		ErrorNotFound(err.Error(), w)
		return
	}

	if chirp.UserID != userid {
		ErrorForbidden("userid does not match chirp owner", w)
		return
	}

	err = cfg.db.query.DeleteChirp(req.Context(), chirp.ID)
	if err != nil {
		ErrorServer(err.Error(), w)
		return
	}

	w.WriteHeader(204)
}
