package main

import (
	"github.com/Pasca11/pkg/storage"
	"log"
)

func initStorage() {
	store, err := storage.NewPostgresStorage()
	if err != nil {
		log.Fatalln(err)
	}
}
