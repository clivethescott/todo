package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/clivethescott/todo/db"
	"github.com/clivethescott/todo/handlers"
)

// Config represents config for the server
type Config struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	DB           db.Config
}

func shutdownListener(srv *http.Server, repo db.Repo) func() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	return func() {
		log.Println("starting server...")
		<-quit

		log.Println("closing DB connections")
		if err := repo.Close(context.TODO()); err != nil {
			log.Printf("error closing DB connections %v\n", err)
		}

		log.Println("shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("error shutting down server: %v\n", err)
		}
	}
}

// Start starts the server
func Start(cfg Config) error {
	repo, err := db.NewRepo(cfg.DB)
	if err != nil {
		return err
	}
	srv := http.Server{
		Handler: handlers.NewMux(repo),

		Addr:         cfg.Addr,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	listener := shutdownListener(&srv, repo)
	go listener()

	return srv.ListenAndServe()
}
