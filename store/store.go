package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/clivethescott/todo/model"
	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sql.DB
}

const timeout = 1 * time.Second

func (s *Store) GetById(ctx context.Context, id string) (*model.Todo, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	row := s.db.QueryRowContext(ctx, "SELECT id, task, completed FROM todo WHERE id = ?", id)
	var todo model.Todo

	if err := row.Scan(&todo.ID, &todo.Task, &todo.Completed); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrTodoNotFound
		}
		return nil, fmt.Errorf("get todo by id=%q: %w", id, err)
	}

	return &todo, nil
}

func (s *Store) GetAll(ctx context.Context) ([]*model.Todo, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	rows, err := s.db.QueryContext(ctx, "SELECT id, task, completed FROM todo")
	if err != nil {
		return nil, fmt.Errorf("get all todos: %w", err)
	}
	defer rows.Close()

	var todos []*model.Todo

	for rows.Next() {
		var todo model.Todo
		if err := rows.Scan(&todo.ID, &todo.Task, &todo.Completed); err != nil {
			return nil, fmt.Errorf("get all todos: %w", err)
		}

		todos = append(todos, &todo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get all todos: %w", err)
	}

	return todos, nil
}

func (s *Store) Create(ctx context.Context, task string) (*model.Todo, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	todo := model.NewTodo(task)
	_, err := s.db.ExecContext(ctx,
		"INSERT INTO todo(id, task, completed) VALUES(?, ?, ?)", todo.ID, todo.Task, todo.Completed)
	if err != nil {
		return nil, fmt.Errorf("create todo=%q: %w", task, err)
	}
	return todo, nil
}

func (s *Store) MarkDone(ctx context.Context, id string) (*model.Todo, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	result, err := s.db.ExecContext(ctx, "UPDATE todo SET completed = 1 WHERE id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("update todo: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("update todo count: %w", err)
	}

	if rows == 0 {
		return nil, model.ErrTodoNotFound
	}

	return s.GetById(ctx, id)
}

func (s *Store) Close() error {
	return s.db.Close()
}

func New(ctx context.Context) (*Store, error) {
	db, err := prepareDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("prepare DB: %w", err)
	}
	store := &Store{db}
	return store, nil
}

func prepareDB(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}

	_, err = db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS todo (
      id        TEXT NOT NULL PRIMARY KEY, 
      task      TEXT NOT NULL, 
      completed INTEGER NOT NULL DEFAULT 0 
    )`)
	if err != nil {
		return nil, fmt.Errorf("prepare table: %w", err)
	}
	return db, nil
}
