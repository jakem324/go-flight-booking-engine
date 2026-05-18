package entities

import "github.com/google/uuid"

type BookingRepository interface {
	CreateBooking() (Booking, error)
	GetBooking(id uuid.UUID) (Booking, error)
	AllocateSeats(bookingId uuid.UUID, isInboundJourney bool, flightId uuid.UUID, seatLockIds []int) error
	DeallocateSeats(bookingId uuid.UUID, isInboundJourney bool)
	UpsertBooking(Booking) error
}

type Journey struct {
	Parent *Booking
	isInboundJourney bool
}

type Booking struct {
	bookingRepository BookingRepository

	Id uuid.UUID
	NumberOfPassengers int
	Outbound Journey
	Return Journey
}

func (journey Journey) ReleaseAllSeats() {
	journey.Parent.bookingRepository.DeallocateSeats(journey.Parent.Id, journey.isInboundJourney)
}

func (booking Booking) FinalizeChanges () error {
	err := booking.bookingRepository.UpsertBooking(booking)
	if err != nil {
		return err
	}

	return nil
}

func (journey Journey) AllocateSeats(flight Flight, seatLockIds []int) error {
	err := journey.Parent.bookingRepository.AllocateSeats(
		journey.Parent.Id,
		journey.isInboundJourney,
		flight.Id,
		seatLockIds)
	
	if err != nil {
		return err
	}

	return nil
}

