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

type CreatePencilBookingResultDto struct {
	BookingId int
	Error error
}

func (service PencilService) CreatePencilBooking(dto CreatePencilBookingDto) chan CreatePencilBookingResultDto {
	out := make(chan CreatePencilBookingResultDto, 1)

	go func() {
		createBookingResult := <-service.bookingRepository.CreateBooking()	
		if createBookingResult.Error != nil {
			out <- CreatePencilBookingResultDto { 0, createBookingResult.Error }
			return
		}

		booking := &createBookingResult.Booking
		seatBookingResult := <- service.tryBookSeats(&booking.Outbound, dto.OutboundJourneyLegs, dto.RequiredNumberOfSeats)

		//	
		


		//
	}()

	return out
}

type bookSeatsResult struct {
	RequestedSeatsAvailable bool
	Error error
}

func (service PencilService) tryBookSeats(journey *entities.Journey, proposedLegs []uuid.UUID, requiredSeats int) chan bookSeatsResult {
	out := make(chan bookSeatsResult, 1)
	go func(){
		result := bookSeatsResult{
			RequestedSeatsAvailable: false,
			Error: nil,
		}

		for _, proposedLeg := range proposedLegs {
			flight := entities.NewFlight(proposedLeg)
			allocationResult := <- flight.TryAllocateSeat(journey)
			if !allocationResult.Available || allocationResult.Error != nil {
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
			
			if allocationResult.Error != nil {
				result.Error = allocationResult.Error
				out <- result
				return
			}

			if !allocationResult.Available {
				out <- result
				return
			}
		}

		result.RequestedSeatsAvailable = true
		out <- result
	}()

	return out
}



