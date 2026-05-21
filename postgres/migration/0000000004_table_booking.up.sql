create table dbo.booking (
	id uuid primary key default gen_random_uuid(),
	number_of_passengers int
)
