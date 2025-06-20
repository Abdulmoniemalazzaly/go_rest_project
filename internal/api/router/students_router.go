package router

import (
	"net/http"
	"restproject/internal/api/handlers"
)

func RegisterStudentsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/students/", handlers.StudentsHandler)
}
