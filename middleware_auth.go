package main

import (
	"fmt"
	"net/http"

	"github.com/My-Golang-Projects/RSS-Scraper/internal/auth"

	"github.com/My-Golang-Projects/RSS-Scraper/internal/database"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header.Clone())
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Auth Error: %v", err))
			return
		}

		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Couldn't get user: %v", err))
			return
		}

		handler(w, r, user)
	}
}
