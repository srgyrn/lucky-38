PROJECT_NAME := "lucky38"

export GO111MODULE=on

## Run server
server-run:
	@echo "Starting docker containers..."
	@docker-compose -f deployment/docker-compose.yml up --build

## Stop server by running docker-compose down
server-stop:
	@echo "Stopping containers..."
	@docker-compose -f deployment/docker-compose.yml down

## Restart server
server-restart: server-stop server-run

## Run unit tests
test:
	@echo "Running tests..."
	@docker exec -it --env APP_ENV=test lucky_api go test -v ./pkg/...