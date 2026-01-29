BASE_STACK = docker compose -f docker-compose.yml

compose-up: ### Run docker compose
	$(BASE_STACK) up --build -d
.PHONY: compose-up

compose-down: ### Down docker compose
	$(BASE_STACK) down -v
.PHONY: compose-down

deps: ### deps tidy + verify
	go mod tidy && go mod verify
.PHONY: deps