package main

import (
	"encoding/json"
	"net/http"
)

func ErrorBadRequest(message string, w http.ResponseWriter) {
	dat, err := errorHelper(message)
	if err != nil {
		ErrorServer("Something Went Wrong", w)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(400)
	w.Write(dat)
}

func ErrorUnauthorized(message string, w http.ResponseWriter) {
	dat, err := errorHelper(message)
	if err != nil {
		ErrorServer("Something Went Wrong", w)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(401)
	w.Write(dat)
}

func ErrorForbidden(message string, w http.ResponseWriter) {
	dat, err := errorHelper(message)
	if err != nil {
		ErrorServer("Something Went Wrong", w)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(403)
	w.Write(dat)
}

func ErrorNotFound(message string, w http.ResponseWriter) {
	dat, err := errorHelper(message)
	if err != nil {
		ErrorServer("Something Went Wrong", w)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(404)
	w.Write(dat)
}

func ErrorServer(message string, w http.ResponseWriter) {
	dat, err := errorHelper(message)
	if err != nil {
		ErrorServer("Something Went Wrong", w)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(500)
	w.Write(dat)
}

func errorHelper(message string) ([]byte, error) {
	type returnVal struct {
		Error string `json:"error"`
	}
	responseBody := returnVal{Error: message}
	dat, err := json.Marshal(responseBody)
	if err != nil {
		return nil, err
	}
	return dat, nil
}
