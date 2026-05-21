// Package repositories contains all postgres repository implementations
package repositories

import "context"
import "errors"
import "github.com/google/uuid"
import "booking.engine/domain/entities"
import "github.com/jackc/pgx/v5/pgxpool"

type BookingRepository struct {
	db *pgxpool.Pool
}

func (bookingRepository BookingRepository) InitializeBooking(
	ctx context.Context, 
	dto entities.InitializeBookingDto,
) (uuid.UUID, error) {
	/*
	command := `
		insert into dbo.booking (number_of_passengers)
		values($1)
		returning id`
	var bookingID int
	err := db.QueryRow(ctx, query, "alice", "alice@example.com", 30).Scan(&bookingID)
	*/
	return uuid.Nil, errors.New("not implemented")
}

func (bookingRepository BookingRepository) ValidateBooking(
	ctx context.Context, 
	ID uuid.UUID,
) (entities.ValidateBookingResult, error) {
	return entities.ValidateBookingResult{}, errors.New("not implemented")
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

