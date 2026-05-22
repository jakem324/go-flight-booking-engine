package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type TypedHandler[T any] func(c *gin.Context, body T)

func WithJSONBody[T any](handler TypedHandler[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req T

		if err := c.ShouldBindJSON(&req); err != nil {
			c.String(http.StatusBadRequest, "Invalid request payload")
			return
		}

		handler(c, req)
	}
}
