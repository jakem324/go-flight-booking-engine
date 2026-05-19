package entities

import "github.com/google/uuid"

type BookingRepository interface {
	CreateBooking() (Booking, error)
	GetBooking(id uuid.UUID) (Booking, error)
	AllocateSeats(bookingID uuid.UUID, isInboundJourney bool, flightID uuid.UUID, seatLockIDs []int) error
	DeallocateSeats(bookingID uuid.UUID, isInboundJourney bool)
	UpsertBooking(*Booking) error
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
	Return Journey
}

func (journey *Journey) ReleaseAllSeats() {
	journey.Parent.bookingRepository.DeallocateSeats(journey.Parent.ID, journey.isInboundJourney)
	for _, leg := range journey.Legs {
		flight := NewFlight(leg.FlightID)
		flight.ReleaseSeats(leg.SeatLockIDs)
	}
}

func (booking *Booking) FinalizeChanges () error {
	err := booking.bookingRepository.UpsertBooking(booking)
	if err != nil {
		return err
	}

	return nil
}

func (journey *Journey) AllocateSeats(flight Flight, seatLockIDs []int) error {
	err := journey.Parent.bookingRepository.AllocateSeats(
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

