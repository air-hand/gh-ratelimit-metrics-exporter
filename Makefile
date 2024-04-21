IS_IN_CONTAINER := $(shell test -f /.dockerenv && echo 0)

SHELL := bash

PROJECT_NAME := gh-ratelimit-metrics-exporter

.PHONY: build
build:
	go mod download && go mod tidy && CGO_ENABLED=0 go build -o build/app ./app

build-image:
	docker build -t $(PROJECT_NAME) .

.PHONY:
run:
ifeq ($(IS_IN_CONTAINER),0)
	go run ./app
else
	make devcontainer-up && devcontainer exec --workspace-folder ./ make run
endif

run-container: build-image
run-container: GH_TOKEN ?= ""
run-container:
	@docker run -d --rm -p 8080:8080 --env GH_TOKEN=$(GH_TOKEN) $(PROJECT_NAME)

.PHONY:
test:
ifeq ($(IS_IN_CONTAINER),0)
	go install github.com/matryer/moq@v0.3.4 && go generate ./app && \
	go test -v --cover ./app
else
	make devcontainer-up && devcontainer exec --workspace-folder ./ make test
endif

.PHONY:
e2e-test:
ifeq ($(IS_IN_CONTAINER),0)
	go install github.com/bitnami/wait-for-port@v1.0.7 && \
	docker kill $(PROJECT_NAME)-app || true
	make build-image && \
	docker run -d --rm -p 8080:8080 --env GH_TOKEN=$${GH_TOKEN} --name=$(PROJECT_NAME)-app $(PROJECT_NAME) && \
	wait-for-port --timeout=30 8080 && \
	curl -X GET http://localhost:8080/metrics
	docker stop $(PROJECT_NAME)-app || true
else
	make devcontainer-up && devcontainer exec --workspace-folder ./ make e2e-test
endif

.PHONY:
devcontainer-build:
	devcontainer build --workspace-folder ./

.PHONY:
devcontainer-up: devcontainer-build
	devcontainer up --workspace-folder ./

.PHONY:
devcontainer-down:
# TODO: hardcoded the project name
	docker compose -p $(PROJECT_NAME)_devcontainer down
