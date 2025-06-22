package handlers

import (
	"encoding/json"
	"net/http"
	"restproject/internal/models"
	"restproject/internal/repository"
)

func AddAdmins(w http.ResponseWriter, r *http.Request) {
	var newAdmins []models.Admin
	err := json.NewDecoder(r.Body).Decode(&newAdmins)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	addedAdmins, err := repository.AddAdmins(newAdmins)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string         `json:"status"`
		Count  int            `json:"count"`
		Data   []models.Admin `json:"data"`
	}{
		Status: "success",
		Count:  len(addedAdmins),
		Data:   addedAdmins,
	}
	json.NewEncoder(w).Encode(response)
}
