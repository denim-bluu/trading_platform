# Variables
APP_NAME=trading_platform
AGGREGATOR_DIR=aggregator
STRATEGY_DIR=strategy
DATA_DIR=aggregator/data
AGGREGATOR_IMAGE=$(APP_NAME)_aggregator
STRATEGY_IMAGE=$(APP_NAME)_strategy
AGGREGATOR_PORT=8080
STRATEGY_PORT=8081

# Go build flags
GO_BUILD_FLAGS=-o

# Docker variables
DOCKER_BUILD=docker build -t
DOCKER_RUN=docker run -p
DOCKER_RM=docker rm -f

# Define targets
.PHONY: all build_aggregator build_strategy run_aggregator run_strategy build_docker_aggregator build_docker_strategy run_docker_aggregator run_docker_strategy clean clean_docker

all: build_aggregator build_strategy

build_aggregator:
	@echo "Building Aggregator Service..."
	cd $(AGGREGATOR_DIR) && go build $(GO_BUILD_FLAGS) aggregator ./cmd/main.go

build_strategy:
	@echo "Building Strategy Service..."
	cd $(STRATEGY_DIR) && go build $(GO_BUILD_FLAGS) strategy ./cmd/main.go

run_aggregator: build_aggregator
	@echo "Running Aggregator Service..."
	./$(AGGREGATOR_DIR)/aggregator

run_strategy: build_strategy
	@echo "Running Strategy Service..."
	AGGREGATOR_URL=http://localhost:$(AGGREGATOR_PORT) ./$(STRATEGY_DIR)/strategy

build_docker_aggregator:
	@echo "Building Docker Image for Aggregator Service..."
	$(DOCKER_BUILD) $(AGGREGATOR_IMAGE) -f $(AGGREGATOR_DIR)/Dockerfile .

build_docker_strategy:
	@echo "Building Docker Image for Strategy Service..."
	$(DOCKER_BUILD) $(STRATEGY_IMAGE) -f $(STRATEGY_DIR)/Dockerfile .

run_docker_aggregator: build_docker_aggregator
	@echo "Running Aggregator Service in Docker..."
	$(DOCKER_RUN) $(AGGREGATOR_PORT):$(AGGREGATOR_PORT) $(AGGREGATOR_IMAGE)

run_docker_strategy: build_docker_strategy
	@echo "Running Strategy Service in Docker..."
	$(DOCKER_RUN) $(STRATEGY_PORT):$(STRATEGY_PORT) --env AGGREGATOR_URL=http://host.docker.internal:$(AGGREGATOR_PORT) $(STRATEGY_IMAGE)

clean:
	@echo "Cleaning up..."
	rm -f $(AGGREGATOR_DIR)/aggregator
	rm -f $(STRATEGY_DIR)/strategy

clean_docker:
	@echo "Cleaning up Docker containers and images..."
	-$(DOCKER_RM) $(AGGREGATOR_IMAGE)
	-$(DOCKER_RM) $(STRATEGY_IMAGE)
	docker rmi $(AGGREGATOR_IMAGE)
	docker rmi $(STRATEGY_IMAGE)
	rm -rf $(DATA_DIR)

# Commands for local development
local: clean all
	@echo "Running local development setup..."
	-mkdir -p $(DATA_DIR)
	-make run_aggregator &
	-make run_strategy &

# Commands for dockerized development
docker: clean clean_docker
	@echo "Running dockerized development setup..."
	-mkdir -p $(DATA_DIR)
	-make run_docker_aggregator &
	-make run_docker_strategy &
