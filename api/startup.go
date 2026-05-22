package api

import "context"
import "fmt"
import "os"
import "booking.engine/domain/entities"
import "booking.engine/domain/commands"
import "booking.engine/postgres/repositories"
import "github.com/jackc/pgx/v5/pgxpool"

type Handlers struct {
	PencilBookingHandler commands.PencilBookingHandler
}

func setup(ctx context.Context) Handlers {
	connString := "postgresql://postgres:password@localhost:5432/postgres"
	dbpool, err := pgxpool.New(ctx, connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	flightRepository := repositories.NewFlightRepository(dbpool)
	bookingRepository := repositories.NewBookingRepository(dbpool)

	flightFactory := entities.NewFlightFactory(flightRepository)
	bookingFactory := entities.NewBookingFactory(bookingRepository, flightFactory)

	pencilBookingHandler := commands.NewPencilBookingHandler(bookingFactory, flightFactory)
	return Handlers{
		PencilBookingHandler: pencilBookingHandler,
	}
}

