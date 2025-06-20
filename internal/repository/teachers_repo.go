package repository

import (
	"database/sql"
	"fmt"
	"net/http"
	"restproject/internal/models"
	"restproject/internal/repository/sqlconnect"
	"restproject/pkg/utils"
)

func DeleteTeacherById(id int) error {
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		return utils.ErrorHandler(err, "error deleting teacher")
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM teachers WHERE id = ?", id)
	if err != nil {
		return utils.ErrorHandler(err, "error deleting teacher")
	}

	fmt.Println(result.RowsAffected())
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utils.ErrorHandler(err, "error deleting teacher")
	}
	if rowsAffected == 0 {
		return utils.ErrorHandler(err, "error deleting teacher")
	}
	return nil
}

func UpdateTeacher(id int, updatedTeacher models.Teacher) (models.Teacher, error) {
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "error updating teacher")
	}
	defer db.Close()

	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT * FROM teachers WHERE id = ?", id).Scan(&existingTeacher.ID, &existingTeacher.Firstname, &existingTeacher.Lastname, &existingTeacher.Subject, &existingTeacher.Class, &existingTeacher.Email)
	if err == sql.ErrNoRows {
		return models.Teacher{}, utils.ErrorHandler(err, "error updating teacher")
	} else if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "error updating teacher")
	}
	updatedTeacher.ID = existingTeacher.ID
	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ? ,subject = ? WHERE id = ?", &updatedTeacher.Firstname, &updatedTeacher.Lastname, &updatedTeacher.Email, &updatedTeacher.Class, &updatedTeacher.Subject, &updatedTeacher.ID)
	if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "error updating teacher")
	}
	return updatedTeacher, nil
}

func AddTeachers(newTeachers []models.Teacher) ([]models.Teacher, error) {
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "error adding teacher")
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES (?,?,?,?,?)")
	if err != nil {
		return nil, utils.ErrorHandler(err, "error adding teacher")
	}
	defer stmt.Close()

	addedTeachers := make([]models.Teacher, len(newTeachers))
	for i, newTeacher := range newTeachers {
		res, err := stmt.Exec(newTeacher.Firstname, newTeacher.Lastname, newTeacher.Email, newTeacher.Class, newTeacher.Subject)
		if err != nil {
			return nil, utils.ErrorHandler(err, "error adding teacher")
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			return nil, utils.ErrorHandler(err, "error adding teacher")
		}
		newTeacher.ID = int(lastID)
		addedTeachers[i] = newTeacher
	}
	return addedTeachers, nil
}

func GetTeacherByID(id int) (models.Teacher, error) {
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "error getting teacher")
	}
	defer db.Close()

	var teacher models.Teacher
	err = db.QueryRow("SELECT * FROM teachers WHERE id = ?", id).Scan(&teacher.ID, &teacher.Firstname, &teacher.Lastname, &teacher.Subject, &teacher.Class, &teacher.Email)
	if err == sql.ErrNoRows {
		return models.Teacher{}, utils.ErrorHandler(err, "error getting teacher")
	} else if err != nil {
		fmt.Println(err)
		return models.Teacher{}, utils.ErrorHandler(err, "error getting teacher")
	}
	return teacher, nil
}

func GetTeachers(teachers []models.Teacher, r *http.Request) ([]models.Teacher, error) {
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "error getting teachers")
	}
	defer db.Close()

	query := "SELECT * FROM teachers WHERE 1=1"
	var args []any
	query, args = addFilters(r, query, args)

	rows, err := db.Query(query, args...)
	if err != nil {
		fmt.Println(err)
		return nil, utils.ErrorHandler(err, "error getting teachers")
	}
	defer rows.Close()

	for rows.Next() {
		var teacher models.Teacher
		err := rows.Scan(&teacher.ID, &teacher.Firstname, &teacher.Lastname, &teacher.Class, &teacher.Subject, &teacher.Email)
		if err != nil {
			return nil, utils.ErrorHandler(err, "error getting teachers")
		}
		teachers = append(teachers, teacher)
	}
	return teachers, nil
}

func addFilters(r *http.Request, query string, args []any) (string, []any) {
	params := map[string]string{
		"first_name": "first_name",
		"last_name":  "last_name",
		"email":      "email",
		"class":      "class",
		"subject":    "subject",
	}

	for param, dbField := range params {
		value := r.URL.Query().Get(param)
		if value != "" {
			query += " AND " + dbField + " = ?"
			args = append(args, value)
		}
	}
	return query, args
}
