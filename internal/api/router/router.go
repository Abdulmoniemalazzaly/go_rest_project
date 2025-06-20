package router

import (
	"net/http"
	"restproject/internal/api/handlers"
)

func Router() *http.ServeMux {
	mux := http.NewServeMux()
	// ROOT
	mux.HandleFunc("/", handlers.RootHandler)

	RegisterTeachersRoutes(mux)
	RegisterStudentsRoutes(mux)
	RegisterExecsRoutes(mux)
	return mux
}
