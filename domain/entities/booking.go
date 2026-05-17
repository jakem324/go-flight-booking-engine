package entities

type BookingRepository interface {
	CreateBooking() (Booking, error)
}

type Journey struct {
	Parent *Booking
}

type Booking struct {
	Outbound Journey
	Return Journey
}

func (journey Journey) ReleaseAllSeats() {

}
