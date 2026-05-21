package repositories

import "context"
import "github.com/google/uuid"
import "github.com/jackc/pgx/v5/pgxpool"
import "booking.engine/domain/entities"

type FlightRepository struct {
	db *pgxpool.Pool
}

func (flightRepository FlightRepository) LockSeats(
	ctx context.Context, 
	flightID uuid.UUID,
	numberOfSeats int,
) (entities.SeatLockResult, error) {
	var flightValid bool
	var seatsAvailable bool
	var seatLockIDs []int
	command := "select (flight_valid, seats_available, seat_lock_ids) from dbo.try_lock_seats($1, $2)"
	err := flightRepository.db.QueryRow(ctx, command, flightID, numberOfSeats).Scan(
		flightValid, seatsAvailable, seatLockIDs)

	if err != nil {
		return entities.SeatLockResult{}, err
	}

	return entities.SeatLockResult{
		ValidFlightID: flightValid,
		SeatsAvailable: seatsAvailable,
		ObtainedSeatLockIDs: seatLockIDs,
	}, nil
}

func (flightRepository FlightRepository) ReleaseSeats(
	ctx context.Context, 
	flightID uuid.UUID,
	seatLockIDs []int,
) {
	// Fire-and-forget; failures unimportant
	command := "delete from dbo.seat_lock where flight_id=$1 id = ANY($2)"
	flightRepository.db.Exec(ctx, command, flightID, seatLockIDs)
}

