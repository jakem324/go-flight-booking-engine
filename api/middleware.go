package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TypedHandler[T any] func(c *gin.Context, body T)

func withJSONBody[T any](handler TypedHandler[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req T

		if err := c.ShouldBindJSON(&req); err != nil {
			c.String(http.StatusBadRequest, "Invalid request payload")
			return
		}

		handler(c, req)
	}
}

func genericErrorResponse(ctx *gin.Context, err error) {
	log.Fatalf("Handler error: %v\n", err)
	ctx.Status(http.StatusInternalServerError)
}
