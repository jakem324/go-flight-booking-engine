package commands

import (
	"strings"
	"errors"
	"testing"
	"github.com/google/uuid"
	"booking.engine/domain/entities"
)

type BookingRepositoryMock struct {
	entities.BookingRepository
	initializeBookingIDFn func() (uuid.UUID, error)
	validateBookingIDFn func(ID uuid.UUID) (bool, error)
	onSeatsAllocatedFn func(bookingID uuid.UUID, isInboundJourney bool, flightID uuid.UUID, seatLockIDs []int) error
	onSeatsDeallocatedFn func(bookingID uuid.UUID, isInboundJourney bool)
	onChangesCompletedFn func(entities.BookingChanges) error
}

func (m BookingRepositoryMock) InitializeBookingID() (uuid.UUID, error) {
	return m.initializeBookingIDFn()
}

func CommandWithZeroRequiredSeatsIsRejected(t* testing.T) {
	mock := BookingRepositoryMock{}
	factory := entities.NewBookingFactory(mock)

	dto := CreatePencilBookingDto{ RequiredNumberOfSeats: 0 }
	expected := "invalid number of passengers"

	_, err := CreatePencilBooking(factory, dto)
	if err == nil || !strings.Contains(err.Error(), expected) {
		t.Errorf("expected error containing %q, got %v", expected, err)	
	}	
}
/*
	mock := BookingRepositoryMock{
		initializeBookingIDFn: func() (uuid.UUID, error) {
			return uuid.Nil, errors.New("")
		},
	}
*/
