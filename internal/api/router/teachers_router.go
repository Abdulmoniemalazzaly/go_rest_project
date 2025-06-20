package router

import (
	"net/http"
	"restproject/internal/api/handlers"
)

func RegisterTeachersRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /teachers/", handlers.GetTeachersHandlers)
	mux.HandleFunc("GET /teachers/{id}", handlers.GetTeacherHandlers)
	mux.HandleFunc("POST /teachers/", handlers.AddTeachersHandler)
	mux.HandleFunc("DELETE /teachers/{id}", handlers.DeleteTeacherHandler)
	mux.HandleFunc("PUT /teachers/{id}", handlers.UpdateTeacherHandler)
}
