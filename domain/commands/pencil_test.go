package commands

import (
	"fmt"
	"testing"

	"github.com/google/uuid"

	"booking.engine/domain/contracts"
)

/*
TestCreatePencilBooking_AllSeatsAvailable
Scenario: Create Pencil Booking when seats are available on specified flights
All happy-path assertions:
  - Booking initialized via booking repo
  - Requested (available) seats locked via flight repo
  - Locked seats allocated via booking repo
  - Changes finalized via booking repo - Created booking ID returned
*/
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

	// fixture.flightRepositoryMock.On("ReleaseSeats", mock.Anything, mock.Anything).Return(nil)

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

/*
 Scenario: Set inbound journey with available seats
 Assertions:
 - Requested (available) seats locked via flight repo
 - Locked seats allocated via booking repo
 - Changes finalized via booking repo
 - Nil error returned
*/

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

/*
 Scenario: Set inbound journey with invalid booking ID
 Assertions:
 - Invalid booking ID error returned
*/

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

/*
 Scenario: Set inbound journey when seats are unavailable on one or more flights
 Assertions:
 - Requested seats locked via flight repo
 - Locked seats allocated via booking repo
 - Seats subsequently deallocated when unavailability discovered on subsequent flight
 - Deallocated seats relleased via flight repo
 - Changes NOT finalized via booking repo
 - Error returned stating seat(s) no longer available
*/

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
