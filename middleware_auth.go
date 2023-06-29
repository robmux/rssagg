package main

import (
	"fmt"
	"net/http"

	"github.com/robmux/rssagg/internal/auth"
	"github.com/robmux/rssagg/internal/database"
)

type authedHandler func(w http.ResponseWriter, r *http.Request, user database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		apiKey, err := auth.GetAPIKey(request.Header)
		if err != nil {
			respondWithError(writer, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}

		user, err := apiCfg.DB.GetUserByAPIKEY(request.Context(), apiKey)

		if err != nil {
			respondWithError(writer, 400, fmt.Sprintf("error getting user: %v", err))
			return
		}
		handler(writer, request, user)
	}
}
