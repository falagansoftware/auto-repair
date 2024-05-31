# Change these variables as necessary.
MAIN_PACKAGE_PATH := ./cmd/web
BINARY_NAME := auto_repair
SHELL := /bin/bash

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

.PHONY: no-dirty
no-dirty:
	git diff --exit-code


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go test -race -buildvcs -vet=off ./...


# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## build tailwind styles: buil tailwind styles in watch mode
.PHONY: build/tailwind
build/tailwind:
	~/go/bin/tailwindcss -i ./internal/http/assets/tailwind.css -o ./internal/http/assets/styles/styles.css

## clean templ: remove go html templates
.PHONY: clean/templ
clean/templ:
	rm -rf internal/http/html/*.go

## build templ: generate go html templates
.PHONY: build/templ
build/templ:
	~/go/bin/templ generate \
		-path internal/http/html

## format templ: format templates files
.PHONY: fmt/templ
fmt/templ:
	~/go/bin/templ fmt internal/http/html

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

## build: build the application
.PHONY: build/go
build: 
	# building go
	go build -o=/tmp/bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

## run: run the  application
.PHONY: run
run: clean/templ build/templ build
	/tmp/bin/${BINARY_NAME}

## run/live: run the application with reloading on file changes
.PHONY: run/live
run/live:
	~/go/bin/air \
		--build.cmd "make build/tailwind build/templ build/go" \
		--build.bin "/tmp/bin/${BINARY_NAME}" \
		--build.delay "1000" \
		--build.exclude_dir "tmp" \
		--build.exclude_regex "_templ.go" \
		--build.include_ext "go, templ" \
		--misc.clean_on_exit "true"

# ==================================================================================== #
# OPERATIONS
# ==================================================================================== #

## push: push changes to the remote Git repository
.PHONY: push
push: tidy audit no-dirty
	git push

## production/deploy: deploy the application to production
# .PHONY: production/deploy
# production/deploy: confirm tidy audit no-dirty
# 	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=/tmp/bin/linux_amd64/${BINARY_NAME} ${MAIN_PACKAGE_PATH}
# 	upx -5 /tmp/bin/linux_amd64/${BINARY_NAME}
# 	# Include additional deployment steps here...