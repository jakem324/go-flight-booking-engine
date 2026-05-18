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

func (flight Flight) TryAllocateSeats(journey *Journey, requiredSeats int) (bool, error) {
	obtainedSeats, err := flight.flightRespository.CreateSeatLock(flight.Id, requiredSeats)
	if err != nil {
		return false, err
	}
	
	if obtainedSeats == nil {
		return false, nil
	}

	err = journey.AllocateSeats(flight, obtainedSeats)
	if err != nil {
		return false, err
	}

	return true, nil
}
