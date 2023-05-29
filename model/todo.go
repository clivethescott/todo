package model

import (
	"errors"

	"github.com/google/uuid"
)

type Todo struct {
	ID        string `json:"id"`
	Task      string `json:"task"`
	Completed bool   `json:"completed"`
}

var ErrTodoNotFound = errors.New("todo not found")

func NewTodo(task string) *Todo {
	return &Todo{
		ID:        uuid.New().String(),
		Task:      task,
		Completed: false,
	}
}
