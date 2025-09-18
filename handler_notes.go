package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bootdotdev/learn-cicd-starter/internal/database"
	"github.com/google/uuid"
)

// Helper-Funktionen für HTTP-Antworten
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if payload != nil {
		_ = json.NewEncoder(w).Encode(payload)
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// GET /notes
func (cfg *apiConfig) handlerNotesGet(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := cfg.DB.GetNotesForUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get posts for user: "+err.Error())
		return
	}

	postsResp, err := databasePostsToPosts(posts)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't convert posts: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, postsResp)
}

// POST /notes
func (cfg *apiConfig) handlerNotesCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Note string `json:"note"`
	}

	decoder := json.NewDecoder(r.Body)
	var params parameters
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters: "+err.Error())
		return
	}

	id := uuid.NewString()
	noteData := database.CreateNoteParams{
		ID:        id,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
		Note:      params.Note,
		UserID:    user.ID,
	}

	if err := cfg.DB.CreateNote(r.Context(), noteData); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create note: "+err.Error())
		return
	}

	note, err := cfg.DB.GetNote(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get note: "+err.Error())
		return
	}

	noteResp, err := databaseNoteToNote(note)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't convert note: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, noteResp)
}
