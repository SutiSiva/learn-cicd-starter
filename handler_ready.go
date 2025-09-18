package main

import "net/http"

// handlerReadiness prüft, ob der Service bereit ist
func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
