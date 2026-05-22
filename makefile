start-local-dependencies:
	@docker rm -f booking-engine-db >/dev/null 2>&1 || true
	@docker run \
		--name booking-engine-db \
		-e POSTGRES_PASSWORD=password \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_DB=booking_engine \
		-p 5432:5432 \
		-d postgres:16 >/dev/null

stop-local-dependencies:
	docker stop booking-engine-db

PID_FILE := api.pid

.PHONY: start-api stop-api restart-api

start-api:
	@if [ -f $(PID_FILE) ]; then \
		echo "API is already running."; \
	else \
		go run . > api.log 2>&1 & \
		PID=$$!; \
		sleep 1; \
		pgrep -P $$PID > $(PID_FILE) || echo $$PID > $(PID_FILE); \
		echo "API started in background."; \
	fi

stop-api:
	@if [ -f $(PID_FILE) ]; then \
		PID=$$(cat $(PID_FILE)); \
		echo "Stopping API..."; \
		pkill -P $$PID 2>/dev/null || true; \
		kill $$PID 2>/dev/null || true; \
		rm -f $(PID_FILE); \
		echo "API stopped."; \
	else \
		echo "API is not running."; \
	fi

lint:
	golangci-lint run --fix

format:
	golangci-lint fmt
