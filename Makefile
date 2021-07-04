format:
	@go mod tidy
	@go fmt github.com/bensooraj/rndmicu/...

lint:
	golangci-lint run

models:
	sqlc generate --file data/sqlc.json

utils-create: utils-dropall
	go run data/utils/utils.go create

utils-dropall:
	go run data/utils/utils.go dropall

utils-seed: utils-create
	go run data/utils/utils.go seed

graph-gen:
	go run github.com/99designs/gqlgen generate

dc-up:
	docker-compose up --build

clean:
	docker system prune --force

run:
	@go run server.go