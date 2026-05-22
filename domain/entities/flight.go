// Package entities houses the domain objects
package entities

import "context"
import "fmt"
import "github.com/google/uuid"

type FlightFactory struct {
	flightRepository FlightRepository
}

func NewFlightFactory(flightRepository FlightRepository) FlightFactory {
	factory := FlightFactory{flightRepository: flightRepository}
	return factory
}

type SeatLockResult struct {
	ValidFlightID       bool
	SeatsAvailable      bool
	ObtainedSeatLockIDs []int
}

type FlightRepository interface {
	LockSeats(ctx context.Context, flightID uuid.UUID, numberOfSeats int) (SeatLockResult, error)
	ReleaseSeats(ctx context.Context, flightID uuid.UUID, seatLockIDs []int)
}

type Flight struct {
	flightRespository FlightRepository

	ID uuid.UUID
}

func (factory FlightFactory) NewFlight(id uuid.UUID) Flight {
	return Flight{ID: id, flightRespository: factory.flightRepository}
}

type FlightIDNotFoundError struct {
	FlightID uuid.UUID
}

func (e *FlightIDNotFoundError) Error() string {
	return fmt.Sprintf("flight ID not found: %v", e.FlightID)
}

func (e *FlightIDNotFoundError) Is(target error) bool {
	_, ok := target.(*FlightIDNotFoundError)
	return ok
}

func (flight Flight) TryBookSeats(ctx context.Context, journey *Journey) (bool, error) {
	result, err := flight.flightRespository.LockSeats(ctx, flight.ID, journey.Parent.numberOfPassengers)
	if err != nil {
		return false, err
	}

	err = journey.AllocateSeats(ctx, flight, result.ObtainedSeatLockIDs)
	if err != nil {
		return false, err
	}

	if !result.ValidFlightID {
		return false, &FlightIDNotFoundError{FlightID: flight.ID}
	}

	return result.SeatsAvailable, nil
}

func (flight Flight) ReleaseSeats(ctx context.Context, seatLockIDs []int) {
	flight.flightRespository.ReleaseSeats(ctx, flight.ID, seatLockIDs)
}
