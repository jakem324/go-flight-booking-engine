start-local-dependencies:
		docker run --name booking-engine-db -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres
