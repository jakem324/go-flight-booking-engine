package repositories

import "context"
import "errors"
import "github.com/google/uuid"
import "github.com/jackc/pgx/v5/pgxpool"

type FlightRepository struct {
	db *pgxpool.Pool
}

func (flightRepository FlightRepository) LockSeats(
	ctx context.Context, 
	flightID uuid.UUID,
	numberOfSeats int,
) ([]int, error) {
	var flightValid bool
	var seatsAvailable bool
	var seatLockIDs []int
	command := "select (flight_valid, seats_available, seat_lock_ids) from dbo.try_lock_seats($1, $2)"
	err := flightRepository.db.QueryRow(ctx, command, flightID, numberOfSeats).Scan(
		flightValid, seatsAvailable, seatLockIDs)
	if err != nil {
		return []int{}, err
	}

	if !flightValid {
		return []int{}, errors.New("flight not found")
	}

	if !seatsAvailable {
		return []int{}, errors.New("seat(s) no longer available")
	}

	return seatLockIDs, nil
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

