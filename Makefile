.PHONY: default install update run

default: run

install:
	@go mod download
update:
	@go mod tidy && go get -u ./...
run:
	@go run ./main.go
