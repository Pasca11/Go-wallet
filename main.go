package main

import (
	"github.com/Pasca11/internal/http/handlers"
	storage2 "github.com/Pasca11/storage"
	"log"
)

func main() {
	storage, err := storage2.NewPostgresStorage()
	if err != nil {
		log.Fatalln(err)
	}
	err = storage.init()
	if err != nil {
		log.Fatalln(err)
	}
	app := handlers.NewApp(":3000", storage)
	app.Start()
}
