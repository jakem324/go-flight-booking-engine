// Package queries provides interfaces for read-only querying. This bypasses the domain entities entirely,
// so that fetching data for the presentation layer may be treated as an altogether separate concern.
package queries

import (
	"time"
	"context"
	"github.com/google/uuid"
)

type JourneyLeg struct {
	DepartureDate        time.Time `json:"departureDate"`
	DepartureAirportCode string    `json:"departureAirportCode"`
	DepartureTerminal    string    `json:"departureTerminal"`
	ArrivalDate          time.Time `json:"arrivalDate"`
	ArrivalAirportCode   string    `json:"arrivalAirportCode"`
	ArrivalTerminal      string    `json:"arrivalTerminal"`
}

type BookingSummary struct {
	BookingID           uuid.UUID    `json:"bookingId"`
	NumberOfPassengers  int          `json:"numberOfPassengers"`
	OutboundJourneyLegs []JourneyLeg `json:"outboundJourneyLegs"`
	InboundJourneyLegs  []JourneyLeg `json:"inboundJourneyLegs"`
}

type BookingQueries interface {
	GetBookingSummary(ctx context.Context, ID uuid.UUID) (*BookingSummary, error)
}

