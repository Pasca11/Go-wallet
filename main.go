package main

import (
	"github.com/Pasca11/internal/http/handlers"
	"github.com/Pasca11/storage"
	"log"
)

func main() {
	store, err := storage.NewPostgresStorage()
	if err != nil {
		log.Fatalln(err)
	}
	err = store.Init()
	if err != nil {
		log.Fatalln(err)
	}
	app := handlers.NewApp(":3000", store)
	app.Start()
}
