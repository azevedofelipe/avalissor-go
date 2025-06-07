package main

import (
	"log"
	"net/http"
)

type apiConfig struct {
	queries     *database.Queries
	tokenSecret string
}

func main() {
	const port = "8080"

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/login", apiCfg.handlerUserLogin)
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Fatal(srv.ListenAndServe())
}
