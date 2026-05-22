// Package API defines app entry point
package api

import (
	//"context"
	"net/http"
	"github.com/gin-gonic/gin"
)

func main() {
	// ctx := context.Background()
	// handlers := setup(ctx)

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
    // Return JSON response
    c.JSON(http.StatusOK, gin.H{
      "message": "pong",
    })
  })
}
