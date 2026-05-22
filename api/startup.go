package api

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"

	"booking.engine/domain/commands"
	"booking.engine/domain/entities"
	"booking.engine/postgres/repositories"
)

type Handlers struct {
	PencilBookingHandler commands.PencilBookingHandler
}

func setup(ctx context.Context) (Handlers, *pgxpool.Pool) {
	connString := "postgresql://postgres:password@localhost:5432/booking_engine"
	err := migrateDB(connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to migrate DB: %v\n", err)
		os.Exit(1)
	}

	dbpool, err := pgxpool.New(ctx, connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	flightRepository := repositories.NewFlightRepository(dbpool)
	bookingRepository := repositories.NewBookingRepository(dbpool)

	flightFactory := entities.NewFlightFactory(flightRepository)
	bookingFactory := entities.NewBookingFactory(bookingRepository, flightFactory)

	pencilBookingHandler := commands.NewPencilBookingHandler(bookingFactory, flightFactory)
	return Handlers{
		PencilBookingHandler: pencilBookingHandler,
	}, dbpool
}

func migrateDB(connString string) error {
	db, err := sql.Open(
		"pgx",
		connString,
	)
	if err != nil {
		return err
	}

	driver, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://postgres/migration",
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
