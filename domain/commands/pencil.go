// Package commands houses the handlers for all commands within the domain
package commands

import "errors"
import "github.com/google/uuid"
import "booking.engine/domain/entities"

type PencilBookingHandler struct {
	bookingFactory *entities.BookingFactory
	flightFactory *entities.FlightFactory
}

func NewPencilBookingHandler (
	bookingFactory entities.BookingFactory,
	flightFactory entities.FlightFactory) PencilBookingHandler {
	return PencilBookingHandler{
		bookingFactory: &bookingFactory,
		flightFactory: &flightFactory,
	}
}

type CreatePencilBookingDto struct {
	RequiredNumberOfSeats int
	OutboundJourneyLegs []uuid.UUID
}

func (handler *PencilBookingHandler) CreatePencilBooking(dto CreatePencilBookingDto) (uuid.UUID, error) {
	booking, err := handler.bookingFactory.NewBooking(dto.RequiredNumberOfSeats)
	if err != nil {
		return uuid.Nil, err
	}

	seatsUnavailable, err := handler.tryBookSeats(&booking.Outbound, dto.OutboundJourneyLegs)

	if seatsUnavailable {
		return uuid.Nil, errors.New("Seat(s) no longer available")
	}

	if err != nil {
		return uuid.Nil, err
	}

	err = booking.FinalizeChanges()
	if err != nil {
		return uuid.Nil, err
	}

	return booking.ID, nil
}

type SetInboundJourneyDto struct {
	BookingID uuid.UUID
	InboundJourneyLegs []uuid.UUID
}

func (handler *PencilBookingHandler) SetInboundJourney(dto SetInboundJourneyDto) error {
	booking, err := handler.bookingFactory.ExistingBooking(dto.BookingID)
	if err != nil {
		return err
	}

	seatsUnavailable, err := handler.tryBookSeats(&booking.Inbound, dto.InboundJourneyLegs)

	if seatsUnavailable {
		return errors.New("Seat(s) no longer available")
	}

	if err != nil {
		return err
	}

	err = booking.FinalizeChanges()
	if err != nil {
		return err
	}

	return nil
}

func (handler *PencilBookingHandler) tryBookSeats(journey *entities.Journey, proposedLegs []uuid.UUID) (bool, error) {
		for _, proposedLeg := range proposedLegs {
			flight := handler.flightFactory.NewFlight(proposedLeg)
			seatsObtained, err := flight.TryBookSeats(journey)
			if !seatsObtained || err != nil {
				// NB: The release of the already-allocated seats could fail, but nothing 
				// can be done about it within this scope. The application will make its best 
				// effort to avoid leaving orphan seat locks, but a background service will need 
				// to clean up any stale bookings and release the seats that were locked and could 
				// not be released by this sync workflow. With this in mind, this workflow treats 
				// the release as a fire-and-forget, hence we are not awaiting any result object.

				// (This workflow being able to lock a seat but unable to subsequently release the 
				// lock is a one-in-a-million edge-case)
				journey.ReleaseAllSeats()	
			}
			
			if err != nil {
				return false, err
			}

			if !seatsObtained {
				return true, nil
			}
		}

		return false, nil
}

