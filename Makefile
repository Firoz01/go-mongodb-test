# Define the command-line arguments
CONFIG_PATH=../../config.json

.PHONY: run

run:
	-@go run ./cmd/server -c ./config.json

seed-mongodb:
	go run ./cmd/seedmongodb .