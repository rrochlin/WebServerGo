package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rrochlin/WebServerGo/internal/auth"
	"github.com/rrochlin/WebServerGo/internal/database"
)

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
	dat, err := json.Marshal(toPublicUser(user, "", ""))
	if err != nil {
		ErrorServer(fmt.Sprintf("Could not convert user to response: %v", err), w)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(dat)
}

func (cfg *apiConfig) HandlerUpdateUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	unverifiedToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		ErrorUnauthorized(err.Error(), w)
		return
	}
	uuid, err := auth.ValidateJWT(unverifiedToken, cfg.api.secret)
	if err != nil {
		ErrorUnauthorized(err.Error(), w)
		return
	}
	var params parameters
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&params)
	if err != nil {
		ErrorServer(err.Error(), w)
		return
	}
	hpass, err := auth.HashPassword(params.Password)
	if err != nil {
		ErrorServer(err.Error(), w)
		return
	}

	user, err := cfg.db.query.UpdateUser(
		req.Context(),
		database.UpdateUserParams{
			Email:          params.Email,
			HashedPassword: hpass,
			ID:             uuid,
		})
	if err != nil {
		ErrorServer(err.Error(), w)
		return
	}
	dat, err := json.Marshal(toPublicUser(user, "", ""))
	if err != nil {
		ErrorServer(err.Error(), w)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)

}
