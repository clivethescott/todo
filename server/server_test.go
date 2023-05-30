package server_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/clivethescott/todo/model"
	"github.com/clivethescott/todo/server"
	"golang.org/x/net/context"
)

type inMemoryStore struct {
	todos []*model.Todo
}

func (s *inMemoryStore) GetById(ctx context.Context, id string) (*model.Todo, error) {
	for _, t := range s.todos {
		if t.ID == id {
			return t, nil
		}
	}
	return nil, errors.New("get todo: no todo by id: " + id)
}

func (s *inMemoryStore) GetAll(ctx context.Context) ([]*model.Todo, error) {
	return s.todos, nil
}

func (s *inMemoryStore) Create(ctx context.Context, task string) (*model.Todo, error) {
	todo := model.NewTodo(task)
	s.todos = append(s.todos, todo)
	return todo, nil
}

func (s *inMemoryStore) MarkDone(ctx context.Context, id string) (*model.Todo, error) {
	for _, t := range s.todos {
		if t.ID == id {
			t.Completed = true
			return t, nil
		}
	}
	return nil, errors.New("mark done: no todo by id: " + id)
}

func TestServer(t *testing.T) {
	store := &inMemoryStore{todos: []*model.Todo{}}
	router := server.New(store, server.WithQuietStartup())

	w := httptest.NewRecorder()
	reqBody := `{"task": "fishing at the lake"}`
	req, _ := http.NewRequest(http.MethodPost, "/todo", strings.NewReader(reqBody))

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status=200, got status=%d", w.Code)
	}

	var todo model.Todo
	if err := json.Unmarshal(w.Body.Bytes(), &todo); err != nil {
		t.Errorf("unmarshall response as todo: %s %v\n", w.Body.String(), err)
	}

	if todo.Task != "fishing at the lake" {
		t.Errorf("expected task=fishing at the lake, got %q", todo.Task)
	}
}
