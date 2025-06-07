package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/azevedofelipe/avalissor-go/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	queries *database.Queries
}

func main() {
	const port = "8080"

	godotenv.Load()

	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error opening database: %s", err)
	}
	dbQueries := database.New(db)
	apiCfg := apiConfig{
		queries: dbQueries,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/users", apiCfg.handlerUserCreation)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
