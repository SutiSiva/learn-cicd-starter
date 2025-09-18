package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bootdotdev/learn-cicd-starter/internal/database"
	"github.com/google/uuid"
)

// handlerNotesGet gibt alle Notizen eines Users zurück
func (cfg *apiConfig) handlerNotesGet(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := cfg.DB.GetNotesForUser(r.Context(), user.ID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't get posts for user", err)
		return
	}

	postsResp, err := databasePostsToPosts(posts)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't convert posts", err)
		return
	}

	RespondWithJSON(w, http.StatusOK, postsResp)
}

// handlerNotesCreate erstellt eine neue Notiz für einen User
func (cfg *apiConfig) handlerNotesCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Note string `json:"note"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	id := uuid.New().String()
	err := cfg.DB.CreateNote(r.Context(), database.CreateNoteParams{
		ID:        id,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
		Note:      params.Note,
		UserID:    user.ID,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't create note", err)
		return
	}

	note, err := cfg.DB.GetNote(r.Context(), id)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Couldn't get note", err)
		return
	}

	noteResp, err := databaseNoteToNote(note)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't convert note", err)
		return
	}

	RespondWithJSON(w, http.StatusCreated, noteResp)
}
