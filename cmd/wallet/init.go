package main

import (
	"github.com/Pasca11/pkg/storage"
	"log"
)

func initStorage() (storage.Storage, error) {
	store, err := storage.NewPostgresStorage()
	if err != nil {
		log.Fatalln(err)
	}
	err = store.Init()
	if err != nil {
		log.Fatalln(err)
	}
	return store, nil
}
