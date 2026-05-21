package repositories

import "errors"
import "github.com/google/uuid"
import "github.com/jackc/pgx/v5/pgxpool"

type FlightRepository struct {
	db *pgxpool.Pool
}

func (flightRepository FlightRepository) LockSeats(flightID uuid.UUID, numberOfSeats int) ([]int, error) {
	return []int{}, errors.New("not implemented")
}

func (flightRepository FlightRepository) ReleaseSeats(flightID uuid.UUID, seatLockIDs []int) {
	// Fire-and-forget; failures unimportant
}

