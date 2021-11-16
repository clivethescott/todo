package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/clivethescott/todo/models"
)

// Config represents config for the store
type Config struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifeTime time.Duration
}

// Repo provides persistence for Todos
type Repo interface {
	// AddTodo adds a new Todo in the store
	AddTodo(models.Todo) error
	Close(context.Context) error
}

// MySQLRepo is a MySQL DB
type MySQLRepo struct {
	db *sql.DB
}

// NewRepo connects to MySQL and returns a repo
func NewRepo(cfg Config) (*MySQLRepo, error) {
	db, err := sql.Open("mysql", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("invalid DB url: %w", err)
	}

	log.Println("MySQL connecting...")
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to the DB: %w", err)
	}
	log.Println("MySQL...connected")

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifeTime)
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	return &MySQLRepo{db: db}, nil
}

// AddTodo adds a new Todo in the DB
func (r *MySQLRepo) AddTodo(models.Todo) error {
	return nil
}

// Close closes the repo
func (r *MySQLRepo) Close(ctx context.Context) error {

	ch := make(chan error)

	go func() {
		ch <- r.db.Close()
	}()

	select {
	case <-ctx.Done():
		return errors.New("deadline exceeded")
	case err := <-ch:
		log.Println("DB connections closed")
		return err
	}
}
