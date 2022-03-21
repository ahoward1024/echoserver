SHELL = /bin/bash

PROJECT_NAME = $(shell yq e '.name' ./k3d-config.yaml) # Gets the name of the cluster from the k3d config.

DOCKER_TAG_VERSION = $(shell grep "VERSION" Dockerfile | sed "s/^.*=\(.*\)$$/\1/g")

.PHONY: uninstall
## Delete all of the plugins/tools versions of the asdf apps. Mainly for testing if install works.
uninstall:
	while read line; do\
		read app version <<< "$${line}"; \
		asdf uninstall $${app} $${version}; \
	done < ./.tool-versions

.PHONY: install
## Install all plugins/tool versions of asdf apps.
install:
	while read line; do \
		read app version <<< "$${line}"; \
		asdf plugin add $${app}; \
		asdf install $${app} $${version}; \
	done < ./.tool-versions

.PHONY: begin
## Install all dependencies and create the cluster, installing all of the infrastructure.
begin: install create kustomize

.PHONY: create
## Creates the cluster and adds the baseline resources.
create:
	k3d cluster create --config ./k3d-config.yaml

.PHONY: kustomize
## Install Kustomize files.
kustomize:
	for directory in $$(ls ./infrastructure/local/kustomize); do \
		kubectl apply --kustomize ./infrastructure/local/kustomize/$${directory}; \
	done

.PHONY: delete
## WARNING: Destroys the cluster entirely!
destroy:
	k3d cluster delete --config ./k3d-config.yaml

.PHONY: start
## Starts the cluster from suspension. The cluster must be created first before this can work.
start:
	k3d cluster start ${PROJECT_NAME}

.PHONY: stop
## Stops the cluster, putting into suspension. The cluster must have been created and running for this to work.
stop:
	k3d cluster stop ${PROJECT_NAME}

.PHONY: up
## Start the Tilt server.
up:
	tilt up

.PHONY: exec
## Exec into a pod in the cluster for testing purposes.
exec:
	kubectl --namespace default exec --stdin=true --tty=true pod/exec-pod -- sh

.PHONY: watch
watch:
## Run gow to restart the server on file saves
	gow run main.go

build:
## Build the server for the current architecture
	mkdir -p build
	go build -o build/echoserver

build-linux:
## Build the server for the linux architecture
	mkdir -p build/linux
	GOOS=linux GOARCH=amd64 go build -o build/linux/echoserver

.PHONY: build-docker
build-docker:
## Build a Docker image of the server
	docker build --tag echoserver:${DOCKER_TAG_VERSION} .

.PHONY: run
run: build-docker
## Run the application in a container
	docker run --rm --publish 127.0.0.1:8080:8080 --name echoserver echoserver:${DOCKER_TAG_VERSION}


# Help
# COLORS
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)
TARGET_MAX_CHAR_NUM=20
.PHONY: help
## Show help
help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${YELLOW}%-$(TARGET_MAX_CHAR_NUM)s${RESET} ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

