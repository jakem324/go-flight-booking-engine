package commands

import (
	"fmt"
	"testing"

	"github.com/google/uuid"

	"booking.engine/domain/contracts"
)

func TestCreatePencilBooking_AllSeatsAvailable(t *testing.T) {
	// Arrange
	fixture := CreateFixture(t)

	bookingID := uuid.New()
	passengers := 5
	firstFlightID := uuid.New()
	secondFlightID := uuid.New()

	expectedChangesDto := contracts.BookingChanges{
		ID:                 bookingID,
		NumberOfPassengers: passengers,
		OutboundLegs: []contracts.JourneyLeg{
			{
				FlightID:    firstFlightID,
				SeatLockIDs: []int{472, 673, 839},
			},
			{
				FlightID:    secondFlightID,
				SeatLockIDs: []int{293, 572, 904},
			},
		},
	}

	// Given
	fixture.bookingRepositoryMock.GivenInitializeBookingMock(bookingID)
	fixture.bookingRepositoryMock.GivenValidateBookingMock(bookingID, passengers)
	fixture.flightRepositoryMock.GivenLockSeatsMock(firstFlightID, passengers, 472, 673, 839)
	fixture.flightRepositoryMock.GivenLockSeatsMock(secondFlightID, passengers, 293, 572, 904)
	fixture.bookingRepositoryMock.GivenSaveBookingMock()

	// When
	fixture.WhenCreatePencilBookingIsCalled(passengers, firstFlightID, secondFlightID)

	// Then
	fixture.HandlerShouldCompleteSuccessfully()
	fixture.HandlerShouldReturnBookingID(bookingID)
	fixture.bookingRepositoryMock.InitializeBookingShouldBeCalled()
	fixture.flightRepositoryMock.LockSeatsShouldBeCalled(firstFlightID, passengers)
	fixture.flightRepositoryMock.LockSeatsShouldBeCalled(secondFlightID, passengers)
	fixture.bookingRepositoryMock.SaveBookingShouldBeCalledWith(expectedChangesDto)
}

func TestCreatePencilBooking_PartialUnavailable(t *testing.T) {
	// Arrange
	fixture := CreateFixture(t)

	bookingID := uuid.New()
	passengers := 5
	firstFlightID := uuid.New()
	secondFlightID := uuid.New()
	thirdFlightID := uuid.New()

	// Given
	fixture.bookingRepositoryMock.GivenInitializeBookingMock(bookingID)
	fixture.bookingRepositoryMock.GivenValidateBookingMock(bookingID, passengers)
	fixture.flightRepositoryMock.GivenLockSeatsMock(firstFlightID, passengers, 472, 673, 839)
	fixture.flightRepositoryMock.GivenLockSeatsMock(secondFlightID, passengers, 582, 612, 783)
	fixture.flightRepositoryMock.GivenLockSeatsMockWithUnavailableResult(thirdFlightID, passengers)
	// ^^ Will cause "seat(s) no longer available" result (returned error is nil yet no seat locks were returned)
	fixture.flightRepositoryMock.GivenReleaseSeatsMock()
	fixture.bookingRepositoryMock.GivenSaveBookingMock()

	// When
	fixture.WhenCreatePencilBookingIsCalled(passengers, firstFlightID, secondFlightID, thirdFlightID)

	// Then
	fixture.bookingRepositoryMock.InitializeBookingShouldBeCalled()
	fixture.flightRepositoryMock.LockSeatsShouldBeCalled(firstFlightID, passengers)
	fixture.flightRepositoryMock.LockSeatsShouldBeCalled(secondFlightID, passengers)
	fixture.flightRepositoryMock.LockSeatsShouldBeCalled(thirdFlightID, passengers) // <-- availability failure here
	// Successfully-locked seats from previous two flights now released
	fixture.flightRepositoryMock.ReleaseSeatsShouldBeCalled(firstFlightID, 472, 673, 839)
	fixture.flightRepositoryMock.ReleaseSeatsShouldBeCalled(secondFlightID, 582, 612, 783)
	// "seat(s) no longer available" returned as error
	fixture.HandlerShouldReturnError("seat(s) no longer available")
	fixture.bookingRepositoryMock.SaveBookingShouldNotBeCalled()
}

func TestSetInboundJourney_AllSeatsAvailable(t *testing.T) {
	// Arrange
	fixture := CreateFixture(t)

	bookingID := uuid.New()
	passengers := 5
	firstFlightID := uuid.New()
	secondFlightID := uuid.New()
	expectedChangesDto := contracts.BookingChanges{
		ID:                 bookingID,
		NumberOfPassengers: passengers,
		InboundLegs: []contracts.JourneyLeg{
			{
				FlightID:    firstFlightID,
				SeatLockIDs: []int{472, 673, 839},
			},
			{
				FlightID:    secondFlightID,
				SeatLockIDs: []int{293, 572, 904},
			},
		},
	}

	// Given
	fixture.bookingRepositoryMock.GivenValidateBookingMock(bookingID, passengers)
	fixture.flightRepositoryMock.GivenLockSeatsMock(firstFlightID, passengers, 472, 673, 839)
	fixture.flightRepositoryMock.GivenLockSeatsMock(secondFlightID, passengers, 293, 572, 904)
	fixture.bookingRepositoryMock.GivenSaveBookingMock()

	// When
	fixture.WhenSetInboundJourneyIsCalled(bookingID, firstFlightID, secondFlightID)

	// Then
	fixture.HandlerShouldCompleteSuccessfully()
	fixture.flightRepositoryMock.LockSeatsShouldBeCalled(firstFlightID, passengers)
	fixture.flightRepositoryMock.LockSeatsShouldBeCalled(secondFlightID, passengers)
	fixture.bookingRepositoryMock.SaveBookingShouldBeCalledWith(expectedChangesDto)
}

func TestSetInboundJourney_InvalidBookingId(t *testing.T) {
	// Arrange
	fixture := CreateFixture(t)
	bookingID := uuid.New()
	firstFlightID := uuid.New()
	secondFlightID := uuid.New()

	fixture.bookingRepositoryMock.GivenValidateBookingMockWithNegativeResult(bookingID)
	fixture.WhenSetInboundJourneyIsCalled(bookingID, firstFlightID, secondFlightID)
	fixture.HandlerShouldReturnError(fmt.Sprintf("booking not found: %v", bookingID))
}

func TestSetInboundJourney_PartialUnavailable(t *testing.T) {
	// Arrange
	fixture := CreateFixture(t)

	bookingID := uuid.New()
	passengers := 5
	firstFlightID := uuid.New()
	secondFlightID := uuid.New()
	thirdFlightID := uuid.New()

	// Given
	fixture.bookingRepositoryMock.GivenValidateBookingMock(bookingID, passengers)
	fixture.flightRepositoryMock.GivenLockSeatsMock(firstFlightID, passengers, 472, 673, 839)
	fixture.flightRepositoryMock.GivenLockSeatsMock(secondFlightID, passengers, 582, 612, 783)
	fixture.flightRepositoryMock.GivenLockSeatsMockWithUnavailableResult(thirdFlightID, passengers)
	// ^ Will cause "seat(s) no longer available" result (error is nil yet no seat locks were returned)
	fixture.flightRepositoryMock.GivenReleaseSeatsMock()

	// When
	fixture.WhenSetInboundJourneyIsCalled(bookingID, firstFlightID, secondFlightID, thirdFlightID)

	// Then
	fixture.HandlerShouldReturnError("seat(s) no longer available")
	fixture.flightRepositoryMock.LockSeatsShouldBeCalled(firstFlightID, passengers)
	fixture.flightRepositoryMock.LockSeatsShouldBeCalled(secondFlightID, passengers)
	fixture.flightRepositoryMock.LockSeatsShouldBeCalled(thirdFlightID, passengers) // <-- availability failure here
	// Successfully-locked seats from previous two flights now released
	fixture.flightRepositoryMock.ReleaseSeatsShouldBeCalled(firstFlightID, 472, 673, 839)
	fixture.flightRepositoryMock.ReleaseSeatsShouldBeCalled(secondFlightID, 582, 612, 783)
}

