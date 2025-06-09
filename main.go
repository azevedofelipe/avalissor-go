package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/azevedofelipe/avalissor-go/internal/auth"
	"github.com/azevedofelipe/avalissor-go/internal/database"
	"github.com/google/uuid"
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
	mux.HandleFunc("POST /api/colleges", apiCfg.handlerCreateCollege)
	mux.HandleFunc("GET /api/colleges/{collegeID}", apiCfg.handlerGetColleges)
	mux.HandleFunc("DELETE /api/colleges/{collegeID}", apiCfg.handlerDeleteCollegeID)

	mux.HandleFunc("POST /api/colleges/{collegeID}/campuses", apiCfg.handlerCreateCampus)

	mux.HandleFunc("POST /admin/reset", apiCfg.handlerResetUsers)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

// Validate JWT from Header
func (cfg *apiConfig) AuthorizeHeader(header http.Header) (userID uuid.UUID, err error) {
	tokenString, err := auth.GetBearerToken(header)
	if err != nil {
		return
	}

	userID, err = auth.ValidateJWT(tokenString, cfg.tokenSecret)
	if err != nil {
		return
	}

	return
}
