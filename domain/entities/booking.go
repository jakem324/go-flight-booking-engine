package entities

import (
	"context"
	"errors"
	"log"
	"github.com/google/uuid"
)

type BookingFactory struct {
	bookingRepository BookingRepository
	flightFactory *FlightFactory
}

func NewBookingFactory(bookingRepository BookingRepository, flightFactory FlightFactory) BookingFactory {
	factory := BookingFactory{ bookingRepository: bookingRepository, flightFactory: &flightFactory }
	return factory
}

func (factory BookingFactory) NewBooking(ctx context.Context, numberOfPassengers int) (Booking, error) {
	if numberOfPassengers < 1 {
		return Booking{}, errors.New("invalid number of passengers")
	}
	booking := Booking{}
	booking.bookingRepository = factory.bookingRepository
	booking.flightFactory = factory.flightFactory
	id, err := booking.bookingRepository.InitializeBooking(ctx, InitializeBookingDto{
		NumberOfPassengers: numberOfPassengers,
	})

	if err != nil {
		return Booking{}, err
	}

	booking.ID = id
	booking.numberOfPassengers = numberOfPassengers
	booking.Outbound = Journey{
		Parent: &booking,
	}
	booking.Inbound = Journey{
		Parent: &booking,
	}
	booking.Inbound.isInboundJourney = true
	return booking, nil
}

func (factory BookingFactory) ExistingBooking(ctx context.Context, ID uuid.UUID) (*Booking, error) {
	booking := Booking{}
	booking.bookingRepository = factory.bookingRepository
	booking.flightFactory = factory.flightFactory
	result, err := booking.bookingRepository.ValidateBooking(ctx, ID)
	if !result.BookingExists || err != nil {
		return nil, err
	}

	booking.ID = ID
	booking.numberOfPassengers = result.NumberOfPassengers
	booking.Outbound = Journey{
		Parent: &booking,
	}
	booking.Inbound = Journey{
		Parent: &booking,
	}
	booking.Inbound.isInboundJourney = true
	return &booking, nil
}

type BookingChanges struct {
	ID uuid.UUID
	NumberOfPassengers int
	InboundLegs []JourneyLeg
	OutboundLegs []JourneyLeg
}

type InitializeBookingDto struct {
	NumberOfPassengers int
}

type ValidateBookingResult struct {
	BookingExists bool
	NumberOfPassengers int
}

type BookingRepository interface {
	InitializeBooking(ctx context.Context, dto InitializeBookingDto) (uuid.UUID, error)
	ValidateBooking(ctx context.Context, ID uuid.UUID) (ValidateBookingResult, error)
	// "Event" language intentional: the implemented repository can choose to either write incrementally
	// using these events, or wait until OnChangesCompleted to write the whole aggregate once at the end.
	// The entity does not care which option is leveraged; it is simply letting the repository know when
	// these state changes are being made. This keeps the entity agnostic of the architectural requirements.
	OnSeatsAllocated(ctx context.Context, bookingID uuid.UUID, isInboundJourney bool, flightID uuid.UUID, seatLockIDs []int) error
	OnSeatsDeallocated(ctx context.Context, bookingID uuid.UUID, isInboundJourney bool)
	OnChangesCompleted(ctx context.Context, changes BookingChanges) error
}

type JourneyLeg struct {
	FlightID uuid.UUID
	SeatLockIDs []int
}

type Journey struct {
	Parent *Booking
	legs []JourneyLeg
	isInboundJourney bool
	modified bool
}

type Booking struct {
	bookingRepository BookingRepository
	flightFactory *FlightFactory

	ID uuid.UUID
	numberOfPassengers int
	Outbound Journey
	Inbound Journey
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

	journey.legs = append(journey.legs, JourneyLeg{ FlightID: flight.ID, SeatLockIDs: seatLockIDs })
	journey.modified = true

	return nil
}

func (booking *Booking) FinalizeChanges (ctx context.Context) error {
	stagedChanges := BookingChanges {
		ID: booking.ID,
		NumberOfPassengers: booking.numberOfPassengers,
	}

	if booking.Inbound.modified {
		stagedChanges.InboundLegs = booking.Inbound.legs
	}

	if booking.Outbound.modified {
		stagedChanges.OutboundLegs = booking.Outbound.legs
	}
	
	err := booking.bookingRepository.OnChangesCompleted(ctx, stagedChanges);

	if err != nil {
		return err
	}

	return nil
}
