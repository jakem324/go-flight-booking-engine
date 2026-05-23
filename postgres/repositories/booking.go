// Package repositories contains all postgres repository implementations
package repositories

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"booking.engine/domain/contracts"
)

type BookingRepository struct {
	db *pgxpool.Pool
}

func NewBookingRepository(db *pgxpool.Pool) BookingRepository {
	return BookingRepository{db: db}
}

func (bookingRepository BookingRepository) InitializeBooking(
	ctx context.Context,
	dto contracts.InitializeBookingDto,
) (uuid.UUID, error) {
	command := `
		insert into dbo.booking (number_of_passengers)
		values($1)
		returning id`
	var result string
	err := bookingRepository.db.QueryRow(ctx, command, dto.NumberOfPassengers).Scan(&result)
	if err != nil {
		return uuid.Nil, err
	}

	createdBookingID, parseErr := uuid.Parse(result)
	if parseErr != nil {
		return uuid.Nil, parseErr
	}

	return createdBookingID, nil
}

func (bookingRepository BookingRepository) ValidateBooking(
	ctx context.Context,
	ID uuid.UUID,
) (contracts.ValidateBookingResult, error) {
	var numberOfPassengers int
	err := bookingRepository.db.QueryRow(
		ctx,
		"select number_of_passengers from dbo.booking where id=$1",
		ID).Scan(&numberOfPassengers)
	if err == pgx.ErrNoRows {
		return contracts.ValidateBookingResult{
			BookingExists: false,
		}, nil
	}
	if err != nil {
		return contracts.ValidateBookingResult{}, err
	}

	return contracts.ValidateBookingResult{
		BookingExists:      true,
		NumberOfPassengers: numberOfPassengers,
	}, nil
}

func (bookingRepository BookingRepository) OnSeatsAllocated(
	ctx context.Context,
	bookingID uuid.UUID,
	isInboundJourney bool,
	flightID uuid.UUID,
	seatLockIDs []int) error {

	convertedLockIDs := make([]int32, len(seatLockIDs))
	for i, v := range seatLockIDs {
		convertedLockIDs[i] = int32(v)
	}

	allocationType := "outbound"
	if isInboundJourney {
		allocationType = "inbound"
	}

	command := `
		insert into dbo.booking_flight_allocation (booking_id, allocation_type, seat_lock_id)
		select $1, $2, unnest($3::int[]);
	`

	_, err := bookingRepository.db.Exec(ctx, command, bookingID, allocationType, convertedLockIDs)
	return err
}

func (bookingRepository BookingRepository) OnSeatsDeallocated(
	ctx context.Context,
	bookingID uuid.UUID,
	isInboundJourney bool) {
	// Fire-and-forget; failures unimportant
	allocationType := "outbound"
	if isInboundJourney {
		allocationType = "inbound"
	}
	command := "delete from dbo.booking_flight_allocation where booking_id = $1 and allocation_type = $2"
	_, err := bookingRepository.db.Exec(ctx, command, bookingID, bookingID, allocationType)
	if err != nil {
		log.Printf("Warning: failed to deallocate seats from booking %v %v journey. Err: %v", bookingID, allocationType, err)
	}
}

func (bookingRepository BookingRepository) OnChangesCompleted(
	ctx context.Context,
	changes contracts.BookingChanges,
) error {
	// Committing at the end for SQL DB implementation rendered unnecessary by incremental writing
	return nil
}
