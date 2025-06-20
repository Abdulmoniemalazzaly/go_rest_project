package repository

import (
	"database/sql"
	"fmt"
	"net/http"
	"restproject/internal/models"
	"restproject/internal/repository/sqlconnect"
)

func DeleteTeacherById(id int) error {
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		// log.Println(w, "Unable to connect to database", http.StatusInternalServerError)
		return err
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM teachers WHERE id = ?", id)
	if err != nil {
		// http.Error(w, "Error deleting teacher", http.StatusInternalServerError)
		return err
	}

	fmt.Println(result.RowsAffected())
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// http.Error(w, "Error retrieving delete result", http.StatusInternalServerError)
		return err
	}
	if rowsAffected == 0 {
		// http.Error(w, "Teacher not found", http.StatusInternalServerError)
		return err
	}
	return nil
}

func UpdateTeacher(id int, updatedTeacher models.Teacher) (models.Teacher, error) {
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		// log.Println(w, "Unable to connect to database", http.StatusInternalServerError)
		return models.Teacher{}, err
	}
	defer db.Close()

	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT * FROM teachers WHERE id = ?", id).Scan(&existingTeacher.ID, &existingTeacher.Firstname, &existingTeacher.Lastname, &existingTeacher.Subject, &existingTeacher.Class, &existingTeacher.Email)
	if err == sql.ErrNoRows {
		// http.Error(w, "Teacher not found", http.StatusNotFound)
		return models.Teacher{}, err
	} else if err != nil {
		// http.Error(w, "Unable to retrieve data", http.StatusInternalServerError)
		return models.Teacher{}, err
	}
	updatedTeacher.ID = existingTeacher.ID
	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ? ,subject = ? WHERE id = ?", &updatedTeacher.Firstname, &updatedTeacher.Lastname, &updatedTeacher.Email, &updatedTeacher.Class, &updatedTeacher.Subject, &updatedTeacher.ID)
	if err != nil {
		// http.Error(w, "Error updating teacher", http.StatusInternalServerError)
		return models.Teacher{}, err
	}
	return updatedTeacher, nil
}

func AddTeachers(newTeachers []models.Teacher) ([]models.Teacher, error) {
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		// http.Error(w, "Error connecting to database", http.StatusInternalServerError)
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES (?,?,?,?,?)")
	if err != nil {
		// http.Error(w, "Error in preparing SQL query", http.StatusInternalServerError)
		return nil, err
	}
	defer stmt.Close()

	addedTeachers := make([]models.Teacher, len(newTeachers))
	for i, newTeacher := range newTeachers {
		res, err := stmt.Exec(newTeacher.Firstname, newTeacher.Lastname, newTeacher.Email, newTeacher.Class, newTeacher.Subject)
		if err != nil {
			// http.Error(w, "Error inserting data into database", http.StatusInternalServerError)
			return nil, err
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			// http.Error(w, "Error getting last insert ID", http.StatusInternalServerError)
			return nil, err
		}
		newTeacher.ID = int(lastID)
		addedTeachers[i] = newTeacher
	}
	return addedTeachers, nil
}

func GetTeacherByID(id int) (models.Teacher, error) {
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		// http.Error(w, "Error connecting to database", http.StatusInternalServerError)
		return models.Teacher{}, err
	}
	defer db.Close()

	var teacher models.Teacher
	err = db.QueryRow("SELECT * FROM teachers WHERE id = ?", id).Scan(&teacher.ID, &teacher.Firstname, &teacher.Lastname, &teacher.Subject, &teacher.Class, &teacher.Email)
	if err == sql.ErrNoRows {
		// http.Error(w, "Teacher not found", http.StatusNotFound)
		return models.Teacher{}, err
	} else if err != nil {
		fmt.Println(err)
		// http.Error(w, "Database Query Error", http.StatusInternalServerError)
		return models.Teacher{}, err
	}
	return teacher, nil
}

func GetTeachers(teachers []models.Teacher, r *http.Request) ([]models.Teacher, error) {
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		// http.Error(w, "Error connecting to database", http.StatusInternalServerError)
		return nil, err
	}
	defer db.Close()

	query := "SELECT * FROM teachers WHERE 1=1"
	var args []any
	query, args = addFilters(r, query, args)

	rows, err := db.Query(query, args...)
	if err != nil {
		fmt.Println(err)
		// http.Error(w, "Database Query Error", http.StatusInternalServerError)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var teacher models.Teacher
		err := rows.Scan(&teacher.ID, &teacher.Firstname, &teacher.Lastname, &teacher.Class, &teacher.Subject, &teacher.Email)
		if err != nil {
			// http.Error(w, "Error scanning database results", http.StatusInternalServerError)
			return nil, err
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
