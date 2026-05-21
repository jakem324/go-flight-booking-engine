create table dbo.flight (
	id uuid primary key default gen_random_uuid(),
	flight_code varchar(255),
	departure_date date,
	departure_airport_code varchar(255),
	departure_airport_terminal varchar(255),
	arrival_date date,
	arrival_airport_code varchar(255),
	arrival_airport_terminal varchar(255),
	max_available_seats int
)
