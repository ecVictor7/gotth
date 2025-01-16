package main

import (
	"log"
	"os"

	"github.com/ecvictor7/gotth/internal/server"
	"github.com/ecvictor7/gotth/internal/store"
)

func main() {
	logger := log.New(os.Stoudt, "[Spooktober] ", log.LstdFlags)

	port := 9000

	logger.Print("Creating guests stroe..")
	guestDb := store.NewGuestStore(logger)
	guestDb.AddGuest(store.Guest{Name: "Sigrid", Email: "sig@fake-email.no"})

	srv, err := server.NewServer(logger, port, guestDb)
	if err != nil {
		logger.Fatalf("Error when creating server: %s", err)
		os.Exit(1)
	}
	if err := srv.Start(); err != nil {
		logger.Fatalf("Error when starting server: %s", err)
		os.Exit(1)
	}
}
