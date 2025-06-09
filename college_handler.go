package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/azevedofelipe/avalissor-go/internal/auth"
	"github.com/azevedofelipe/avalissor-go/internal/database"
)

type College struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (cfg *apiConfig) handlerGetColleges(w http.ResponseWriter, r *http.Request) {

	colleges, err := cfg.queries.GetColleges(r.Context())
	if err != nil {
		http.Error(w, "Erro obtendo faculdades", 500)
		log.Printf("Erro buscando faculdades: %v", err)
		return
	}

	response := make([]College, len(colleges))
	for idx, college := range colleges {
		response[idx] = College{
			ID:   int(college.ID),
			Name: college.NameCollege,
		}
	}

	dat, err := json.Marshal(&response)
	if err != nil {
		http.Error(w, "Erro gerando resposta", 500)
		log.Printf("Error marshalling response: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)

}

func (cfg *apiConfig) handlerCreateCollege(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, "Unauthorized, please provide valid authorization", http.StatusUnauthorized)
		log.Printf("Error getting token: %v", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.tokenSecret)
	if err != nil {
		http.Error(w, "Unauthorized, please provide valid authorization", http.StatusUnauthorized)
		log.Printf("Error validating JWT: %v", err)
		return
	}

	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)

	college, err := cfg.queries.CreateCollege(r.Context(), database.CreateCollegeParams{
		NameCollege: params.Name,
		CreatedBy:   userID,
	})

	if err != nil {
		http.Error(w, "Erro criando faculdade", 500)
		log.Printf("Erro criando faculdade %s: %v", params.Name, err)
		return
	}

	response := College{
		ID:   int(college.ID),
		Name: college.NameCollege,
	}

	dat, err := json.Marshal(&response)
	if err != nil {
		http.Error(w, "Erro gerando resposta", 500)
		log.Printf("Error marshalling response: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)

}

func (cfg *apiConfig) handlerGetCollegeID(w http.ResponseWriter, r *http.Request) {
	college_id := r.PathValue("chirpID")

	colleges, err := cfg.queries.GetColleges(r.Context())
	if err != nil {
		http.Error(w, "Erro obtendo faculdades", 500)
		log.Printf("Erro buscando faculdades: %v", err)
		return
	}

	response := make([]College, len(colleges))
	for idx, college := range colleges {
		response[idx] = College{
			ID:   int(college.ID),
			Name: college.NameCollege,
		}
	}

	dat, err := json.Marshal(&response)
	if err != nil {
		http.Error(w, "Erro gerando resposta", 500)
		log.Printf("Error marshalling response: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)

}
