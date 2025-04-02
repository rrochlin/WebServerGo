package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rrochlin/WebServerGo/internal/auth"
)

func (cfg *apiConfig) HandlerRefresh(w http.ResponseWriter, req *http.Request) {
	untrustedToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		ErrorBadRequest(err.Error(), w)
		return
	}
	token, err := cfg.db.query.GetRToken(req.Context(), untrustedToken)
	if err != nil {
		ErrorUnauthorized(err.Error(), w)
		return
	}
	if token.RevokedAt.Valid {
		ErrorUnauthorized(fmt.Sprintf("Refresh Token has been revoked at %v", token.RevokedAt), w)
		return
	}

	jtoken, err := auth.MakeJWT(token.UserID, cfg.api.secret)
	if err != nil {
		ErrorServer(err.Error(), w)
		return
	}
	type response struct {
		Token string `json:"token"`
	}
	res := response{Token: jtoken}

	dat, err := json.Marshal(res)
	if err != nil {
		ErrorServer(fmt.Sprintf("failed to encode token for response %v", err), w)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)

}

func (cfg *apiConfig) HandlerRevoke(w http.ResponseWriter, req *http.Request) {
	untrustedToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		ErrorBadRequest(err.Error(), w)
		return
	}
	token, err := cfg.db.query.RevokeToken(req.Context(), untrustedToken)
	if err != nil {
		ErrorUnauthorized(err.Error(), w)
		return
	}
	if !token.RevokedAt.Valid {
		ErrorServer("Token Revocation Failed", w)
		return
	}
	w.WriteHeader(204)

}
