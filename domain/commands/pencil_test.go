package commands

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/google/uuid"
	"booking.engine/domain/entities"
)

type BookingRepositoryMock struct {
	mock.Mock
	entities.BookingRepository
}

func (m *BookingRepositoryMock) InitializeBookingID() (uuid.UUID, error) {
	args := m.Called()
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func TestCreatePencilBooking_CommandWithZeroRequiredSeatsIsRejected(t* testing.T) {
	mock := new(BookingRepositoryMock)
	factory := entities.NewBookingFactory(mock)

	dto := CreatePencilBookingDto{ RequiredNumberOfSeats: 0 }
	_, err := CreatePencilBooking(factory, dto)

	if assert.NotNil(t, err) {
		assert.Equal(t, "invalid number of passengers", err.Error())
	}

	mock.AssertNotCalled(t, "InitializeBooking")
	mock.AssertNotCalled(t, "OnChangesCompleted")
}

func TestCreatePencilBooking_BookingIsInitialized(t* testing.T) {
	mock := new(BookingRepositoryMock)

	bookingID := uuid.New()
	expectedInitializationDto := entities.InitializeBookingDto{
		NumberOfPassengers: 5,
	}
	mock.On("InitializeBooking", expectedInitializationDto).Return(bookingID, nil)

	factory := entities.NewBookingFactory(mock)
	dto := CreatePencilBookingDto{ RequiredNumberOfSeats: 5 }
	_, err := CreatePencilBooking(factory, dto)

	if assert.Nil(t, err) {
		mock.AssertCalled(t, "InitializeBooking", expectedInitializationDto)
	}
}

