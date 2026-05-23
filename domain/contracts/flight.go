package contracts

import (
	"context"

	"github.com/google/uuid"
)

type SeatLockResult struct {
	ValidFlightID       bool
	SeatsAvailable      bool
	ObtainedSeatLockIDs []int
}

type FlightRepository interface {
	LockSeats(ctx context.Context, flightID uuid.UUID, numberOfSeats int) (SeatLockResult, error)
	ReleaseSeats(ctx context.Context, flightID uuid.UUID, seatLockIDs []int)
}
