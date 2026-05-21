create table dbo.booking_flight_allocation (
	booking_id uuid references dbo.booking (id),
	allocation_type varchar(255),
	seat_lock_id int references dbo.seat_lock (id)
)
