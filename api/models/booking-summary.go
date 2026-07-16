package models

import (
	"time"

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

