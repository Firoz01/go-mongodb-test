# Define the path to the docker-compose.yml file
COMPOSE_FILE = cmd/seedmongodb/docker-compose.yml
COMPOSE_CMD = docker-compose -f $(COMPOSE_FILE)

# Define the path to the config file
CONFIG_PATH = ./config.json

.PHONY: run mongodb-seed mongodb-up mongodb-down

# Run the server with the configuration file
run:
	@echo "Running the Go server..."
	@go run ./cmd/server -c $(CONFIG_PATH)

# Seed MongoDB using the seedmongodb command
mongodb-seed:
	@echo "Seeding MongoDB..."
	@go run ./cmd/seedmongodb

# Start the Docker containers defined in docker-compose.yml
mongodb-up:
	@echo "Starting MongoDB containers..."
	@$(COMPOSE_CMD) up -d

# Stop and remove the Docker containers defined in docker-compose.yml
mongodb-down:
	@echo "Stopping and removing MongoDB containers..."
	@$(COMPOSE_CMD) down
