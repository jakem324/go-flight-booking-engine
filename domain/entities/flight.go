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

type SeatLockResult struct {
	ValidFlightID bool
	SeatsAvailable bool
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
	return Flight{ ID: id, flightRespository: factory.flightRepository }
}

type TryBookSeatsOutcome struct {
	FlightIDFound bool
	SeatsObtained bool
}

func (flight Flight) TryBookSeats(ctx context.Context, journey *Journey) (TryBookSeatsOutcome, error) {
	result, err := flight.flightRespository.LockSeats(ctx, flight.ID, journey.Parent.numberOfPassengers)
	if err != nil {
		return TryBookSeatsOutcome{}, err
	}

	err = journey.AllocateSeats(ctx, flight, result.ObtainedSeatLockIDs)
	if err != nil {
		return TryBookSeatsOutcome{}, err
	}

	return TryBookSeatsOutcome{
		FlightIDFound: result.ValidFlightID,
		SeatsObtained: result.SeatsAvailable,
	}, nil
}

func (flight Flight) ReleaseSeats(ctx context.Context, seatLockIDs []int) {
	flight.flightRespository.ReleaseSeats(ctx, flight.ID, seatLockIDs)
}

