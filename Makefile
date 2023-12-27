VERSION := 1.0.0
SHELL := /usr/bin/env bash
MAKEFLAGS += --silent
ARCHS := amd64 arm64
OS_LIST := windows linux darwin
DOCKER_COMPOSE := docker-compose --log-level ERROR

.PHONY: clean \
				test \
				build

clean:
	rm -rf $(PWD)/out;

build: clean
build:
	sha=$$(git rev-parse --short HEAD); \
	COMMIT_SHA="$$sha" VERSION=$(VERSION) $(DOCKER_COMPOSE) run --rm copy-source || exit 1; \
	mkdir $(PWD)/out; \
	for arch in $(ARCHS); \
	do \
		for os in $(OS_LIST); \
		do \
			>&2 echo "===> Building: $${os}-$${arch}"; \
			GOOS=$$os GOARCH=$$arch VERSION=$(VERSION) $(DOCKER_COMPOSE) run --rm build; \
		done; \
	done

test:
	$(DOCKER_COMPOSE) run --rm unit-tests
