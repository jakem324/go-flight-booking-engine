package commands

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"booking.engine/domain/contracts"
	"booking.engine/domain/entities"
)

type BookingRepositoryMock struct {
	fixture *Fixture
	mock.Mock
	contracts.BookingRepository
}

func (m *BookingRepositoryMock) InitializeBooking(ctx context.Context) (uuid.UUID, error) {
	args := m.Called()
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *BookingRepositoryMock) ValidateBooking(ctx context.Context, ID uuid.UUID) (contracts.ValidateBookingResult, error) {
	args := m.Called(ID)
	return args.Get(0).(contracts.ValidateBookingResult), args.Error(1)
}

func (m *BookingRepositoryMock) SaveBooking(ctx context.Context, changes contracts.BookingChanges) error {
	args := m.Called(changes)
	return args.Error(0)
}

type FlightRepositoryMock struct {
	mock.Mock
	contracts.FlightRepository
	fixture *Fixture
}

func (m *FlightRepositoryMock) LockSeats(ctx context.Context, flightID uuid.UUID, numberOfSeats int) (contracts.SeatLockResult, error) {
	args := m.Called(flightID, numberOfSeats)
	var result contracts.SeatLockResult
	if args.Get(0) != nil {
		result = args.Get(0).(contracts.SeatLockResult)
	}
	return result, args.Error(1)
}

func (m *FlightRepositoryMock) ReleaseSeats(ctx context.Context, flightID uuid.UUID, seatLockIDs []int) {
	m.Called(flightID, seatLockIDs)
}

type Fixture struct {
	bookingRepositoryMock *BookingRepositoryMock
	flightRepositoryMock  *FlightRepositoryMock

	t      *testing.T
	result uuid.UUID
	err    error
}

func CreateFixture(t *testing.T) Fixture {
	bookingRepositoryMock := new(BookingRepositoryMock)
	flightRepositoryMock := new(FlightRepositoryMock)
	fixture := Fixture{
		bookingRepositoryMock: bookingRepositoryMock,
		flightRepositoryMock:  flightRepositoryMock,
		t:                     t,
	}

	fixture.bookingRepositoryMock.fixture = &fixture
	fixture.flightRepositoryMock.fixture = &fixture
	return fixture
}

// Helpers

func (m *BookingRepositoryMock) GivenInitializeBookingMock(bookingID uuid.UUID) {
	m.On("InitializeBooking", mock.Anything).Return(bookingID, nil)
}

func (m *BookingRepositoryMock) GivenValidateBookingMock(bookingID uuid.UUID, numberOfPassengers int) {
	m.On("ValidateBooking", bookingID).Return(
		contracts.ValidateBookingResult{
			BookingExists:      true,
			NumberOfPassengers: numberOfPassengers,
		}, nil)
}

func (m *BookingRepositoryMock) GivenValidateBookingMockWithNegativeResult(bookingID uuid.UUID) {
	m.On("ValidateBooking", bookingID).Return(
		contracts.ValidateBookingResult{
			BookingExists:      false,
			NumberOfPassengers: 0,
		}, nil)
}

func (m *FlightRepositoryMock) GivenLockSeatsMock(flightID uuid.UUID, passengers int, returnSeatLockIDs ...int) {
	m.On(
		"LockSeats",
		flightID,
		passengers,
	).Return(contracts.SeatLockResult{
		ValidFlightID:       true,
		SeatsAvailable:      true,
		ObtainedSeatLockIDs: returnSeatLockIDs}, nil)
}

func (m *FlightRepositoryMock) GivenLockSeatsMockWithUnavailableResult(flightID uuid.UUID, passengers int, returnSeatLockIDs ...int) {
	m.On(
		"LockSeats",
		flightID,
		passengers,
	).Return(contracts.SeatLockResult{
		ValidFlightID:  true,
		SeatsAvailable: false}, nil)
}

func (m *FlightRepositoryMock) GivenReleaseSeatsMock() {
	m.On("ReleaseSeats", mock.Anything, mock.Anything).Return(nil)
}

func (m *BookingRepositoryMock) GivenSaveBookingMock() {
	m.On("SaveBooking", mock.Anything).Return(nil)
}

func (f *Fixture) WhenCreatePencilBookingIsCalled(requiredNumberOfSeats int, outboundLegs ...uuid.UUID) {
	flightFactory := entities.NewFlightFactory(f.flightRepositoryMock)
	bookingFactory := entities.NewBookingFactory(f.bookingRepositoryMock, flightFactory)
	handler := NewPencilBookingHandler(bookingFactory, flightFactory)

	dto := CreatePencilBookingDto{
		RequiredNumberOfSeats: requiredNumberOfSeats,
		OutboundJourneyLegs:   outboundLegs,
	}
	ctx := context.Background()
	result, err := handler.CreatePencilBooking(ctx, dto)

	f.result = result
	f.err = err
}

func (f *Fixture) WhenSetInboundJourneyIsCalled(bookingID uuid.UUID, legs ...uuid.UUID) {
	flightFactory := entities.NewFlightFactory(f.flightRepositoryMock)
	bookingFactory := entities.NewBookingFactory(f.bookingRepositoryMock, flightFactory)
	handler := NewPencilBookingHandler(bookingFactory, flightFactory)

	dto := SetInboundJourneyDto{
		BookingID:          bookingID,
		InboundJourneyLegs: legs,
	}
	ctx := context.Background()
	err := handler.SetInboundJourney(ctx, dto)
	f.err = err
}

func (f *Fixture) HandlerShouldCompleteSuccessfully() {
	assert.Nil(f.t, f.err)
}

func (f *Fixture) HandlerShouldReturnBookingID(bookingID uuid.UUID) {
	assert.Equal(f.t, bookingID, f.result)
}

func (f *Fixture) HandlerShouldReturnError(message string) {
	assert.Equal(f.t, message, f.err.Error())
}

func (m *BookingRepositoryMock) InitializeBookingShouldBeCalled() {
	m.AssertCalled(m.fixture.t, "InitializeBooking")
}

func (m *FlightRepositoryMock) LockSeatsShouldBeCalled(flightID uuid.UUID, passengers int) {
	m.AssertCalled(m.fixture.t, "LockSeats", flightID, passengers)
}

func (m *FlightRepositoryMock) ReleaseSeatsShouldBeCalled(flightID uuid.UUID, seatLockIDs ...int) {
	m.AssertCalled(m.fixture.t, "ReleaseSeats", flightID, seatLockIDs)
}

func (m *BookingRepositoryMock) SaveBookingShouldBeCalledWith(dto contracts.BookingChanges) {
	m.AssertCalled(m.fixture.t, "SaveBooking", dto)
}

func (m *BookingRepositoryMock) SaveBookingShouldNotBeCalled() {
	m.AssertNotCalled(m.fixture.t, "SaveBooking", mock.Anything)
}
