// Package contracts defines the interfaces between the domain entities and any dependencies external to the
// domain boundary
package contracts

import (
	"context"

	"github.com/google/uuid"
)

type BookingChanges struct {
	ID                 uuid.UUID
	NumberOfPassengers int
	InboundLegs        []JourneyLeg
	OutboundLegs       []JourneyLeg
}

type ValidateBookingResult struct {
	BookingExists      bool
	NumberOfPassengers int
}

type JourneyLeg struct {
	FlightID    uuid.UUID
	SeatLockIDs []int
}

type BookingRepository interface {
	InitializeBooking(ctx context.Context) (uuid.UUID, error)
	ValidateBooking(ctx context.Context, ID uuid.UUID) (ValidateBookingResult, error)
	SaveBooking(ctx context.Context, changes BookingChanges) error
}
