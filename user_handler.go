package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/azevedofelipe/avalissor-go/internal/auth"
	"github.com/azevedofelipe/avalissor-go/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (cfg *apiConfig) handlerUserCreation(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&params); err != nil {
		log.Printf("Error reading json: %s", err)
		w.WriteHeader(400)
		return
	}

	if params.Username == "" || params.Email == "" || params.Password == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		log.Print("Missing required fields to create account")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error reading json: %s", err)
		w.WriteHeader(500)
		return
	}

	user, err := cfg.queries.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		Username:       params.Username,
		HashedPassword: hashedPassword,
	})

	if err != nil {
		log.Printf("Error creating user: %s", err)
		w.WriteHeader(500)
		return
	}

	response := User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		UpdatedAt: user.UpdatedAt,
		CreatedAt: user.CreatedAt,
	}

	dat, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshalling response json: %s", err)
		w.WriteHeader(500)
		return
	}

	log.Print("Created user successfully")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}

func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&params); err != nil {
		log.Printf("Error reading json: %s", err)
		w.WriteHeader(400)
		return
	}

	user, err := cfg.queries.GetUserByUsername(r.Context(), params.Username)
	if err != nil {
		log.Printf("Error getting user in database: %s", err)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(401)
		w.Write([]byte("Invalid username or password"))
		return
	}

	err = auth.CheckPasswordHash(user.HashedPassword, params.Password)
	if err != nil {
		log.Printf("Passwords dont match: %s", err)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(401)
		w.Write([]byte("Invalid username or password"))
		return
	}

	jwtToken, err := auth.MakeJWT(user.ID, cfg.tokenSecret)
	if err != nil {
		log.Printf("Error making JWT: %s", err)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(500)
		w.Write([]byte("Error creating JWT"))
		return
	}

	type jsonResp struct {
		Token string `json:"token"`
	}

	response := jsonResp{Token: jwtToken}

	dat, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshalling response json: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}

func (cfg *apiConfig) handlerResetUsers(w http.ResponseWriter, r *http.Request) {
	err := cfg.queries.DeleteUsers(r.Context())
	if err != nil {
		log.Printf("Error deleting all users: %s", err)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(401)
		w.Write([]byte("Error deleting all users"))
		return
	}
}
