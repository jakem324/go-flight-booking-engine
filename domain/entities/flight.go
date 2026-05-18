package entities

import "github.com/google/uuid"

type SeatLockResult struct {
	LockIds []int
	RequestedSeatsAvailable bool
	Error error
}

type FlightRepository interface {
	LockSeats(flightId uuid.UUID, numberOfSeats int) ([]int, error)
}

type Flight struct {
	flightRespository FlightRepository

	Id uuid.UUID
}

func NewFlight(id uuid.UUID) Flight {
	return Flight{ Id: id }
}

func (flight Flight) TryBookSeats(journey *Journey, requiredSeats int) (bool, error) {
	obtainedSeatLocks, err := flight.flightRespository.LockSeats(flight.Id, requiredSeats)
	if err != nil {
		return false, err
	}
	
	if obtainedSeatLocks == nil {
		return false, nil
	}

	err = journey.AllocateSeats(flight, obtainedSeatLocks)
	if err != nil {
		return false, err
	}

	return true, nil
}
