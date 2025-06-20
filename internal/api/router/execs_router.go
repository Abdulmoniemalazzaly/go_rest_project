package router

import (
	"net/http"
	"restproject/internal/api/handlers"
)

func RegisterExecsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/execs/", handlers.ExecsHandler)
}
