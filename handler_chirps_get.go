package main

import (
	"chirpy/internal/database"
	"net/http"
	"sort"
	"time"

	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	Body      string    `json:"body"`
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	authorID := r.URL.Query().Get("author_id")
	sortDirection := r.URL.Query().Get("sort")
	var chirps []database.Chirp
	var err error

	if authorID == "" {
		chirps, err = cfg.queries.GetChirps(r.Context())
	} else {
		authorUUID, parseErr := uuid.Parse(authorID)
		if parseErr != nil {
			respondWithError(w, http.StatusBadRequest, parseErr.Error(), parseErr)
			return
		}

		chirps, err = cfg.queries.GetChirpsByUserId(r.Context(), authorUUID)
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	respondWithJSON(w, http.StatusOK, MapChirps(chirps, sortDirection))
}

func MapChirps(chirps []database.Chirp, sortDirection string) []Chirp {
	slc := make([]Chirp, len(chirps))

	for i, chirp := range chirps {
		slc[i] = Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
	}

	sort.Slice(slc, func(i, j int) bool {
		if sortDirection == "desc" {
			return slc[i].CreatedAt.After(slc[j].CreatedAt)
		}
		return slc[i].CreatedAt.Before(slc[j].CreatedAt)
	})

	return slc
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}
	chirpDb, err := cfg.queries.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error(), err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirpID,
		CreatedAt: chirpDb.CreatedAt,
		UpdatedAt: chirpDb.UpdatedAt,
		Body:      chirpDb.Body,
		UserID:    chirpDb.UserID,
	})
}
