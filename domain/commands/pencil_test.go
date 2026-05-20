package commands

import (
	"testing"
	"booking.engine/domain/entities"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

/*
	TestCreatePencilBooking_AllSeatsAvailable
	Scenario: Create Pencil Booking when seats are available on specified flights
	All happy-path assertions:
	 - Booking initialized via booking repo
	 - Requested (available) seats locked via flight repo
	 - Locked seats allocated via booking repo
	 - Changes finalized via booking repo
	 - Created booking ID returned

	TestCreatePencilBooking_Partial
	Scenario: Create Pencil Booking when seats are unavailable on one or more specified flights
	All assertions:
	 - Booking initialized via booking repo
	 - Requested seats locked via flight repo
	 - Locked seats allocated via booking repo
	 - Seats subsequently deallocated when unavailability discovered on subsequent flight
	 - Deallocated seats relleased via flight repo
	 - Changes NOT finalized via booking repo
	 - Error returned stating seat(s) no longer available

	 Scenario: Set inbound journey with available seats
	 Assertions:
	 - Requested (available) seats locked via flight repo
	 - Locked seats allocated via booking repo
	 - Changes finalized via booking repo
	 - Nil error returned

	 Scenario: Set inbound journey with invalid booking ID
	 Assertions:
	 - Invalid booking ID error returned

	 Scenario: Set inbound journey when sets are unavailable on one or more flights
	 Assertions:
	 - Requested seats locked via flight repo
	 - Locked seats allocated via booking repo
	 - Seats subsequently deallocated when unavailability discovered on subsequent flight
	 - Deallocated seats relleased via flight repo
	 - Changes NOT finalized via booking repo
	 - Error returned stating seat(s) no longer available

*/

/*
	TestCreatePencilBooking_AllSeatsAvailable
	Scenario: Create Pencil Booking when seats are available on specified flights
	All happy-path assertions:
	 - Booking initialized via booking repo
	 - Requested (available) seats locked via flight repo
	 - Locked seats allocated via booking repo
	 - Changes finalized via booking repo
	 - Created booking ID returned
*/

func TestCreatePencilBooking_AllSeatsAvailable(t* testing.T) {
	fixture := CreateFixture()

	bookingID := uuid.New()
	passengers := 5
	firstFlightID := uuid.New()
	secondFlightID := uuid.New()

	dto := CreatePencilBookingDto{
		RequiredNumberOfSeats: passengers,
		OutboundJourneyLegs: []uuid.UUID { firstFlightID, secondFlightID },	
	}

	expectedInitializationDto := entities.InitializeBookingDto{
		NumberOfPassengers: passengers,
	}

	fixture.bookingRepositoryMock.On("InitializeBooking", expectedInitializationDto).Return(bookingID, nil)
	fixture.bookingRepositoryMock.On("ValidateBooking", bookingID).Return(
		entities.ValidateBookingResult{NumberOfPassengers: 5}, nil)
	fixture.flightRepositoryMock.On(
		"LockSeats",
		firstFlightID,
		passengers,
		).Return([]int {472, 673, 839}, nil)
	fixture.flightRepositoryMock.On(
		"LockSeats",
		secondFlightID,
		passengers,
		).Return([]int {293, 572, 904}, nil)
	fixture.bookingRepositoryMock.On("OnSeatsAllocated", bookingID, false, mock.Anything, mock.Anything).Return(nil)
	fixture.bookingRepositoryMock.On("OnChangesCompleted", mock.Anything).Return(nil)

	result, err := fixture.handler.CreatePencilBooking(dto)

	expectedChangesDto := entities.BookingChanges{
		ID: bookingID,
		NumberOfPassengers: passengers,
		OutboundLegs: []entities.JourneyLeg {
			{
				FlightID: firstFlightID,
				SeatLockIDs: []int {472, 673, 839},
			},
			{
				FlightID: secondFlightID,
				SeatLockIDs: []int {293, 572, 904},
			},
		},
	}

	if assert.Nil(t, err) {
		fixture.bookingRepositoryMock.AssertCalled(t, "InitializeBooking", expectedInitializationDto)
		fixture.flightRepositoryMock.AssertCalled(t, "LockSeats", firstFlightID, passengers)
		fixture.flightRepositoryMock.AssertCalled(t, "LockSeats", secondFlightID, passengers)
		fixture.bookingRepositoryMock.AssertCalled(t, "OnSeatsAllocated", bookingID, false, firstFlightID, []int {472, 673, 839})
		fixture.bookingRepositoryMock.AssertCalled(t, "OnSeatsAllocated", bookingID, false, secondFlightID, []int {293, 572, 904})
		fixture.bookingRepositoryMock.AssertCalled(t, "OnChangesCompleted", expectedChangesDto)
		assert.Equal(t, bookingID, result)
	}
}

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

func (m *BookingRepositoryMock) OnSeatsAllocated(
	bookingID uuid.UUID,
	isInboundJourney bool,
	flightID uuid.UUID,
	seatLockIDs []int) error {

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





/*
func TestCreatePencilBooking_CommandWithZeroRequiredSeatsIsRejected(t* testing.T) {
	fixture := CreateFixture()
	dto := CreatePencilBookingDto{ RequiredNumberOfSeats: 0 }
	_, err := fixture.handler.CreatePencilBooking(dto)

	if assert.NotNil(t, err) {
		assert.Equal(t, "invalid number of passengers", err.Error())
	}

	fixture.bookingRepositoryMock.AssertNotCalled(t, "InitializeBooking")
	fixture.bookingRepositoryMock.AssertNotCalled(t, "OnChangesCompleted")
}

func TestCreatePencilBooking_BookingIsInitialized(t* testing.T) {
	fixture := CreateFixture()

	bookingID := uuid.New()
	expectedInitializationDto := entities.InitializeBookingDto{
		NumberOfPassengers: 5,
	}

	fixture.bookingRepositoryMock.On("InitializeBooking", expectedInitializationDto).Return(bookingID, nil)
	fixture.bookingRepositoryMock.On("ValidateBooking", bookingID).Return(
		entities.ValidateBookingResult{NumberOfPassengers: 5}, nil)
	fixture.flightRepositoryMock.On(
		"LockSeats",
		mock.AnythingOfType("uuid.UUID"),
		mock.AnythingOfType("int")).Return([]int {1,2,3}, nil)
	fixture.bookingRepositoryMock.On("OnSeatsAllocated", bookingID, false).Return(nil)
	fixture.bookingRepositoryMock.On("OnChangesCompleted", mock.Anything).Return(nil)

	dto := CreatePencilBookingDto{ RequiredNumberOfSeats: 5 }
	_, err := fixture.handler.CreatePencilBooking(dto)

	if assert.Nil(t, err) {
		fixture.bookingRepositoryMock.AssertCalled(t, "InitializeBooking", expectedInitializationDto)
	}
}

func TestCreatePencilBooking_PartialUnavailable(t* testing.T) {
	fixture := CreateFixture()

	bookingID := uuid.New()
	passengers := 5
	firstFlightID := uuid.New()
	secondFlightID := uuid.New()

	dto := CreatePencilBookingDto{
		RequiredNumberOfSeats: passengers,
		OutboundJourneyLegs: []uuid.UUID { firstFlightID, secondFlightID },	
	}

	expectedInitializationDto := entities.InitializeBookingDto{
		NumberOfPassengers: passengers,
	}

	fixture.bookingRepositoryMock.On("InitializeBooking", expectedInitializationDto).Return(bookingID, nil)
	fixture.bookingRepositoryMock.On("ValidateBooking", bookingID).Return(
		entities.ValidateBookingResult{NumberOfPassengers: 5}, nil)
	fixture.flightRepositoryMock.On(
		"LockSeats",
		firstFlightID,
		).Return([]int {472, 673, 839}, nil)
	fixture.flightRepositoryMock.On(
		"LockSeats",
		secondFlightID,
		).Return(nil, nil)
	// ^ will produce "seat(s) no longer available" due to nil error
	fixture.bookingRepositoryMock.On("OnSeatsAllocated", bookingID, false).Return(nil)
	fixture.bookingRepositoryMock.On("OnChangesCompleted", mock.Anything).Return(nil)

	_, err := fixture.handler.CreatePencilBooking(dto)

	if assert.Nil(t, err) {
		fixture.bookingRepositoryMock.AssertCalled(t, "InitializeBooking", expectedInitializationDto)
		fixture.flightRepositoryMock.AssertCalled(t, "LockSeats", firstFlightID, passengers)
		fixture.flightRepositoryMock.AssertCalled(t, "LockSeats", secondFlightID, passengers)
		fixture.flightRepositoryMock.AssertCalled(t, "ReleaseSeats", firstFlightID, []int{472, 673, 839})
	//LockSeats(flightID uuid.UUID, numberOfSeats int) ([]int, error)
	}
}
*/
