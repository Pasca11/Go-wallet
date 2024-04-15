package handlers

import (
	"encoding/json"
	"github.com/Pasca11/internal/http/utils"
	storage2 "github.com/Pasca11/storage"
	"github.com/Pasca11/types"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
)

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
	utils.RenderJSON(w, 200, acc)
}
func (a *App) handleGetAccount(w http.ResponseWriter, r *http.Request) {
	accs, err := a.storage.GetAllAccounts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.RenderJSON(w, 200, accs)
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

	tokenStr, err := utils.CreateJWT(acc)
	if err != nil {
		http.Error(w, "Cant create token", http.StatusInternalServerError)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(req.Password)); err != nil {
		http.Error(w, "Incorrect password/login", http.StatusBadRequest)
		return
	}

	utils.RenderJSON(w, 200, types.LoginResponse{
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

	utils.RenderJSON(w, 200, acc)
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
	utils.RenderJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func (a *App) handelTransfer(w http.ResponseWriter, r *http.Request) {
	var req types.TransferRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error()+"1", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	claims, err := utils.GetClaimsFromRequest(r)
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

	utils.RenderJSON(w, http.StatusOK, "transfer succeeded")
}

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
