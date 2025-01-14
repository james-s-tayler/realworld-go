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

## run/api: run the /api application in the foreground
.PHONY: run/api
run/api:
	go run ./cmd

## run/api/background: run the /api application in the background
.PHONY: run/api/background
run/api/background:
	go run ./cmd &

## test/api: run the /api application in the background, then run the postman collection in docker and kill the api application once finished
test/api: run/api/background
	sleep 1 && docker compose up && pkill cmd

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	@migrate -path=./migrations -database=sqlite3://conduit.db up