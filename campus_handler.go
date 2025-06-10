package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/azevedofelipe/avalissor-go/internal/database"
)

type Campus struct {
	ID        int32     `json:"id"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	CollegeID int32     `json:"college_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (cfg *apiConfig) handlerCreateCampus(w http.ResponseWriter, r *http.Request) {
	collegeIdString := r.PathValue("collegeID")

	if collegeIdString == "" {
		http.Error(w, "Nenhum collegeID passado para POST", http.StatusBadRequest)
		log.Printf("No collegeID passed")
		return
	}

	collegeId, err := strconv.Atoi(collegeIdString)
	if err != nil {
		http.Error(w, "Invalid collegeID", http.StatusBadRequest)
		log.Printf("Error converting string: %v", err)
		return
	}

	_, err = cfg.AuthorizeHeader(r.Header)
	if err != nil {
		http.Error(w, "Error authorizing header", http.StatusUnauthorized)
		log.Printf("Error authorizing header: %v", err)
		return
	}

	type parameters struct {
		Name     string `json:"campus_name"`
		Location string `json:"location"`
	}

	params, err := decode[parameters](r)
	if err != nil {
		http.Error(w, "Error reading parameters", 400)
		log.Printf("Error decoding params: %v", err)
		return
	}

	campus, err := cfg.queries.CreateCampus(r.Context(), database.CreateCampusParams{
		Name: params.Name,
		Location: sql.NullString{
			Valid:  params.Location != "",
			String: params.Location,
		},
		CollegeID: int32(collegeId),
	})

	response := Campus{
		ID:        campus.ID,
		Name:      campus.Name,
		Location:  campus.Location.String,
		CollegeID: campus.CollegeID,
		CreatedAt: campus.CreatedAt,
		UpdatedAt: campus.UpdatedAt,
	}

	dat, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Erro gerando resposta", 500)
		log.Printf("Error marshalling response: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)

}

func (cfg *apiConfig) handlerDeleteCampusID(w http.ResponseWriter, r *http.Request) {
	campusIdString := r.PathValue("campusID")

	if campusIdString == "" {
		http.Error(w, "Nenhum collegeID passado para POST", http.StatusBadRequest)
		log.Printf("No collegeID passed")
		return
	}

	campusID, err := strconv.Atoi(campusIdString)
	if err != nil {
		http.Error(w, "Invalid collegeID", http.StatusBadRequest)
		log.Printf("Error converting string: %v", err)
		return
	}

	userID, err := cfg.AuthorizeHeader(r.Header)
	if err != nil {
		http.Error(w, "Invalid header", http.StatusBadRequest)
		log.Printf("Error validating header: %v", err)
		return
	}

	log.Printf("User %s deleting campus %d", userID, campusID)
	err = cfg.queries.DeleteCampus(r.Context(), int32(campusID))
	if err != nil {
		http.Error(w, "Unable to delete campus", 500)
		log.Printf("Error deleting campus %d: %v", campusID, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("Succesfully deleted campus %d", campusID)))
}
