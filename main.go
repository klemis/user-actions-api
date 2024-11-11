package main

import (
	"flag"
	"log"

	"github.com/klemis/user-actions-api/api"
	"github.com/klemis/user-actions-api/storage"
)

func main() {
	listenAddr := flag.String("listenaddr", ":8080", "api server address")
	flag.Parse()

	store, err := storage.NewInMemoryStorage("users.json", "actions.json")
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	server := api.NewServer(*listenAddr, store)
	log.Println("API server running on port: ", *listenAddr)
	log.Fatal(server.Start())
}
