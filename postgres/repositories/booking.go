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
) (uuid.UUID, error) {
	/*
	command := `
		insert into dbo.booking (number_of_passengers)
		values($1)
		returning id`
	var result string
	err := bookingRepository.db.QueryRow(ctx, command, dto.NumberOfPassengers).Scan(&result)
	*/
	command := `
		insert into dbo.booking
		default values
		returning id`
	var result string
	err := bookingRepository.db.QueryRow(ctx, command).Scan(&result)
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

func (bookingRepository BookingRepository) SaveBooking(
	ctx context.Context,
	changes contracts.BookingChanges,
) error {
	err := bookingRepository.writeBookingDetails(ctx, changes)
	if err != nil {
		return err
	}

	var outboundJourneySeats []int32
	for _, leg := range changes.OutboundLegs {
		for _, seat := range leg.SeatLockIDs {
			outboundJourneySeats = append(outboundJourneySeats, int32(seat))
		}
	}

	var inboundJourneySeats []int32
	for _, leg := range changes.InboundLegs {
		for _, seat := range leg.SeatLockIDs {
			inboundJourneySeats = append(inboundJourneySeats, int32(seat))
		}
	}

	err = bookingRepository.allocateSeats(
		ctx,
		changes.ID,
		false,
		outboundJourneySeats,
	)
	if err != nil {
		return err
	}

	err = bookingRepository.allocateSeats(
		ctx,
		changes.ID,
		true,
		inboundJourneySeats,
	)
	if err != nil {
		return err
	}

	return nil
}

func (bookingRepository BookingRepository) writeBookingDetails(
	ctx context.Context,
	changes contracts.BookingChanges,
) error {
	command := `
		update dbo.booking
		set number_of_passengers = $1
		where id=$2`
	_, err := bookingRepository.db.Exec(ctx, command, changes.ID, changes.NumberOfPassengers)
	return err
}

func (bookingRepository BookingRepository) allocateSeats(
	ctx context.Context,
	bookingID uuid.UUID,
	isInboundJourney bool,
	seatLockIDs []int32) error {

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

func (bookingRepository BookingRepository) deallocateSeats(
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
