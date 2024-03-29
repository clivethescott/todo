package server

import (
	"errors"
	"log"
	"net/http"

	"github.com/clivethescott/todo/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

type (
	routes struct {
		store TodoStore
	}
	createTodo struct {
		Task string `json:"task"`
	}
	TodoStore interface {
		GetById(ctx context.Context, id string) (*model.Todo, error)
		GetAll(ctx context.Context) ([]*model.Todo, error)
		Create(ctx context.Context, task string) (*model.Todo, error)
		MarkDone(ctx context.Context, id string) (*model.Todo, error)
	}
)

type ServerOptsFunc func(*gin.Engine)

func WithQuietStartup() ServerOptsFunc {
	return func(_ *gin.Engine) {
		gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {}
	}
}

func New(db TodoStore, opts ...ServerOptsFunc) *gin.Engine {
	r := gin.Default()
	for _, opt := range opts {
		opt(r)
	}
	if err := r.SetTrustedProxies(nil); err != nil {
		log.Printf("set trusted proxies: %v\n", err)
	}
	routes := &routes{db}

	r.GET("/todos", routes.GetAllTodos())
	r.GET("/todo/:id", routes.GetTodoById())
	r.PATCH("/todo/:id", routes.MarkDone())
	r.POST("/todo", routes.CreateTodo())
	return r
}

func (r *routes) GetAllTodos() gin.HandlerFunc {
	return func(c *gin.Context) {
		todos, err := r.store.GetAll(c.Request.Context())
		if err != nil {
			log.Printf("get todos: %v\n", err)
			internalServerError(c)
			return
		}

		if len(todos) == 0 {
			c.JSON(http.StatusOK, []model.Todo{})
			return
		}
		c.JSON(http.StatusOK, todos)
	}
}

func (r *routes) GetTodoById() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		todo, err := r.store.GetById(c.Request.Context(), id)
		if err != nil {
			if errors.Is(err, model.ErrTodoNotFound) {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
			log.Printf("get todo[id=%s]: %v\n", id, err)
			internalServerError(c)
			return
		}

		c.JSON(http.StatusOK, todo)
	}
}

func (r *routes) CreateTodo() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := new(createTodo)
		if err := c.BindJSON(req); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		task, err := r.store.Create(c.Request.Context(), req.Task)
		if err != nil {
			log.Printf("create todo: %v\n", err)
			internalServerError(c)
			return
		}

		c.JSON(http.StatusCreated, task)
	}
}

func (r *routes) MarkDone() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		todo, err := r.store.MarkDone(c.Request.Context(), id)
		if err != nil {
			if errors.Is(err, model.ErrTodoNotFound) {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
			log.Printf("get todo[id=%s]: %v\n", id, err)
			internalServerError(c)
			return
		}

		c.JSON(http.StatusOK, todo)
	}
}

func internalServerError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": "Unable to process request",
	})
}
