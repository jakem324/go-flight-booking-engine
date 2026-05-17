package entities

import "github.com/google/uuid"

type SeatLockResult struct {
	LockIds []int
	RequestedSeatsAvailable bool
	Error error
}

type FlightRepository interface {
	CreateSeatLock(flightId uuid.UUID, numberOfSeats int) ([]int, error)
}

type Flight struct {
	flightRespository FlightRepository

	Id uuid.UUID
}

func NewFlight(id uuid.UUID) Flight {
	return Flight{ Id: id }
}

func (flight Flight) TryAllocateSeat(*Journey) (bool, error) {
	return false, nil
}
