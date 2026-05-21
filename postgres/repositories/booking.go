// Package repositories contains all postgres repository implementations
package repositories

import (
	"context"
	"errors"

	"booking.engine/domain/entities"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookingRepository struct {
	db *pgxpool.Pool
}

func (bookingRepository BookingRepository) InitializeBooking(
	ctx context.Context, 
	dto entities.InitializeBookingDto,
) (uuid.UUID, error) {
	command := `
		insert into dbo.booking (number_of_passengers)
		values($1)
		returning id`
	var result string
	err := bookingRepository.db.QueryRow(ctx, command, dto.NumberOfPassengers).Scan(&result)
	createdBookingID, parseErr := uuid.Parse(result)
	if parseErr != nil {
		return uuid.Nil, parseErr
	}
	if err != nil {
		return uuid.Nil, err
	}
	return createdBookingID, nil
}

func (bookingRepository BookingRepository) ValidateBooking(
	ctx context.Context, 
	ID uuid.UUID,
) (entities.ValidateBookingResult, error) {
	var numberOfPassengers int
	err := bookingRepository.db.QueryRow(
		ctx,
		"select number_of_passengers from dbo.booking where id=$1",
		ID).Scan(&numberOfPassengers)
	if err == pgx.ErrNoRows {
		return entities.ValidateBookingResult{
			BookingExists: false,
		}, nil
	}
	if err != nil {
		return entities.ValidateBookingResult{}, err
	}
	
	return entities.ValidateBookingResult{
		BookingExists: true,
		NumberOfPassengers: numberOfPassengers,
	}, nil
}

func (bookingRepository BookingRepository) OnSeatsAllocated(
	ctx context.Context, 
	bookingID uuid.UUID,
	isInboundJourney bool,
	flightID uuid.UUID,
	seatLockIDs []int) error {
	return errors.New("not implemented")
}

func (bookingRepository BookingRepository) OnSeatsDeallocated(
	ctx context.Context, 
	bookingID uuid.UUID,
	isInboundJourney bool) {
	// Fire-and-forget; failures unimportant
}

func (bookingRepository BookingRepository)OnChangesCompleted(
	ctx context.Context, 
	changes entities.BookingChanges,
) error {
	// Committing at the end for SQL DB implementation rendered unnecessary by incremental writing
	return nil
}

