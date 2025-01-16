# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## build/api/dev: builds the api for the local dev environment
build/api/dev:
	go build -o ./bin/api ./cmd

## run/api: run the /api application in the foreground
.PHONY: run/api
run/api:
	go run ./cmd

## run/api/background: run the /api application in the background
.PHONY: run/api/background
run/api/background:
	go run ./cmd &

## test/api: run the /api application in the background, then run the postman collection in docker and kill the api application once finished
.PHONY: test/api
test/api: db/reset build/api/dev run/api/background
	sleep 1 && docker compose up && pkill cmd

## test/api/auth: run the /api application in the background, then run the tests in the Auth folder of the postman collection in docker and kill the api application once finished
.PHONY: test/api/auth
test/api/auth: db/reset build/api/dev run/api/background
	sleep 1 && FOLDER=Auth docker compose up && pkill cmd

## db/reset: delete the db and recreate it via running the migrations
.PHONY: db/reset
db/reset: db/delete db/migrations/up

## db/delete: delete the sqlite database file
.PHONY: db/delete
db/delete:
	rm conduit.db

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up:
	@echo 'Running up migrations...'
	@migrate -path=./migrations -database=sqlite3://conduit.db up