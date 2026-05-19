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

func (m *BookingRepositoryMock) InitializeBooking(dto entities.InitializeBookingDto) (uuid.UUID, error) {
	args := m.Called(dto)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *BookingRepositoryMock) ValidateBooking(ID uuid.UUID) (entities.ValidateBookingResult, error) {
	args := m.Called(ID)
	return args.Get(0).(entities.ValidateBookingResult), args.Error(1)
}

func (m *BookingRepositoryMock) OnSeatsAllocated(bookingID uuid.UUID, isInboundJourney bool, flightID uuid.UUID, seatLockIDs []int) error {
	args := m.Called(bookingID, isInboundJourney, flightID, seatLockIDs)
	return args.Error(0)
}

func (m *BookingRepositoryMock) OnSeatsDeallocated(bookingID uuid.UUID, isInboundJourney bool) {
	m.Called(bookingID, isInboundJourney)
}

func (m *BookingRepositoryMock) OnChangesCompleted(changes entities.BookingChanges) error {
	args := m.Called(changes)
	return args.Error(0)
}

type FlightRepositoryMock struct {
	mock.Mock
	entities.FlightRepository
}

func (m *FlightRepositoryMock) LockSeats(flightID uuid.UUID, numberOfSeats int) ([]int, error) {
	args := m.Called(flightID, numberOfSeats)
	return args.Get(0).([]int), args.Error(1)
}

func (m *FlightRepositoryMock) ReleaseSeats(flightID uuid.UUID, seatLockIDs []int) {
	m.Called(flightID, seatLockIDs)
}





func TestCreatePencilBooking_CommandWithZeroRequiredSeatsIsRejected(t* testing.T) {
	bookingRepositoryMock := new(BookingRepositoryMock)
	flightRepositoryMock := new(FlightRepositoryMock)
	flightFactory := entities.NewFlightFactory(flightRepositoryMock)
	bookingFactory := entities.NewBookingFactory(bookingRepositoryMock, flightFactory)

	dto := CreatePencilBookingDto{ RequiredNumberOfSeats: 0 }
	_, err := CreatePencilBooking(bookingFactory, flightFactory, dto)

	if assert.NotNil(t, err) {
		assert.Equal(t, "invalid number of passengers", err.Error())
	}

	bookingRepositoryMock.AssertNotCalled(t, "InitializeBooking")
	bookingRepositoryMock.AssertNotCalled(t, "OnChangesCompleted")
}

/*
func TestCreatePencilBooking_BookingIsInitialized(t* testing.T) {
	bookingRepositoryMock := new(BookingRepositoryMock)
	flightRepositoryMock := new(FlightRepositoryMock)

	bookingID := uuid.New()
	expectedInitializationDto := entities.InitializeBookingDto{
		NumberOfPassengers: 5,
	}
	bookingRepositoryMock.On("InitializeBooking", expectedInitializationDto).Return(bookingID, nil)
	bookingRepositoryMock.On("ValidateBooking", bookingID).Return(entities.ValidateBookingResult{NumberOfPassengers: 5}, nil)
	flightRepositoryMock.On(
		"LockSeats",
		mock.AnythingOfType("uuid.UUID"),
		mock.AnythingOfType("int")).Return([]int {1,2,3}, nil)
	//mock.On("OnSeatsAllocated", bookingID, false).Return(entities.ValidateBookingResult{NumberOfPassengers: 5}, nil)

	flightFactory := entities.NewFlightFactory(flightRepositoryMock)
	bookingFactory := entities.NewBookingFactory(bookingRepositoryMock, flightFactory)

	dto := CreatePencilBookingDto{ RequiredNumberOfSeats: 5 }
	_, err := CreatePencilBooking(factory, dto)

	if assert.Nil(t, err) {
		mock.AssertCalled(t, "InitializeBooking", expectedInitializationDto)
	}
}
*/
