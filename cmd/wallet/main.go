package main

import (
	"github.com/Pasca11/internal/transport/http/handlers"
	"log"
)

func main() {
	store, err := initStorage()
	if err != nil {
		log.Fatalln(err)
		return
	}
	app := handlers.NewApp(":3000", store)
	app.Start()
}
