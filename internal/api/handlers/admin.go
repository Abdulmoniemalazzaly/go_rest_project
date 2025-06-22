package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"restproject/internal/models"
	"restproject/internal/repository"
	"restproject/pkg/utils"
)

func AddAdmins(w http.ResponseWriter, r *http.Request) {
	var newAdmins []models.Admin
	err := json.NewDecoder(r.Body).Decode(&newAdmins)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = validateAndEncryptPassword(newAdmins)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

func validateAndEncryptPassword(newAdmins []models.Admin) error {
	for i, newAdmin := range newAdmins {
		if newAdmin.Password == "" {
			return utils.ErrorHandler(errors.New("password is blank"), "please enter password")
		}
		if newAdmin.Username == "" {
			return utils.ErrorHandler(errors.New("username is blank"), "please enter username")
		}

		encryptedPassword, err := utils.EncryptPassword(newAdmin.Password)
		if err != nil {
			return utils.ErrorHandler(errors.New("error handling password"), "error handling password")
		}
		newAdmins[i].Password = encryptedPassword
	}
	return nil
}
