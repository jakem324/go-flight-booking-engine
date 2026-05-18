package entities

import "github.com/google/uuid"

type BookingRepository interface {
	CreateBooking() (Booking, error)
	GetBooking(id uuid.UUID) (Booking, error)
}

type Journey struct {
	Parent *Booking
}

type Booking struct {
	bookingRepository BookingRepository

	Id uuid.UUID
	NumberOfPassengers int
	Outbound Journey
	Return Journey
}

func (journey Journey) ReleaseAllSeats() {

}

func (booking Booking) FinalizeChanges () {

}

func (journey Journey) AllocateSeats(flight Flight, seats []int) error {

}
