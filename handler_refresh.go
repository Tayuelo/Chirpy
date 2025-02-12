package main

import (
	"chirpy/internal/auth"
	"errors"
	"net/http"
	"time"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	refreshToken, err := cfg.queries.GetUserFromRefreshToken(r.Context(), bearer)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		errMessage := "token expired or revoked"
		respondWithError(w, http.StatusUnauthorized, errMessage, errors.New(errMessage))
		return
	}

	if refreshToken.RevokedAt.Valid {
		errMessage := "token expired or revoked"
		respondWithError(w, http.StatusUnauthorized, errMessage, errors.New(errMessage))
		return
	}

	token, err := auth.MakeJWT(refreshToken.UserID, cfg.jwtSecret, accessTokenExpirationTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	type response struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: token,
	})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	_, err = cfg.queries.RevokeRefreshToken(r.Context(), bearer)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
