package router

import (
	"net/http"
	"restproject/internal/api/handlers"
)

func RegisterAdminRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /admins/", handlers.AddAdmins)
}
