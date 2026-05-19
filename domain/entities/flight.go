// Package entities houses the domain objects
package entities

import "github.com/google/uuid"

type SeatLockResult struct {
	LockIDs []int
	RequestedSeatsAvailable bool
	Error error
}

type FlightRepository interface {
	LockSeats(flightID uuid.UUID, numberOfSeats int) ([]int, error)
	ReleaseSeats(flightID uuid.UUID, seatLockIDs []int) 
}

type Flight struct {
	flightRespository FlightRepository

	ID uuid.UUID
}

func NewFlight(id uuid.UUID) Flight {
	return Flight{ ID: id }
}

func (flight Flight) TryBookSeats(journey *Journey, requiredSeats int) (bool, error) {
	obtainedSeatLocks, err := flight.flightRespository.LockSeats(flight.ID, requiredSeats)
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

func (flight Flight) ReleaseSeats(seatLockIDs []int) {
	flight.flightRespository.ReleaseSeats(flight.ID, seatLockIDs)
}

