package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ecVictor7/gotth/store"
	"github.com/ecVictor7/gotth/templates"
)

type GuestStore interface {
	AddGuest(guest store.Guest) error
	GetGuests() ([]store.Guest, error)
}

type server struct {
	logger     *log.Logger
	port       int
	httpServer *http.Server
	guestDb    GuestStore
}

// Create a new server instance with the given logger and port
func NewServer(
	logger *log.Logger, port int, guestDb GuestStore) (*server, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}
	if guestDb == nil {
		return nil, fmt.Errorf("guestDb is required")
	}
	return &server{
		logger:  logger,
		port:    port,
		guestDb: guestDb}, nil

}

func (s *server) Start() error {
	s.logger.Printf("Starting server on port %d", s.port)
	var stopChan chan os.Signal

	//define router
	router := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static"))
	router.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	//routes
	router.HandleFunc("GET /", s.defaultHandler)
	router.HandleFunc("GET /guests", s.getGuestsHandler)
	router.HandleFunc("POST /guests", s.addGuestHandler)

	//define server
	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: router}

	//creeate channel to listen for signals
	stopChan = make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Fatalf("Error when running server: %s", err)
		}
	}()
	<-stopChan

	if err := s.httpServer.Shutdown(context.Background()); err != nil {
		log.Fatalf("Error when shutting down server: %v", err)
		return err
	}
	return nil
}

// GET /
func (s *server) defaultHandler(w http.ResponseWriter, r *http.Request) {
	//w.WriteHeader(http.StatusOK)
	//w.Write([]byte("My spooky halloween party!"))

	homeTpl := templates.Home()
	templates.Layout(homeTpl, "Home").Render(r.Context(), w)
}

func (s *server) getGuestsHandler(w http.ResponseWriter, r *http.Request) {
	guests, err := s.guestDb.GetGuests()
	if err != nil {
		s.logger.Printf("Error when fetching guests: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	templates.Guests(guests).Render(r.Context(), w)

}

func (s *server) addGuestHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.Printf("Error when reading request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	guest, err := store.DecodeGuest(payload)
	if err != nil {
		s.logger.Printf("Error when decoding guest: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := s.guestDb.AddGuest(guest); err != nil {
		s.logger.Printf("Error when adding guest: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	guests, err := s.guestDb.GetGuests()
	if err != nil {
		s.logger.Printf("Error when getting guests: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	templates.Guests(guests).Render(r.Context(), w)
}
