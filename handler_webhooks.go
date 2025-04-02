package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/rrochlin/WebServerGo/internal/auth"
)

func (cfg *apiConfig) HandlerUpgradeUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}
	key, err := auth.GetAPIKey(req.Header)
	if err != nil {
		ErrorUnauthorized(err.Error(), w)
		return
	}

	if key != cfg.api.polkaKey {
		ErrorUnauthorized("Invalid Api key", w)
		return
	}

	var params parameters
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&params)
	if err != nil {
		ErrorBadRequest(err.Error(), w)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}
	userid, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		ErrorBadRequest(err.Error(), w)
		return
	}

	_, err = cfg.db.query.UpgradeUser(req.Context(), userid)
	if err != nil {
		ErrorNotFound(err.Error(), w)
		return
	}
	w.WriteHeader(204)

}
