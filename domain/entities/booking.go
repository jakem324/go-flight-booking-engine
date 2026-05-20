package entities

import (
	"errors"
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

func (factory BookingFactory) NewBooking(numberOfPassengers int) (Booking, error) {
	if numberOfPassengers < 1 {
		return Booking{}, errors.New("invalid number of passengers")
	}
	booking := Booking{}
	booking.bookingRepository = factory.bookingRepository
	booking.flightFactory = factory.flightFactory
	id, err := booking.bookingRepository.InitializeBooking(InitializeBookingDto{
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

func (factory BookingFactory) ExistingBooking(ID uuid.UUID) (*Booking, error) {
	booking := Booking{}
	booking.bookingRepository = factory.bookingRepository
	result, err := booking.bookingRepository.ValidateBooking(ID)
	if !result.BookingExists || err != nil {
		return nil, err
	}

	booking.Outbound = Journey{}
	booking.Inbound = Journey{}
	booking.Inbound.isInboundJourney = true
	booking.ID = ID
	booking.numberOfPassengers = result.NumberOfPassengers
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
	InitializeBooking(dto InitializeBookingDto) (uuid.UUID, error)
	ValidateBooking(ID uuid.UUID) (ValidateBookingResult, error)
	// "Event" language intentional: the implemented repository can choose to either write incrementally
	// using these events, or wait until OnChangesCompleted to write the whole aggregate once at the end.
	// The entity does not care which option is leveraged; it is simply letting the repository know when
	// these state changes are being made. This keeps the entity agnostic of the architectural requirements.
	OnSeatsAllocated(bookingID uuid.UUID, isInboundJourney bool, flightID uuid.UUID, seatLockIDs []int) error
	OnSeatsDeallocated(bookingID uuid.UUID, isInboundJourney bool)
	OnChangesCompleted(BookingChanges) error
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

func (journey *Journey) ReleaseAllSeats() {
	journey.Parent.bookingRepository.OnSeatsDeallocated(journey.Parent.ID, journey.isInboundJourney)
	for _, leg := range journey.legs {
		flight := journey.Parent.flightFactory.NewFlight(leg.FlightID)
		flight.ReleaseSeats(leg.SeatLockIDs)
	}
	journey.legs = nil
	journey.modified = true
}

func (journey *Journey) AllocateSeats(flight Flight, seatLockIDs []int) error {
	err := journey.Parent.bookingRepository.OnSeatsAllocated(
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

func (booking *Booking) FinalizeChanges () error {
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
	
	err := booking.bookingRepository.OnChangesCompleted(stagedChanges);

	if err != nil {
		return err
	}

	return nil
}
