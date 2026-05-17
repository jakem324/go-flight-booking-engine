package entities

import "github.com/google/uuid"

type BookingRepository interface {
	CreateBooking() (Booking, error)
}

type Journey struct {
	Parent *Booking
}

type Booking struct {
	Id uuid.UUID
	Outbound Journey
	Return Journey
}

func (journey Journey) ReleaseAllSeats() {

}

func (booking Booking) Finalize () {

}
