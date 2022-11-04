all : build-bin build-image
.PHONY : all build-price-periodic-collector build-influxdb-datasource build-datasource-gw \
build-binance-datasource build-docker-image-binance-datasource build-docker-image-influxdb-datasource \
build-docker-image-datasource-gw build-docker-image-price-periodic-collector build-bin build-image\
build-n-up up run-binance-datasource run-influxdb-datasource run-datasource-gw run-price-periodic-collector \
docker-compose go-test

SHELL := /bin/bash

GIT_TAG := $(shell git tag -l --contains HEAD)
ifeq ($(GIT_TAG),)
GIT_TAG := none
endif

GIT_COMMIT_ID = $(shell git rev-parse --short HEAD)
IS_TAINTED := $(shell git diff --no-ext-diff --quiet --exit-code; echo $$?)
DATE=$(shell date +%s)
VERSION=$(GIT_TAG)-$(GIT_COMMIT_ID)-$(IS_TAINTED)

build-bin: build-binance-datasource build-influxdb-datasource build-datasource-gw build-price-periodic-collector
build-image: build-docker-image-binance-datasource build-docker-image-influxdb-datasource \
build-docker-image-datasource-gw build-docker-image-price-periodic-collector
build-n-up: build-image docker-compose
up: docker-compose

build-binance-datasource:
	go build -o bin/binance-datasource -ldflags '-X "main.Version=$(VERSION)"' ./cmd/binance-datasource

build-influxdb-datasource:
	go build -o bin/influxdb-datasource -ldflags '-X "main.Version=$(VERSION)"' ./cmd/influxdb-datasource

build-datasource-gw:
	go build -o bin/datasource-gw -ldflags '-X "main.Version=$(VERSION)"' ./cmd/datasource-gw

build-price-periodic-collector:
	go build -o bin/price-periodic-collector -ldflags '-X "main.Version=$(VERSION)"' ./cmd/price-periodic-collector

build-docker-image-binance-datasource:
	docker build --quiet . -f cmd/binance-datasource/Dockerfile -t ty2/binance-datasource:$(VERSION) -t ty2/binance-datasource:dev

build-docker-image-influxdb-datasource:
	docker build --quiet . -f cmd/influxdb-datasource/Dockerfile -t ty2/influxdb-datasource:$(VERSION) -t ty2/influxdb-datasource:dev

build-docker-image-datasource-gw:
	docker build --quiet . -f cmd/datasource-gw/Dockerfile -t ty2/datasource-gw:$(VERSION) -t ty2/datasource-gw:dev

build-docker-image-price-periodic-collector:
	docker build --quiet . -f cmd/price-periodic-collector/Dockerfile -t ty2/price-periodic-collector:$(VERSION) -t ty2/price-periodic-collector:dev

run-binance-datasource:
	source ./set-local-env.sh && go run ./cmd/binance-datasource/*.go

run-influxdb-datasource:
	source ./set-local-env.sh && go run ./cmd/influxdb-datasource/*.go

run-datasource-gw:
	source ./set-local-env.sh && go run ./cmd/datasource-gw/*.go

run-price-periodic-collector:
	source ./set-local-env.sh && go run ./cmd/price-periodic-collector/*.go

docker-compose:
	docker-compose up

go-test:
	source ./set-local-env.sh && go test -coverprofile test-coverage.out ./...
	go tool cover -html test-coverage.out -o test-coverage.html
