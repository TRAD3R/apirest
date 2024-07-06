DEFAULT:
	run

PHONY: build
build: GOMAXPROCS=4 GOMEMLIMIT=4GiB go build -o apirest cmd/app/main.go

PHONY: run
run: go run cmd/app/main.go