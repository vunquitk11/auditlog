# Makefile for audit-log-service

# Variables
PROJECT_NAME := auditlog
DOCKER_BIN := docker
DOCKER_COMPOSE_BIN := docker-compose

COMPOSE := PROJECT_NAME=${PROJECT_NAME} ${DOCKER_COMPOSE_BIN} -f docker-compose.yaml
API_COMPOSE = ${COMPOSE} run --name ${PROJECT_NAME}-api-$${CONTAINER_SUFFIX:-local} --rm -w /app app

# Base Targets
setup: migrate mocks up
teardown:
	${COMPOSE} down -v --remove-orphans
	@echo "WARNING: All data has been removed including volumes"

# Stop containers but preserve data (default behavior)
stop:
	${COMPOSE} down
	@echo "Containers stopped. Data is preserved in volumes"

# Stop containers and remove volumes (data will be lost)
stop-clean:
	${COMPOSE} down -v --remove-orphans
	@echo "Containers stopped and all data removed"

up:
	${COMPOSE} up -d --build

# Build & Run
build:
	${API_COMPOSE} sh -c "go build -mod=vendor -o bin/auditlog cmd/serverd/main.go"

# Run main application only (without producers)
run-app:
	${API_COMPOSE} sh -c "go run -mod=vendor cmd/serverd/main.go"

# Legacy run command (kept for backward compatibility)
run: run-app

# Kafka Producers - Separate from main app
build-producers:
	${API_COMPOSE} sh -c "go build -mod=vendor -o bin/user-producer cmd/producer-user/main.go"
	${API_COMPOSE} sh -c "go build -mod=vendor -o bin/payment-producer cmd/producer-payment/main.go"

# Run individual producers
run-user-producer:
	${COMPOSE} run --name ${PROJECT_NAME}-user-producer --rm -w /app app sh -c "go run -mod=vendor cmd/producer-user/main.go"

run-payment-producer:
	${COMPOSE} run --name ${PROJECT_NAME}-payment-producer --rm -w /app app sh -c "go run -mod=vendor cmd/producer-payment/main.go"

# Run both producers in separate containers
run-producers:
	${COMPOSE} run --name ${PROJECT_NAME}-user-producer --rm -w /app app sh -c "go run -mod=vendor cmd/producer-user/main.go" &
	${COMPOSE} run --name ${PROJECT_NAME}-payment-producer --rm -w /app app sh -c "go run -mod=vendor cmd/producer-payment/main.go" &
	@echo "Both producers started. Use 'make stop-producers' to stop them."

stop-producers:
	@docker stop ${PROJECT_NAME}-user-producer 2>/dev/null || true
	@docker stop ${PROJECT_NAME}-payment-producer 2>/dev/null || true
	@echo "Producers stopped."

# Demo workflow commands
demo-start: up migrate
	@echo "Starting demo: Infrastructure is ready"
	@echo "Run 'make demo-app' to start the main application"
	@echo "Run 'make demo-producers' to start the producers after app is running"

demo-app: run-app

demo-producers: run-producers
	@echo "Producers started. Data will appear in the database within 1 minute."
	@echo "Check Kafka UI at http://localhost:8081 to monitor messages"

demo-stop: stop-producers
	@echo "Demo stopped. Run 'make stop' to preserve data"

# Reset database (clear all data but keep infrastructure)
reset-db: drop migrate
	@echo "Database reset completed. All audit logs have been cleared."

# Test & Mocks
test:
	${API_COMPOSE} sh -c "go test -mod=vendor -coverprofile=c.out -failfast -timeout 5m ./..."
mocks:
	${API_COMPOSE} sh -c "mockery --all --recursive --inpackage --output internal/mocks"
api-gen-mocks:
	${COMPOSE} run --name ${PROJECT_NAME}-mockery-$${CONTAINER_SUFFIX:-local} --rm -w /app --entrypoint '' app /bin/sh -c "\
		mockery --dir internal/gateway --all --recursive --inpackage && \
		mockery --dir internal/controller --all --recursive --inpackage && \
		mockery --dir internal/repository --all --recursive --inpackage"

# Automatically bring up db before migrating
migrate:
	@docker-compose up -d db
	@sleep 3
	@docker-compose exec db sh -c 'for f in /migrations/*.up.sql; do psql -U postgres -d auditlog -f "$$f"; done'

# Drop tables using all down.sql files in the migrations folder
drop:
	${COMPOSE} exec db sh -c 'for f in /migrations/*.down.sql; do psql -U $$DB_USER -d $$DB_NAME -f "$$f"; done'

# Redo by dropping then migrating again
redo: drop migrate

# Vendor Management
vendor:
	${API_COMPOSE} sh -c "go mod tidy -compat=1.17 && go mod vendor"

# Clean
clean:
	${API_COMPOSE} sh -c "go clean -mod=vendor -i -x ./..."

# Start auditlog consumer manually
run-auditlog-consumer:
	${COMPOSE} run --name auditlog-consumer --rm -w /app app sh -c "go run cmd/auditlogconsumer/main.go"

# Start/stop user-producer container
start-user-producer:
	${COMPOSE} up -d user-producer
stop-user-producer:
	${COMPOSE} stop user-producer

# Start/stop payment-producer container
start-payment-producer:
	${COMPOSE} up -d payment-producer
stop-payment-producer:
	${COMPOSE} stop payment-producer

# Start/stop all producers
start-all-producers:
	${COMPOSE} up -d user-producer payment-producer
stop-all-producers:
	${COMPOSE} stop user-producer payment-producer

# Start/stop Prometheus
start-prometheus:
	${COMPOSE} up -d prometheus
stop-prometheus:
	${COMPOSE} stop prometheus

# Start/stop Grafana
start-grafana:
	${COMPOSE} up -d grafana
stop-grafana:
	${COMPOSE} stop grafana
	
