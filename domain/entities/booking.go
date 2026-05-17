package entities

type CreateBookingResult struct {
	Booking Booking
	Error error
}

type BookingRepository interface {
	CreateBooking() chan CreateBookingResult
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
