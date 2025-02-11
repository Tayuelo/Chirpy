package main

import (
	"chirpy/internal/auth"
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
	}

	payload := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	user, err := cfg.queries.GetUserByEmail(r.Context(), payload.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	if err = auth.CheckPasswordHash(payload.Password, user.HashedPassword); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:    user.ID,
			Email: user.Email,
		},
	})
}
