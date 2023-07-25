include .envrc

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

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	go run ./cmd/api

## run/api/nobgp: run the cmd/api application without BGP
.PHONY: run/api/nobgp
run/api/nobgp:
	@echo 'Running without BGP...'
	go run ./cmd/api -enable-bgp=false

## db/migrations/new name=<VALUE>ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC6y+7TPzfVoZXmRCFkZTC+oz3gFiGTxrLHGg/Vcn7i6Vvm5aIk8aA2Xfa+XU2zn+y3l+8VaQPXLlfpYPadpLb/0r3owreKBB43cQNyant+yOv4VKcy2c01K6Ralk5WDQ4LwO7/d/NU9Ue/6moiRe/c/Ds1xpJKDqS+TC22gclFoMuuotDeFCybXDWGaz0cKIFxVPFEkrB9538n56IbHeIp/aGj7WDx3c0Z9uYViPN+DVw0ecmLnTMSGM7K72f+DinfioyZr+8IQ8AeYljsn0JOMapAwYeeR6dcnCfMk90FPINs3kJCNNxLh4UHnO4K/TSnQBu3BkNNer0w1gJVpDDh1CCAMq/TDA1Tdxkvm8AN/FprfUdE7GXvxHHUNFEV943syJwKslYz8XUnfz/b7s9e3QqIg6G4La2h/2vL3HwZTmKShQGbeNZg4KNECzxQ+2J6+fw91+vNACmYEGmmzerMKg1Emjz5qcKfJiSWoBOpVe9LhGmmvH5lIz8+tQrhXr8= naruto@rogers-mbp.lan: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${GREENLIGHT_DB_DSN} up

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor

# ==================================================================================== #
# BUILD
# ==================================================================================== #

## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/api ./cmd/api
	go build -ldflags='-s' -o=./bin/mac_arm64/api ./cmd/api

## build/api: build the cmd/api application
.PHONY: build/api/debug
build/api/debug:
	@echo 'Building cmd/api for debug'
	go build -gcflags "all=-N -l" -o web ./cmd/api
	dlv --listen=localhost:2345 --headless=true --api-version=2 exec web