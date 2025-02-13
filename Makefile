# Makefile

# Project-specific variables
DOCKER_COMPOSE_FILE=deployments/docker-compose.yml
DOCKER_COMPOSE_DB_FILE=deployments/docker-compose.db.yml

# Target to bring up the Docker Compose services
up:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

# Target to bring down the Docker Compose services
down:
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

# Delete volumes and images
destroy:
	docker-compose -f $(DOCKER_COMPOSE_FILE) down -v --rmi all

# Target to build the Docker Compose services
build:
	docker-compose -f $(DOCKER_COMPOSE_FILE) build

# Target to start the Docker Compose services
start:
	docker-compose -f $(DOCKER_COMPOSE_FILE) start

# Target to stop the Docker Compose services
stop:
	docker-compose -f $(DOCKER_COMPOSE_FILE) stop

# Target to view logs from the Docker Compose services
logs:
	docker-compose -f $(DOCKER_COMPOSE_FILE) logs

# Target to restart the Docker Compose services
restart:
	docker-compose -f $(DOCKER_COMPOSE_FILE) restart

# Target to generate proto files (assuming the script is already set up)
generate-proto:
	./scripts/generate_proto.sh

# Optional: add a target to clean the generated files
clean-proto:
	rm -rf api/proto/*

# Initiate PostgreSQL database
setup-db:
	docker-compose -f $(DOCKER_COMPOSE_DB_FILE) up -d

# Destroy PostgreSQL database
destroy-db:
	docker-compose -f $(DOCKER_COMPOSE_DB_FILE) down -v

# To run the data-service
run-data-service:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up data_service -d


# Default target
.PHONY: all
all: build up generate-proto

# Phony targets to prevent conflicts with file names
.PHONY: up down build start stop logs restart generate-proto clean-proto setup-db destroy-db