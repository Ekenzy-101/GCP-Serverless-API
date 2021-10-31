include .env
export $(shell sed 's/=.*//' .env)

dev:
	@go run cmd/main.go
	