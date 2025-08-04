# Makefile for building bpp
# Reference Guide - https://www.gnu.org/software/make/manual/make.html

#
# Internal variables or constants.
# NOTE - These will be executed when any make target is invoked.
#
IS_DOCKER_INSTALLED = $(shell which docker >> /dev/null 2>&1; echo $$?)

# Docker info
DOCKER_REGISTRY ?= docker.io
DOCKER_REPO ?= neelanjan00
DOCKER_IMAGE ?= bpp
DOCKER_TAG ?= main-latest

.PHONY: help
help:
	@echo ""
	@echo "Usage:-"
	@echo "\tmake deps          -- sets up dependencies for image build"
	@echo "\tmake push          -- pushes the machine ifs multi-arch image"
	@echo "\tmake build-amd64   -- builds the machine ifs binary & docker amd64 image"
	@echo "\tmake push-amd64    -- pushes the machine ifs amd64 image"
	@echo ""

.PHONY: all
all: deps gotasks build trivy-check

.PHONY: deps
deps: _build_check_docker godeps

_build_check_docker:
	@echo "------------------"
	@echo "--> Check the Docker deps"
	@echo "------------------"
	@if [ $(IS_DOCKER_INSTALLED) -eq 1 ]; \
		then echo "" \
		&& echo "ERROR:\tdocker is not installed. Please install it before build." \
		&& echo "" \
		&& exit 1; \
		fi;

.PHONY: godeps
godeps:
	@echo ""
	@echo "INFO: verifying dependencies for bpp build ..."
	@go get -u -v golang.org/x/lint/golint
	@go get -u -v golang.org/x/tools/cmd/goimports

.PHONY: test
test:
	@echo "------------------"
	@echo "--> Run Go Test"
	@echo "------------------"
	@go test ./... -coverprofile=coverage.txt -covermode=atomic -v

.PHONY: gotasks
gotasks: unused-package-check gofmt-check

.PHONY: unused-package-check
unused-package-check:
	@echo "------------------"
	@echo "--> Check unused packages for the bpp"
	@echo "------------------"
	@tidy=$$(go mod tidy); \
	if [ -n "$${tidy}" ]; then \
		echo "go mod tidy checking failed!"; echo "$${tidy}"; echo; \
	fi

.PHONY: gofmt-check
gofmt-check:
	@echo "------------------"
	@echo "--> Check go formatting"
	@echo "------------------"
	@gfmt=$$(gofmt -s -l . | wc -l); \
	if [ "$${gfmt}" -ne 0 ]; then \
		echo "The following files were found to be not go formatted:"; \
   		gofmt -s -l .; \
   		exit 1; \
  	fi

.PHONY: docker.buildx
docker.buildx:
	@echo "------------------------------"
	@echo "--> Setting up Builder        "
	@echo "------------------------------"
	@if ! docker buildx ls | grep -q multibuilder; then\
		docker buildx create --name multibuilder;\
		docker buildx inspect multibuilder --bootstrap;\
		docker buildx use multibuilder;\
	fi

.PHONY: build
build: docker.buildx build-multiarch

build-multiarch:
	@echo "-------------------------"
	@echo "--> Build bpp image"
	@echo "-------------------------"
	@docker buildx build . --file build/Dockerfile --build-arg BPP_PATH=./bpp --progress plain --no-cache --platform linux/arm64,linux/amd64 --tag $(DOCKER_REGISTRY)/$(DOCKER_REPO)/$(DOCKER_IMAGE):$(DOCKER_TAG)

.PHONY: push
push: docker.buildx push-multiarch

push-multiarch:
	@echo "------------------------"
	@echo "--> Push bpp image"
	@echo "------------------------"
	@echo "Pushing $(DOCKER_REPO)/$(DOCKER_IMAGE):$(DOCKER_TAG)"
	@docker buildx build . --push --file build/Dockerfile --build-arg BPP_PATH=./bpp --progress plain --no-cache --platform linux/arm64,linux/amd64 --tag $(DOCKER_REGISTRY)/$(DOCKER_REPO)/$(DOCKER_IMAGE):$(DOCKER_TAG)

.PHONY: build-amd64
build-amd64:
	@echo "-------------------------"
	@echo "--> Build bpp amd64 image"
	@echo "-------------------------"
	@docker build --file build/Dockerfile --build-arg BPP_PATH=./bpp --tag $(DOCKER_REGISTRY)/$(DOCKER_REPO)/$(DOCKER_IMAGE):$(DOCKER_TAG) . --build-arg TARGETARCH=amd64

.PHONY: push-amd64
push-amd64:
	@echo "------------------------------"
	@echo "--> Pushing bpp amd64 image"
	@echo "------------------------------"
	@docker push $(DOCKER_REGISTRY)/$(DOCKER_REPO)/$(DOCKER_IMAGE):$(DOCKER_TAG) . --build-arg TARGETARCH=amd64

.PHONY: trivy-check
trivy-check:
	@echo "------------------------"
	@echo "---> Running Trivy Check"
	@echo "------------------------"
	@./trivy --exit-code 0 --severity HIGH --no-progress $(DOCKER_REGISTRY)/$(DOCKER_REPO)/$(DOCKER_IMAGE):$(DOCKER_TAG)
	@./trivy --exit-code 0 --severity CRITICAL --no-progress $(DOCKER_REGISTRY)/$(DOCKER_REPO)/$(DOCKER_IMAGE):$(DOCKER_TAG)

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)
SWAG ?= $(LOCALBIN)/swag

.PHONY: swag
swag: $(SWAG) ## Download swag locally if necessary.
$(SWAG): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install github.com/swaggo/swag/cmd/swag@v1.16.1

.PHONY: server-api-doc
server-api-doc: swag
	$(SWAG) init --parseDependency --parseInternal -g main.go
