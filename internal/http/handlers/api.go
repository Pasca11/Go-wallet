package handlers

import (
	storage2 "github.com/Pasca11/storage"
	"github.com/gorilla/mux"
	"log"
	"net/http"
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
