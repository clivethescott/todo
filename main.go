package main

import (
	"context"
	"log"

	"github.com/clivethescott/todo/server"
	"github.com/clivethescott/todo/store"
)

func main() {
	store, err := store.New(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	if err := server.Start(8080, store); err != nil {
		log.Fatal(err)
	}
}
