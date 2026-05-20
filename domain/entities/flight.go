// Package entities houses the domain objects
package entities

import "log"
import "github.com/google/uuid"

type FlightFactory struct {
	flightRepository FlightRepository
}

func NewFlightFactory(flightRepository FlightRepository) FlightFactory {
	factory := FlightFactory{ flightRepository: flightRepository }
	return factory
}

type FlightRepository interface {
	LockSeats(flightID uuid.UUID, numberOfSeats int) ([]int, error)
	ReleaseSeats(flightID uuid.UUID, seatLockIDs []int) 
}

type Flight struct {
	flightRespository FlightRepository

	ID uuid.UUID
}

func (factory FlightFactory) NewFlight(id uuid.UUID) Flight {
	return Flight{ ID: id, flightRespository: factory.flightRepository }
}

func (flight Flight) TryBookSeats(journey *Journey) (bool, error) {
	log.Printf("TBS %v", journey.Parent.numberOfPassengers)
	obtainedSeatLocks, err := flight.flightRespository.LockSeats(flight.ID, journey.Parent.numberOfPassengers)
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

