package handlers

import (
	"fmt"
	"net/http"
)

func TeachersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Handling orders")
}
