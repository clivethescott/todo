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
	srv := server.New(store)
	log.Fatal(srv.Run(":8080"))
}
