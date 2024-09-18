DEFAULT:
	run

PHONE: lint
lint:
	golangci-lint run

PHONY: build
build: lint
	GOMAXPROCS=4 GOMEMLIMIT=4GiB go build -o apirest cmd/app/main.go

PHONY: run
run:
	go run cmd/app/main.go

PHONY: docker-up
docker-up: lint
	docker compose -f ./deployments/docker-compose.yml up --build -d --remove-orphans

PHONY: tests
tests:
	go test -coverprofile=coverage.out ./...

PHONY: cover
cover:
	go tool cover -html=coverage.out