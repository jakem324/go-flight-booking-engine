package services

import "github.com/google/uuid"
import "booking.engine/domain/entities"

type PencilService struct {
	bookingRepository entities.BookingRepository
}

type CreatePencilBookingDto struct {
	RequiredNumberOfSeats int
	OutboundJourneyLegs []uuid.UUID
}

func (service PencilService) CreatePencilBooking(dto CreatePencilBookingDto) (uuid.UUID, error) {
	booking, err := service.bookingRepository.CreateBooking()	
	if err != nil {
		return uuid.Nil, err
	}

	seatsObtained, err := service.tryBookSeats(&booking.Outbound, dto.OutboundJourneyLegs, dto.RequiredNumberOfSeats)

}

func (service PencilService) tryBookSeats(journey *entities.Journey, proposedLegs []uuid.UUID, requiredSeats int) (bool, error) {
		for _, proposedLeg := range proposedLegs {
			flight := entities.NewFlight(proposedLeg)
			seatObtained, err := flight.TryAllocateSeat(journey)
			if !seatObtained || err != nil {
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

			if !seatObtained {
				return false, nil
			}
		}

		return true, nil
}

