.PHONY: gen test up down

PACKAGES := $(shell go list ./... | grep -v /mocks)

gen:
	go generate ./...
test:
	# run tests without mocks
	go test -tags=integration -coverprofile=coverage.out $(PACKAGES)
	go tool cover -func=coverage.out
	rm coverage.out
up:
	docker compose up -d
down:
	docker compose down
