package handlers

import (
	"net/http"

	"github.com/clivethescott/todo/db"
	"github.com/gin-gonic/gin"
)

type handler struct {
	repo db.Repo
}

func (h *handler) hello() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "Hello")
	}
}

// NewMux creates a new Mux and registers routes
func NewMux(repo db.Repo) http.Handler {
	h := &handler{repo}

	mux := gin.Default()
	mux.GET("/", h.hello())

	return mux
}
