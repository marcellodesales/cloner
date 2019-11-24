PWD = $(CURDIR)
DIST_DIR = dist
APP_NAME = cloner
ORG = marcellodesales
PUBLISH_GITHUB_USER = marcellodesales 
#PUBLISH_GITHUB_TOKEN 
PUBLISH_GITHUB_HOST = github.com
PUBLISH_GITHUB_ORG = marcellodesales

# BUILD_NUMBER expects a tag with format month-day-build_number
BUILD_NUMBER_PREFIX := $(shell date +%y.%-m)
BUILD_NUMBER := $(shell git describe --tags --abbrev=0 | cut -d . -f 3)
BUILD_NUMBER := $(shell [ ! -z "$(BUILD_NUMBER)" ] && echo $(BUILD_NUMBER) || echo "0")

# https://stackoverflow.com/questions/34142638/makefile-how-to-increment-a-variable-when-you-call-it-var-in-bash/34145171#34145171
$(eval BIN_VERSION=$(shell echo "$(BUILD_NUMBER_PREFIX).$$((  $(BUILD_NUMBER)+1 ))"))

# HELP
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

clean: ## Deletes the directory ./build
	rm -rf $(DIST_DIR)

build: clean ## Builds the docker image with binaries
	#docker build -t $(APP_NAME) --build-arg BIN_NAME=$(APP_NAME) --build-arg BIN_VERSION=$(BIN_VERSION) .
	@echo "Building next version $(BIN_VERSION)"
	BIN_VERSION=$(BIN_VERSION) docker-compose build --build-arg BIN_VERSION=$(BIN_VERSION) cli

dist: build ## Makes the dir ./dist with binaries from docker image
	@echo "Distribution libraries for version $(BIN_VERSION)"
	docker run --rm -ti --entrypoint sh -v $(PWD)/$(DIST_DIR):/bins $(ORG)/$(APP_NAME):$(BIN_VERSION) -c "cp /bin/$(APP_NAME)-darwin-amd64 /bins"
	docker run --rm -ti --entrypoint sh -v $(PWD)/$(DIST_DIR):/bins $(ORG)/$(APP_NAME):$(BIN_VERSION) -c "cp /bin/$(APP_NAME)-linux-amd64 /bins"
	docker run --rm -ti --entrypoint sh -v $(PWD)/$(DIST_DIR):/bins $(ORG)/$(APP_NAME):$(BIN_VERSION) -c "cp /bin/$(APP_NAME)-windows-amd64.exe /bins"
	ls -la $(PWD)/$(DIST_DIR)

release: dist ## Publishes the built binaries in Github Releases
	echo "Releasing next version $(BIN_VERSION)"
	git tag v$(BIN_VERSION) || true
	git push origin v$(BIN_VERSION) || true
	docker run -ti -e GITHUB_HOST=$(PUBLISH_GITHUB_HOST) -e GITHUB_USER=$(PUBLISH_GITHUB_USER) -e GITHUB_TOKEN=$(PUBLISH_GITHUB_TOKEN) -e GITHUB_REPOSITORY=$(PUBLISH_GITHUB_ORG)/$(APP_NAME) -e HUB_PROTOCOL=https -v $(PWD):/git marcellodesales/github-hub release create --prerelease --attach dist/$(APP_NAME)-darwin-amd64 --attach dist/$(APP_NAME)-linux-amd64 --attach dist/$(APP_NAME)-windows-amd64.exe -m "$(APP_NAME) $(BIN_VERSION) release" v$(BIN_VERSION)