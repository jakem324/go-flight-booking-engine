package repositories

import "errors"
import "github.com/google/uuid"

type FlightRepository struct {}

func (flightRepository FlightRepository) LockSeats(flightID uuid.UUID, numberOfSeats int) ([]int, error) {
	return []int{}, errors.New("not implemented")
}

func (flightRepository FlightRepository) ReleaseSeats(flightID uuid.UUID, seatLockIDs []int) {
	// Fire-and-forget; failures unimportant
}

