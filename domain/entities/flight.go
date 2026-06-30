// Package entities houses the domain objects
package entities

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"booking.engine/domain/contracts"
)

type FlightFactory struct {
	flightRepository contracts.FlightRepository
}

func NewFlightFactory(flightRepository contracts.FlightRepository) FlightFactory {
	factory := FlightFactory{flightRepository: flightRepository}
	return factory
}

type Flight struct {
	flightRespository contracts.FlightRepository

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

	journey.AllocateSeats(ctx, flight, result.ObtainedSeatLockIDs)

	if !result.ValidFlightID {
		return false, &FlightIDNotFoundError{FlightID: flight.ID}
	}

	return result.SeatsAvailable, nil
}

func (flight Flight) ReleaseSeats(ctx context.Context, seatLockIDs []int) {
	flight.flightRespository.ReleaseSeats(ctx, flight.ID, seatLockIDs)
}
