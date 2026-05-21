// Package commands houses the handlers for all commands within the domain
package commands

import "context"
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

func (handler *PencilBookingHandler) CreatePencilBooking(ctx context.Context, dto CreatePencilBookingDto) (uuid.UUID, error) {
	booking, err := handler.bookingFactory.NewBooking(ctx, dto.RequiredNumberOfSeats)
	if err != nil {
		return uuid.Nil, err
	}

	seatsUnavailable, err := handler.tryBookSeats(ctx, &booking.Outbound, dto.OutboundJourneyLegs)

	if err != nil {
		return uuid.Nil, err
	}

	if seatsUnavailable {
		return uuid.Nil, errors.New("seat(s) no longer available")
	}

	err = booking.FinalizeChanges(ctx)
	if err != nil {
		return uuid.Nil, err
	}

	return booking.ID, nil
}

type SetInboundJourneyDto struct {
	BookingID uuid.UUID
	InboundJourneyLegs []uuid.UUID
}

func (handler *PencilBookingHandler) SetInboundJourney(ctx context.Context, dto SetInboundJourneyDto) error {
	booking, err := handler.bookingFactory.ExistingBooking(ctx, dto.BookingID)
	if err != nil {
		return err
	}

	if booking == nil {
		return errors.New("booking not found")
	}

	seatsUnavailable, err := handler.tryBookSeats(ctx, &booking.Inbound, dto.InboundJourneyLegs)
	if seatsUnavailable {
		return errors.New("seat(s) no longer available")
	}

	if err != nil {
		return err
	}

	err = booking.FinalizeChanges(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (handler *PencilBookingHandler) tryBookSeats(ctx context.Context, journey *entities.Journey, proposedLegs []uuid.UUID) (bool, error) {
		for _, proposedLeg := range proposedLegs {
			flight := handler.flightFactory.NewFlight(proposedLeg)
			seatsObtained, err := flight.TryBookSeats(ctx, journey)
			if !seatsObtained || err != nil {
				// NB: The release of the already-allocated seats could fail, but nothing 
				// can be done about it within this scope. The application will make its best 
				// effort to avoid leaving orphan seat locks, but a background service will need 
				// to clean up any stale bookings and release the seats that were locked and could 
				// not be released by this sync workflow. With this in mind, this workflow treats 
				// the release as a fire-and-forget, hence we are not awaiting any result object.

				// (This workflow being able to lock a seat but unable to subsequently release the 
				// lock is a one-in-a-million edge-case)
				journey.ReleaseAllSeats(ctx)	
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

