// Package entities houses the domain objects
package entities

import "context"
import "github.com/google/uuid"

type FlightFactory struct {
	flightRepository FlightRepository
}

func NewFlightFactory(flightRepository FlightRepository) FlightFactory {
	factory := FlightFactory{ flightRepository: flightRepository }
	return factory
}

type FlightRepository interface {
	LockSeats(ctx context.Context, flightID uuid.UUID, numberOfSeats int) ([]int, error)
	ReleaseSeats(ctx context.Context, flightID uuid.UUID, seatLockIDs []int) 
}

type Flight struct {
	flightRespository FlightRepository

	ID uuid.UUID
}

func (factory FlightFactory) NewFlight(id uuid.UUID) Flight {
	return Flight{ ID: id, flightRespository: factory.flightRepository }
}

func (flight Flight) TryBookSeats(ctx context.Context, journey *Journey) (bool, error) {
	obtainedSeatLocks, err := flight.flightRespository.LockSeats(ctx, flight.ID, journey.Parent.numberOfPassengers)
	if err != nil {
		return false, err
	}
	
	if obtainedSeatLocks == nil {
		return false, nil
	}

	err = journey.AllocateSeats(ctx, flight, obtainedSeatLocks)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (flight Flight) ReleaseSeats(ctx context.Context, seatLockIDs []int) {
	flight.flightRespository.ReleaseSeats(ctx, flight.ID, seatLockIDs)
}

