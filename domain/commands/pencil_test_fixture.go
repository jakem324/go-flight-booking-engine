package commands

import (
	"context"
	"booking.engine/domain/entities"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type BookingRepositoryMock struct {
	mock.Mock
	entities.BookingRepository
}

func (m *BookingRepositoryMock) InitializeBooking(ctx context.Context, dto entities.InitializeBookingDto) (uuid.UUID, error) {
	args := m.Called(dto)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *BookingRepositoryMock) ValidateBooking(ctx context.Context, ID uuid.UUID) (entities.ValidateBookingResult, error) {
	args := m.Called(ID)
	return args.Get(0).(entities.ValidateBookingResult), args.Error(1)
}

func (m *BookingRepositoryMock) OnSeatsAllocated(
	ctx context.Context, 
	bookingID uuid.UUID,
	isInboundJourney bool,
	flightID uuid.UUID,
	seatLockIDs []int) error {

	args := m.Called(bookingID, isInboundJourney, flightID, seatLockIDs)
	return args.Error(0)
}

func (m *BookingRepositoryMock) OnSeatsDeallocated(ctx context.Context, bookingID uuid.UUID, isInboundJourney bool) {
	m.Called(bookingID, isInboundJourney)
}

func (m *BookingRepositoryMock) OnChangesCompleted(ctx context.Context, changes entities.BookingChanges) error {
	args := m.Called(changes)
	return args.Error(0)
}

type FlightRepositoryMock struct {
	mock.Mock
	entities.FlightRepository
}

func (m *FlightRepositoryMock) LockSeats(ctx context.Context, flightID uuid.UUID, numberOfSeats int) ([]int, error) {
	args := m.Called(flightID, numberOfSeats)
	var seats []int
	if args.Get(0) != nil {
		seats = args.Get(0).([]int)
	}
	return seats, args.Error(1)
}

func (m *FlightRepositoryMock) ReleaseSeats(ctx context.Context, flightID uuid.UUID, seatLockIDs []int) {
	m.Called(flightID, seatLockIDs)
}

type Fixture struct {
	bookingRepositoryMock *BookingRepositoryMock
	flightRepositoryMock *FlightRepositoryMock
	
	handler PencilBookingHandler
}

func CreateFixture () Fixture {
	bookingRepositoryMock := new(BookingRepositoryMock)
	flightRepositoryMock := new(FlightRepositoryMock)

	flightFactory := entities.NewFlightFactory(flightRepositoryMock)
	bookingFactory := entities.NewBookingFactory(bookingRepositoryMock, flightFactory)

	handler := NewPencilBookingHandler(bookingFactory, flightFactory)

	return Fixture{
		bookingRepositoryMock: bookingRepositoryMock,
		flightRepositoryMock: flightRepositoryMock,
		handler: handler,
	}
}

