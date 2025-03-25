package main

import "net/http"

func main() {
	var mu = http.NewServeMux()
	var server = http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mu,
	}
	server.ListenAndServe()

}
