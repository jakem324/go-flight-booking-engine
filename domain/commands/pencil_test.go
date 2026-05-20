package commands

import (
	"testing"
	"booking.engine/domain/entities"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

/*
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
	// Arrange
	fixture := CreateFixture()

	bookingID := uuid.New()
	passengers := 5
	firstFlightID := uuid.New()
	secondFlightID := uuid.New()

	dto := CreatePencilBookingDto{
		RequiredNumberOfSeats: passengers,
		OutboundJourneyLegs: []uuid.UUID { firstFlightID, secondFlightID },	
	}

	fixture.bookingRepositoryMock.On("InitializeBooking", mock.Anything).Return(bookingID, nil)
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

	// Act
	result, err := fixture.handler.CreatePencilBooking(dto)

	// Assert
	expectedInitializationDto := entities.InitializeBookingDto{
		NumberOfPassengers: passengers,
	}

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

/*
	TestCreatePencilBooking_PartialUnavailable
	Scenario: Create Pencil Booking when seats are unavailable on one or more specified flights
	All assertions:
	 - Booking initialized via booking repo
	 - Requested seats locked via flight repo
	 - Locked seats allocated via booking repo
	 - Seats subsequently deallocated when unavailability discovered on subsequent flight
	 - Deallocated seats relleased via flight repo
	 - Changes NOT finalized via booking repo
	 - Error returned stating seat(s) no longer available
*/
func TestCreatePencilBooking_PartialUnavailable(t* testing.T) {
	// Arrange
	fixture := CreateFixture()

	bookingID := uuid.New()
	passengers := 5
	firstFlightID := uuid.New()
	secondFlightID := uuid.New()
	thirdFlightID := uuid.New()

	dto := CreatePencilBookingDto{
		RequiredNumberOfSeats: passengers,
		OutboundJourneyLegs: []uuid.UUID { firstFlightID, secondFlightID, thirdFlightID },	
	}

	fixture.bookingRepositoryMock.On("InitializeBooking", mock.Anything).Return(bookingID, nil)
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
		).Return([]int {582, 612, 783}, nil)
	fixture.flightRepositoryMock.On(
		"LockSeats",
		thirdFlightID,
		passengers,
		).Return(nil, nil) // <-- Will cause "seat(s) no longer available" result (error is nil yet no seat locks were returned)
	fixture.bookingRepositoryMock.On("OnSeatsAllocated", bookingID, false, mock.Anything, mock.Anything).Return(nil)
	fixture.flightRepositoryMock.On("ReleaseSeats", mock.Anything, mock.Anything).Return(nil)
	fixture.bookingRepositoryMock.On("OnSeatsDeallocated", bookingID, false).Return(nil)

	// Act
	_, err := fixture.handler.CreatePencilBooking(dto)

	// Assert
	expectedInitializationDto := entities.InitializeBookingDto{
		NumberOfPassengers: passengers,
	}

	fixture.bookingRepositoryMock.AssertCalled(t, "InitializeBooking", expectedInitializationDto)
	fixture.flightRepositoryMock.AssertCalled(t, "LockSeats", firstFlightID, passengers)
	fixture.bookingRepositoryMock.AssertCalled(t, "OnSeatsAllocated", bookingID, false, firstFlightID, []int {472, 673, 839})
	fixture.flightRepositoryMock.AssertCalled(t, "LockSeats", secondFlightID, passengers)
	fixture.bookingRepositoryMock.AssertCalled(t, "OnSeatsAllocated", bookingID, false, secondFlightID, []int {582, 612, 783})
	fixture.flightRepositoryMock.AssertCalled(t, "LockSeats", thirdFlightID, passengers) // <-- availability failure here
	// Seats deallocated for entire journey (including previous two flights)
	fixture.bookingRepositoryMock.AssertCalled(t, "OnSeatsDeallocated", bookingID, false)
	// Successfully-locked seats from previous two flights now released
	fixture.flightRepositoryMock.AssertCalled(t, "ReleaseSeats", firstFlightID, []int{472, 673, 839})
	fixture.flightRepositoryMock.AssertCalled(t, "ReleaseSeats", secondFlightID, []int{582, 612, 783})
	// "seat(s) no longer available" returned as error
	assert.Equal(t, "Seat(s) no longer available", err.Error())
}

/*
	 Scenario: Set inbound journey with available seats
	 Assertions:
	 - Requested (available) seats locked via flight repo
	 - Locked seats allocated via booking repo
	 - Changes finalized via booking repo
	 - Nil error returned
*/

func TestSetInboundJourney_AllSeatsAvailable(t* testing.T) {
	// Arrange
	fixture := CreateFixture()

	bookingID := uuid.New()
	passengers := 5
	firstFlightID := uuid.New()
	secondFlightID := uuid.New()

	dto := SetInboundJourneyDto{
		BookingID: bookingID,
		InboundJourneyLegs: []uuid.UUID { firstFlightID, secondFlightID },	
	}

	fixture.bookingRepositoryMock.On("ValidateBooking", bookingID).Return(
		entities.ValidateBookingResult{
			BookingExists: true,
			NumberOfPassengers: passengers,
		}, nil)
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
	fixture.bookingRepositoryMock.On("OnSeatsAllocated", bookingID, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	fixture.bookingRepositoryMock.On("OnChangesCompleted", mock.Anything).Return(nil)

	// Act
	err := fixture.handler.SetInboundJourney(dto)

	// Assert
	expectedChangesDto := entities.BookingChanges{
		ID: bookingID,
		NumberOfPassengers: passengers,
		InboundLegs: []entities.JourneyLeg {
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
		fixture.flightRepositoryMock.AssertCalled(t, "LockSeats", firstFlightID, passengers)
		fixture.flightRepositoryMock.AssertCalled(t, "LockSeats", secondFlightID, passengers)
		fixture.bookingRepositoryMock.AssertCalled(t, "OnSeatsAllocated", bookingID, true, firstFlightID, []int {472, 673, 839})
		fixture.bookingRepositoryMock.AssertCalled(t, "OnSeatsAllocated", bookingID, true, secondFlightID, []int {293, 572, 904})
		fixture.bookingRepositoryMock.AssertCalled(t, "OnChangesCompleted", expectedChangesDto)
	}
}

/*
	 Scenario: Set inbound journey with invalid booking ID
	 Assertions:
	 - Invalid booking ID error returned
*/

func TestSetInboundJourney_InvalidBookingId(t* testing.T) {
	// Arrange
	fixture := CreateFixture()
	fixture.bookingRepositoryMock.On("ValidateBooking", mock.Anything).Return(
		entities.ValidateBookingResult{
			BookingExists: false,
			NumberOfPassengers: 0,
		}, nil)

	bookingID := uuid.New()
	firstFlightID := uuid.New()
	secondFlightID := uuid.New()
	dto := SetInboundJourneyDto{
		BookingID: bookingID,
		InboundJourneyLegs: []uuid.UUID { firstFlightID, secondFlightID },	
	}

	// Act
	err := fixture.handler.SetInboundJourney(dto)

	// Assert
	assert.Equal(t, "booking not found", err.Error())
}

