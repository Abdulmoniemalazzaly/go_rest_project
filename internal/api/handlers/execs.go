package handlers

import (
	"fmt"
	"net/http"
)

func ExecsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Handling execs")
}
