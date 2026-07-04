// Package queryhandlers contains all postgres query handler implementations
package queryhandlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"booking.engine/domain/queries"
)

type BookingQueryHandler struct {
	db *pgxpool.Pool
}

func NewBookingQueryHandler(db *pgxpool.Pool) BookingQueryHandler {
	return BookingQueryHandler{db: db}
}

func (handler BookingQueryHandler) GetBookingSummary(ctx context.Context, ID uuid.UUID) (*queries.BookingSummary, error) {
	query := `
		with unique_legs as (
			select
				b.id as booking_id,
				b.number_of_passengers,
				bfa.allocation_type,
				f.departure_date,
				f.departure_airport_code,
				f.departure_airport_terminal,
				f.arrival_date,
				f.arrival_airport_code,
				f.arrival_airport_terminal
			from dbo.booking b 
			join dbo.booking_flight_allocation bfa on b.id = bfa.booking_id 
			join dbo.seat_lock sl on bfa.seat_lock_id = sl.id 
			join dbo.flight f on sl.flight_id = f.id
			where b.id = $1
			group by b.id, f.id, bfa.allocation_type
		)
		select json_build_object(
			'BookingID', booking_id,
			'NumberOfPassengers', number_of_passengers,
			'OutboundJourneyLegs', coalesce(json_agg(json_build_object(
				'DepartureDate', to_char(departure_date, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
				'DepartureAirportCode', departure_airport_code,
				'DepartureTerminal', departure_airport_terminal,
				'ArrivalDate', to_char(arrival_date, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
				'ArrivalAirportCode', arrival_airport_code,
				'ArrivalTerminal', arrival_airport_terminal
			)) filter (where allocation_type = 'outbound'), '[]'::json),
			'InboundJourneyLegs', coalesce(json_agg(json_build_object(
				'DepartureDate', to_char(departure_date, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
				'DepartureAirportCode', departure_airport_code,
				'DepartureTerminal', departure_airport_terminal,
				'ArrivalDate', to_char(arrival_date, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
				'ArrivalAirportCode', arrival_airport_code,
				'ArrivalTerminal', arrival_airport_terminal
			)) filter (where allocation_type = 'inbound'), '[]'::json)
		)
		from unique_legs
		group by booking_id, number_of_passengers;
	`

	var jsonBytes []byte

	err := handler.db.QueryRow(ctx, query, ID).Scan(&jsonBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	var summary queries.BookingSummary
	if err := json.Unmarshal(jsonBytes, &summary); err != nil {
		return nil, fmt.Errorf("failed to unmarshal booking summary: %w", err)
	}

	return &summary, nil

}
