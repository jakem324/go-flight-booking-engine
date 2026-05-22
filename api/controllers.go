// Package api defines app entry point
package api

import (
	"context"
	"errors"
	"net/http"

	"booking.engine/domain/commands"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	ctx := context.Background()
	handlers := setup(ctx)

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
    // Return JSON response
    c.JSON(http.StatusOK, gin.H{
      "message": "pong",
    })
  })

	type CreatePencilBookingRequest struct {
		RequiredNumberOfSeats int `json:"requiredNumberOfSeats" binding:"required"`
		OutboundJourneyLegs []uuid.UUID `json:"outboundJourneyLegs" binding:"required"`
	}

	router.POST("/user", func(c *gin.Context) {
		var json CreatePencilBookingRequest

		// Bind incoming JSON to the struct
		// ShouldBindJSON returns an error if the JSON is invalid or missing required fields
		if jsonErr := c.ShouldBindJSON(&json); jsonErr != nil {
				c.Status(http.StatusBadRequest)
				return
		}

		createdBookingID, err := handlers.PencilBookingHandler.CreatePencilBooking(c, commands.CreatePencilBookingDto{
			RequiredNumberOfSeats: json.RequiredNumberOfSeats,
			OutboundJourneyLegs: json.OutboundJourneyLegs,
		})

		var seatsUnavailableError *commands.SeatsUnavailableError
		if errors.As(err, &seatsUnavailableError) {
			c.Status(http.StatusConflict)	
			return
		} else if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		c.String(http.StatusOK, createdBookingID.String())
	})
}
