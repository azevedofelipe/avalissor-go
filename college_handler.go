package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/azevedofelipe/avalissor-go/internal/database"
)

type College struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (cfg *apiConfig) handlerGetColleges(w http.ResponseWriter, r *http.Request) {
	collegeIdString := r.PathValue("collegeID")

	var colleges []database.College
	var err error

	if collegeIdString != "" {
		collegeId, err := strconv.Atoi(collegeIdString)
		if err != nil {
			w.WriteHeader(500)
			log.Printf("Erro convertendo collegeID: %v", err)
			return
		}

		college, err := cfg.queries.GetCollegeByID(r.Context(), int32(collegeId))
		if err != nil {
			http.Error(w, "Erro obtendo faculdades", 500)
			log.Printf("Erro buscando faculdades: %v", err)
			return
		}

		colleges = []database.College{college}
		log.Printf("Getting single college")

	} else {
		colleges, err = cfg.queries.GetColleges(r.Context())
		if err != nil {
			http.Error(w, "Erro obtendo faculdades", 500)
			log.Printf("Erro buscando faculdades: %v", err)
			return
		}
		log.Printf("Getting all colleges")
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
	userID, err := cfg.AuthorizeHeader(r.Header)
	if err != nil {
		http.Error(w, "Error authorizing header", http.StatusUnauthorized)
		log.Printf("Error authorizing header: %v", err)
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

func (cfg *apiConfig) handlerDeleteCollegeID(w http.ResponseWriter, r *http.Request) {
	collegeIdString := r.PathValue("collegeID")

	if collegeIdString == "" {
		http.Error(w, "Error authorizing header", http.StatusBadRequest)
		log.Printf("No collegeID passed")
		return
	}

	collegeId, err := strconv.Atoi(collegeIdString)
	if err != nil {
		http.Error(w, "Invalid collegeID", http.StatusBadRequest)
		log.Printf("Error converting string: %v", err)
		return
	}

	userID, err := cfg.AuthorizeHeader(r.Header)
	if err != nil {
		http.Error(w, "Error authorizing header", http.StatusUnauthorized)
		log.Printf("Error authorizing header: %v", err)
		return
	}

	log.Printf("UserID %s Deleting college %d", userID, collegeId)

	err = cfg.queries.DeleteCollegeID(r.Context(), int32(collegeId))
	if err != nil {
		http.Error(w, "Error deleting college", http.StatusBadRequest)
		log.Printf("Error deleting college %d: %v", collegeId, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("Successfuly deleted college %d", collegeId)))
}
