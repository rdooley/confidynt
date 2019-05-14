SHELL=/bin/bash

BUILD_DIR=build
PROJECT_NAME=confidynt ## @build Project name
DEFAULT_TARGET=$(BUILD_DIR)/$(PROJECT_NAME)
MAKEFILE_LIST = Makefile

.PHONY: clean build generate install help help-variables

clean: ## @build Clean stuff
		rm -rf $(BUILD_DIR)

build $(DEFAULT_TARGET): ## @build Build the actual thing
		mkdir -p $(BUILD_DIR)
		# default build
		go build -o $(DEFAULT_TARGET)
		# build for specific operating systems
		BUILD_DIR=$(BUILD_DIR) PROJECT_NAME=$(PROJECT_NAME) ./build.sh "linux/amd64" "darwin/amd64"

generate:
	go generate ./...

install: build ## @build Install make-doc to /usr/local/bin
	sudo cp $(DEFAULT_TARGET) /usr/local/bin

help: ## @help show this help
	@make-doc $(MAKEFILE_LIST)

help-variables: ## @help show makefile customizable variables
	@make-doc $(MAKEFILE_LIST) --variables
