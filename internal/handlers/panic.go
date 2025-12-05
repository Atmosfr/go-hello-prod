package handlers

import "net/http"

func PanicHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/panic" {
		http.NotFound(w, r)
		return
	}
	panic("intentional panic for testing recovery middleware")
}
