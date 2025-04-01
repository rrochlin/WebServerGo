package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rrochlin/WebServerGo/internal/auth"
	"github.com/rrochlin/WebServerGo/internal/database"
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

func (cfg *apiConfig) HandlerChirps(w http.ResponseWriter, req *http.Request) {
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

func (cfg *apiConfig) HandlerUsers(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		ErrorBadRequest("failed to parse request body", w)
		return
	}
	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		ErrorServer(fmt.Sprintf("Passowrd hash failed: %v", err), w)
		return
	}

	user, err := cfg.db.query.CreateUser(req.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPass,
	})
	if err != nil {
		ErrorServer(fmt.Sprintf("Could not create user: %v", err), w)
		return
	}
	dat, err := json.Marshal(user)
	if err != nil {
		ErrorServer(fmt.Sprintf("Could not convert user to response: %v", err), w)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(dat)
}

func (cfg *apiConfig) HandlerGetChirps(w http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.db.query.GetAllChirps(req.Context())
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

func (cfg *apiConfig) HandlerLogin(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Expires  *int   `json:"expires_in_seconds,omitempty"`
	}
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		ErrorBadRequest("failed to parse request body", w)
		return
	}

	user, err := cfg.db.query.GetUser(req.Context(), params.Email)
	if err != nil {
		ErrorNotFound(fmt.Sprintf("%v", err), w)
		return
	}
	err = auth.CheckPasswordHash(user.HashedPassword, params.Password)
	if err != nil {
		ErrorUnauthorized("Incorrect Password", w)
		return
	}
	var expirey time.Duration
	if params.Expires == nil {
		expirey, _ = time.ParseDuration("1h")
	} else {
		expirey, err = time.ParseDuration(fmt.Sprintf("%vs", *params.Expires))
		if err != nil {
			ErrorServer(fmt.Sprintf("failed to parse time duration: %v", err), w)
			return
		}
	}

	token, err := auth.MakeJWT(user.ID, cfg.api.secret, expirey)
	if err != nil {
		ErrorServer(fmt.Sprintf("failed to construct JWT %v", err), w)
		return
	}

	dat, err := json.Marshal(toPublicUser(user, token))
	if err != nil {
		ErrorServer(fmt.Sprintf("failed to user for response %v", err), w)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)

}

type PublicUser struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

func toPublicUser(user database.User, token string) PublicUser {
	return PublicUser{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Token:     token,
	}
}
