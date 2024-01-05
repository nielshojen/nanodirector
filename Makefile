current_dir = $(shell pwd)

SHELL = /bin/sh

ifneq ($(OS), Windows_NT)
	CURRENT_PLATFORM = linux
	ifeq ($(shell uname), Darwin)
		SHELL := /bin/sh
		CURRENT_PLATFORM = darwin
	endif
else
	CURRENT_PLATFORM = windows
endif

all: xp-build

.PHONY: postgres

.pre-build:
	mkdir -p build/darwin
	mkdir -p build/linux

lint:
	golangci-lint run

fix:
	golangci-lint run --fix

gomodcheck:
	@go help mod > /dev/null || (@echo nanodirector requires Go version 1.11 or  higher for module support && exit 1)

deps: gomodcheck
	@go mod download

clean:
	rm -rf build

build: clean .pre-build
	echo "Building..."
	go build -o build/$(CURRENT_PLATFORM)/nanodirector

xp-build:  clean .pre-build
	GOOS=darwin GOARCH=amd64 go build -o build/darwin/nanodirector-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -o build/darwin/nanodirector-darwin-arm64
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/linux/nanodirector-linux-amd64
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o build/linux/nanodirector-linux-arm64

postgres-clean:
	rm -rf postgres

postgres:
	docker rm -f nanodirector-postgres || true
	docker run --name nanodirector-postgres -p 5432:5432 -e POSTGRES_PASSWORD=password -v ${current_dir}/postgres:/var/lib/postgresql/data -d postgres:11
	docker rm -f nanodirector-redis || true
	docker run --name nanodirector-redis -d -p 6379:6379 redis
	sleep 5

redis:
	docker rm -f nanodirector-redis || true
	docker run --name nanodirector-redis -d -p 6379:6379 redis:6
	sleep 5

nanodirector_nosign: build
	build/$(CURRENT_PLATFORM)/nanodirector -nanomdmurl="${NANO_URL}" -nanomdmapikey="supersecret" -debug

curlprofile:
	rm -f EnrollmentProfile.mobileconfig
	curl -o EnrollmentProfile.mobileconfig ${NANO_URL}/mdm/enroll

nanodirector: build curlprofile
	build/$(CURRENT_PLATFORM)/nanodirector -micromdmurl="${NANO_URL}" -micromdmapikey="supersecret" -debug -sign -cert=SigningCert.p12 -key-password=password -password=secret  -db-username=postgres -db-host=127.0.0.1 -db-port=5432 -db-name=postgres -db-password=password -db-sslmode=disable -loglevel=debug -escrowurl="${ESCROW_URL}" -clear-device-on-enroll -enrollment-profile=EnrollmentProfile.mobileconfig #-enrollment-profile-signed #-log-format=json

test:
	go test -cover -v ./...