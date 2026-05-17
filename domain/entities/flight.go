package entities

import "github.com/google/uuid"

type SeatLockResult struct {
	LockIds []int
	RequestedSeatsAvailable bool
	Error error
}

type FlightRepository interface {
	CreateSeatLock(flightId uuid.UUID, numberOfSeats int) <- chan SeatLockResult
}

type Flight struct {
	flightRespository FlightRepository

	Id uuid.UUID
}

func NewFlight(id uuid.UUID) Flight {
	return Flight{ Id: id }
}

type SeatAllocationResult struct {
	Available bool
	Error error
}

func (flight Flight) TryAllocateSeat(*Journey) chan SeatAllocationResult {
	return nil
}
