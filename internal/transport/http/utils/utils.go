package utils

import (
	"encoding/json"
	"fmt"
	"github.com/Pasca11/types"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

// TODO take key from env
const secretKey = "secret"

func GetClaimsFromRequest(r *http.Request) (jwt.MapClaims, error) {
	token, err := ExtractRequestToken(r)
	if err != nil {
		return nil, err
	}
	claims, err := GetClaimsFromToken(token)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func GetClaimsFromToken(token *jwt.Token) (jwt.MapClaims, error) {
	claims := token.Claims.(jwt.MapClaims)
	if !claims.VerifyExpiresAt(time.Now().Unix(), false) {
		return nil, fmt.Errorf("acces denied (expired)")
	}
	return claims, nil
}

func ExtractRequestToken(r *http.Request) (*jwt.Token, error) {
	jwtKey := r.Header.Get("Authorization")
	if jwtKey == "" {
		return nil, fmt.Errorf("acces denied 3")
	}
	token, err := ValidateJWT(jwtKey)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, err
	}
	return token, nil
}

func ValidateJWT(stringToken string) (*jwt.Token, error) {
	return jwt.Parse(stringToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

}

func CreateJWT(acc *types.Account) (string, error) {
	claims := &jwt.MapClaims{
		"exp": time.Second * 300,
		"sub": acc.ID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secretKey))

}

func RenderJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(&v)
	if err != nil {
		http.Error(w, "Can`t render result", http.StatusInternalServerError)
	}
}

func GetIDFromRequest(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	return strconv.Atoi(id)
}
