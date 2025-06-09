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
	queries     *database.Queries
	tokenSecret string
}

func main() {
	const port = "8080"

	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	tokenSecret := os.Getenv("TOKEN_SECRET")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error opening database: %s", err)
	}
	dbQueries := database.New(db)

	apiCfg := apiConfig{
		queries:     dbQueries,
		tokenSecret: tokenSecret,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/users", apiCfg.handlerUserCreation)
	mux.HandleFunc("POST /api/login", apiCfg.handlerUserLogin)

	mux.HandleFunc("GET /api/colleges", apiCfg.handlerGetColleges)
	mux.HandleFunc("GET /api/colleges/{collegeID}", apiCfg.handlerGetColleges)
	mux.HandleFunc("POST /api/colleges", apiCfg.handlerCreateCollege)

	mux.HandleFunc("POST /admin/reset", apiCfg.handlerResetUsers)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
