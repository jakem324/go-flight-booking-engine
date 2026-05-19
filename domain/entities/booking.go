package entities

import "github.com/google/uuid"

type BookingRepository interface {
	InitializeBookingID() (uuid.UUID, error)
	ValidateBookingID(ID uuid.UUID) (bool, error)
	// "Event" language intentional: the implemented repository can choose to either write incrementally
	// using these events, or wait until OnChangesCompleted to write the whole aggregate once at the end.
	// The entity does not care which option is leveraged; it is simply letting the repository know when
	// these state changes are being made. This keeps the entity agnostic of the architectural requirements.
	OnSeatsAllocated(bookingID uuid.UUID, isInboundJourney bool, flightID uuid.UUID, seatLockIDs []int) error
	OnSeatsDeallocated(bookingID uuid.UUID, isInboundJourney bool)
	OnChangesCompleted(*Booking) error
}

type JourneyLeg struct {
	FlightID uuid.UUID
	SeatLockIDs []int
}

type Journey struct {
	Parent *Booking
	Legs []JourneyLeg
	isInboundJourney bool
}

type Booking struct {
	bookingRepository BookingRepository

	ID uuid.UUID
	NumberOfPassengers int
	Outbound Journey
	Inbound Journey
}

func NewBooking() (*Booking, error) {
	booking := Booking{}
	id, err := booking.bookingRepository.InitializeBookingID()
	if err != nil {
		return nil, err
	}

	booking.ID = id
	booking.Inbound.isInboundJourney = true
	return &booking, nil
}

func ExistingBooking(ID uuid.UUID) (*Booking, error) {
	booking := Booking{}
	valid, err := booking.bookingRepository.ValidateBookingID(ID)
	if !valid || err != nil {
		return nil, err
	}

	booking.ID = ID
	booking.Inbound.isInboundJourney = true
	return &booking, nil
}

func (journey *Journey) ReleaseAllSeats() {
	journey.Parent.bookingRepository.OnSeatsDeallocated(journey.Parent.ID, journey.isInboundJourney)
	for _, leg := range journey.Legs {
		flight := NewFlight(leg.FlightID)
		flight.ReleaseSeats(leg.SeatLockIDs)
	}
	journey.Legs = nil
}

func (booking *Booking) FinalizeChanges () error {
	err := booking.bookingRepository.OnChangesCompleted(booking)
	if err != nil {
		return err
	}

	return nil
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

	journey.Legs = append(journey.Legs, JourneyLeg{ FlightID: flight.ID, SeatLockIDs: seatLockIDs })

	return nil
}

