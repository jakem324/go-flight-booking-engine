package entities

import (
	"context"
	"log"

	"github.com/google/uuid"

	"booking.engine/domain/contracts"
)

type BookingFactory struct {
	bookingRepository contracts.BookingRepository
	flightFactory     *FlightFactory
}

func NewBookingFactory(bookingRepository contracts.BookingRepository, flightFactory FlightFactory) BookingFactory {
	factory := BookingFactory{bookingRepository: bookingRepository, flightFactory: &flightFactory}
	return factory
}

func (factory BookingFactory) NewBooking(ctx context.Context, numberOfPassengers int) (*Booking, error) {
	id, err := factory.bookingRepository.InitializeBooking(ctx, contracts.InitializeBookingDto{
		NumberOfPassengers: numberOfPassengers,
	})

	if err != nil {
		return &Booking{}, err
	}

	booking, err := factory.constructBooking(id, numberOfPassengers)
	return booking, err
}

func (factory BookingFactory) ExistingBooking(ctx context.Context, ID uuid.UUID) (*Booking, error) {
	result, err := factory.bookingRepository.ValidateBooking(ctx, ID)
	if !result.BookingExists || err != nil {
		return nil, err
	}

	booking, err := factory.constructBooking(ID, result.NumberOfPassengers)
	return booking, err
}

func (factory BookingFactory) constructBooking(ID uuid.UUID, numberOfPassengers int) (*Booking, error) {
	booking := Booking{}
	booking.bookingRepository = factory.bookingRepository
	booking.flightFactory = factory.flightFactory
	booking.ID = ID
	booking.numberOfPassengers = numberOfPassengers
	booking.Outbound = Journey{
		Parent: &booking,
	}
	booking.Inbound = Journey{
		Parent: &booking,
	}
	booking.Inbound.isInboundJourney = true
	return &booking, nil
}

type Journey struct {
	Parent           *Booking
	legs             []contracts.JourneyLeg
	isInboundJourney bool
	modified         bool
}

type Booking struct {
	bookingRepository contracts.BookingRepository
	flightFactory     *FlightFactory

	ID                 uuid.UUID
	numberOfPassengers int
	Outbound           Journey
	Inbound            Journey
}

func (journey *Journey) ReleaseAllSeats(ctx context.Context) {
	log.Printf("ReleaseAllSeats %v", journey.Parent.flightFactory)
	journey.Parent.bookingRepository.OnSeatsDeallocated(ctx, journey.Parent.ID, journey.isInboundJourney)
	for _, leg := range journey.legs {
		flight := journey.Parent.flightFactory.NewFlight(leg.FlightID)
		flight.ReleaseSeats(ctx, leg.SeatLockIDs)
	}
	journey.legs = nil
	journey.modified = true
}

func (journey *Journey) AllocateSeats(ctx context.Context, flight Flight, seatLockIDs []int) error {
	err := journey.Parent.bookingRepository.OnSeatsAllocated(
		ctx,
		journey.Parent.ID,
		journey.isInboundJourney,
		flight.ID,
		seatLockIDs)

	if err != nil {
		return err
	}

	journey.legs = append(journey.legs, contracts.JourneyLeg{FlightID: flight.ID, SeatLockIDs: seatLockIDs})
	journey.modified = true

	return nil
}

func (booking *Booking) FinalizeChanges(ctx context.Context) error {
	stagedChanges := contracts.BookingChanges{
		ID:                 booking.ID,
		NumberOfPassengers: booking.numberOfPassengers,
	}

	if booking.Inbound.modified {
		stagedChanges.InboundLegs = booking.Inbound.legs
	}

	if booking.Outbound.modified {
		stagedChanges.OutboundLegs = booking.Outbound.legs
	}

	err := booking.bookingRepository.OnChangesCompleted(ctx, stagedChanges)

	if err != nil {
		return err
	}

	return nil
}
