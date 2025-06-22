package repository

import (
	"restproject/internal/models"
	"restproject/internal/repository/sqlconnect"
	"restproject/pkg/utils"
)

func AddAdmins(newAdmins []models.Admin) ([]models.Admin, error) {
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "error adding admin")
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO admins (username, password) VALUES (?,?)")
	if err != nil {
		return nil, utils.ErrorHandler(err, "error adding admin")
	}
	defer stmt.Close()

	addedAdmins := make([]models.Admin, len(newAdmins))
	for i, newAdmin := range newAdmins {
		res, err := stmt.Exec(newAdmin.Username, newAdmin.Password)
		if err != nil {
			return nil, utils.ErrorHandler(err, "error adding teacher")
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			return nil, utils.ErrorHandler(err, "error adding teacher")
		}
		newAdmin.ID = int(lastID)
		addedAdmins[i] = newAdmin
	}
	return addedAdmins, nil
}
