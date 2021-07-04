format:
	@go fmt github.com/bensooraj/rndmicu/...

models:
	sqlc generate --file data/sqlc.json

utils-create: utils-dropall
	go run data/utils/utils.go create

utils-dropall:
	go run data/utils/utils.go dropall

utils-seed: utils-create
	go run data/utils/utils.go seed
