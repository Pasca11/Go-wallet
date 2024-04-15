package middlewares

import (
	utils "github.com/Pasca11/internal/http/utils"
	storage2 "github.com/Pasca11/storage"
	"net/http"
)

// TODO take key from env
const secretKey = "secret"

func JWTMiddleware(h http.HandlerFunc, s storage2.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := utils.GetClaimsFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := utils.GetIDFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		claimID := claims["sub"]
		acc, err := s.GetAccountByID(id)
		if acc.ID != int((claimID.(float64))) {
			http.Error(w, "access denied", http.StatusForbidden)
			return
		}
		h.ServeHTTP(w, r)
	}
}
