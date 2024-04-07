package handlers

import (
	"encoding/json"
	"fmt"
	storage2 "github.com/Pasca11/storage"
	"github.com/Pasca11/types"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type App struct {
	ListenAddr string
	storage    storage2.Storage
}

func NewApp(address string, storage storage2.Storage) *App {
	return &App{
		ListenAddr: address,
		storage:    storage,
	}
}

func (a *App) Start() {
	router := mux.NewRouter()

	router.HandleFunc("/login", a.handleLogin).Methods("POST")
	//router.Use(JWTMiddleware)
	router.HandleFunc("/account", a.handleGetAccount).Methods("GET")
	router.HandleFunc("/account/{id}", JWTMiddleware(a.handleGetAccountByID, a.storage)).Methods("GET")
	router.HandleFunc("/account", a.handleCreateAccount).Methods("POST")
	router.HandleFunc("/account/{id}", a.handleDeleteAccount).Methods("DELETE")

	router.HandleFunc("/transfer", a.handelTransfer).Methods("POST")

	log.Println("Server is started")
	log.Fatalln(http.ListenAndServe(a.ListenAddr, router))
}

func (a *App) handleGetAccountByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	acc, err := a.storage.GetAccountByID(id)
	if err != nil {
		http.Error(w, "Cant find user. Try again", http.StatusInternalServerError)
		return
	}
	RenderJSON(w, 200, acc)
}
func (a *App) handleGetAccount(w http.ResponseWriter, r *http.Request) {
	accs, err := a.storage.GetAllAccounts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	RenderJSON(w, 200, accs)
}

func (a *App) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req types.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid login/password", http.StatusBadRequest)
		return
	}

	acc, err := a.storage.GetAccountByWallet(req.Wallet)
	if err != nil {
		http.Error(w, "Cant get account", http.StatusBadRequest)
		return
	}

	tokenStr, err := createJWT(acc)
	if err != nil {
		http.Error(w, "Cant create token", http.StatusInternalServerError)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(req.Password)); err != nil {
		http.Error(w, "Incorrect password/login", http.StatusBadRequest)
		return
	}

	RenderJSON(w, 200, types.LoginResponse{
		Wallet: acc.Wallet,
		Token:  tokenStr,
	})
}

func (a *App) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	var req types.CreateAccountRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	accProto, err := types.NewAccount(req.FirstName, req.LastName, req.Patronymic, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	acc, err := a.storage.CreateAccount(accProto)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Some error", http.StatusInternalServerError)
		return
	}
	//
	//tokenString, err := createJWT(acc)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}

	//fmt.Println(tokenString)
	//w.Header().Add("Authorization", tokenString)

	RenderJSON(w, 200, acc)
}

func (a *App) handleDeleteAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	err = a.storage.DeleteAccount(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	RenderJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func (a *App) handelTransfer(w http.ResponseWriter, r *http.Request) {
	var req types.TransferRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error()+"1", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	claims, err := getClaimsFromRequest(r)
	if err != nil {
		http.Error(w, err.Error()+"2", http.StatusBadRequest)
		return
	}

	claimsID, ok := claims["sub"]
	acc, _ := a.storage.GetAccountByID(int(claimsID.(float64)))

	if !ok || int(acc.Wallet) != req.From {
		http.Error(w, "access denied", http.StatusBadRequest)
		return
	}

	err = a.storage.Transfer(req)
	if err != nil {
		http.Error(w, "cant process transfer", http.StatusBadRequest)
		return
	}

	RenderJSON(w, http.StatusOK, "transfer succeeded")
}

// TODO take key from env
const secretKey = "secret"

func JWTMiddleware(h http.HandlerFunc, s storage2.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := getClaimsFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := getIDFromRequest(r)
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

func getClaimsFromRequest(r *http.Request) (jwt.MapClaims, error) {
	token, err := extractRequestToken(r)
	if err != nil {
		return nil, err
	}
	claims, err := getClaimsFromToken(token)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func getClaimsFromToken(token *jwt.Token) (jwt.MapClaims, error) {
	claims := token.Claims.(jwt.MapClaims)
	if !claims.VerifyExpiresAt(time.Now().Unix(), false) {
		return nil, fmt.Errorf("acces denied (expired)")
	}
	return claims, nil
}

func extractRequestToken(r *http.Request) (*jwt.Token, error) {
	jwtKey := r.Header.Get("Authorization")
	if jwtKey == "" {
		return nil, fmt.Errorf("acces denied 3")
	}
	token, err := validateJWT(jwtKey)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, err
	}
	return token, nil
}

func validateJWT(stringToken string) (*jwt.Token, error) {
	return jwt.Parse(stringToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

}

func createJWT(acc *types.Account) (string, error) {
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

func getIDFromRequest(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	return strconv.Atoi(id)
}
