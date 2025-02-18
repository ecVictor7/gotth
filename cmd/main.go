package main

import (
	"log"
	"os"

	"github.com/ecVictor7/gotth/server"
	"github.com/ecVictor7/gotth/store"
)

func main() {
	logger := log.New(os.Stdout, "[Spooktober] ", log.LstdFlags)

	port := 9000

	logger.Print("Creating guests store ..")
	guestDb := store.NewGuestStore(logger)
	guestDb.AddGuest(store.Guest{Name: "Sigrid", Email: "sig@fake-email.no"})

	srv, err := server.NewServer(logger, port, guestDb)
	if err != nil {
		logger.Fatalf("Error when  creating server: %s", err)
		os.Exit(1)
	}
	if err := srv.Start(); err != nil {
		logger.Fatalf("Error when starting server: %s", err)
		os.Exit(1)
	}
}
