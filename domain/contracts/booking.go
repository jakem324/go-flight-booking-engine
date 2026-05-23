// Package contracts houses the interfaces which the corresponding entities will use to interface with
// components outside of the domain boundary
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

type InitializeBookingDto struct {
	NumberOfPassengers int
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
	InitializeBooking(ctx context.Context, dto InitializeBookingDto) (uuid.UUID, error)
	ValidateBooking(ctx context.Context, ID uuid.UUID) (ValidateBookingResult, error)
	// "Event" language intentional: the implemented repository can choose to either write incrementally
	// using these events, or wait until OnChangesCompleted to write the whole aggregate once at the end.
	// The entity does not care which option is leveraged; it is simply letting the repository know when
	// these state changes are being made. This keeps the entity agnostic of the architectural requirements.
	OnSeatsAllocated(ctx context.Context, bookingID uuid.UUID, isInboundJourney bool, flightID uuid.UUID, seatLockIDs []int) error
	OnSeatsDeallocated(ctx context.Context, bookingID uuid.UUID, isInboundJourney bool)
	OnChangesCompleted(ctx context.Context, changes BookingChanges) error
}
