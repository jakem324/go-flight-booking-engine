// Package repositories contains all postgres repository implementations
package repositories

import "errors"
import "github.com/google/uuid"
import "booking.engine/domain/entities"

type BookingRepository struct {}

func (bookingRepository BookingRepository) InitializeBooking(
	dto entities.InitializeBookingDto,
) (uuid.UUID, error) {
	return uuid.Nil, errors.New("not implemented")
}

func (bookingRepository BookingRepository) ValidateBooking(
	ID uuid.UUID,
) (entities.ValidateBookingResult, error) {
	return entities.ValidateBookingResult{}, errors.New("not implemented")
}

func (bookingRepository BookingRepository) OnSeatsAllocated(
	bookingID uuid.UUID,
	isInboundJourney bool,
	flightID uuid.UUID,
	seatLockIDs []int) error {
	return errors.New("not implemented")
}

func (bookingRepository BookingRepository) OnSeatsDeallocated(
	bookingID uuid.UUID,
	isInboundJourney bool) {
	// Fire-and-forget; failures unimportant
}

func (bookingRepository BookingRepository)OnChangesCompleted(
	entities.BookingChanges,
) error {
	// Committing at the end for SQL DB implementation rendered unnecessary by incremental writing
	return nil
}

