create table dbo.seat_lock (
	id int generated always as identity primary key,
	flight_id uuid references dbo.flight (id)
)
