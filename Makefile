PROJECT_NAME := "filler"
PKG_LIST := $(shell go list ./... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)
VERSION := $(shell cat version)

.SHELLFLAGS = -c # Run commands in a -c flag
.SILENT: ; # no need for @
.ONESHELL: ; # recipes execute in same shell
.PHONY: install all build clean test coverage coverhtml lint fmt cpd errcheck staticcheck

default: build
all: test build

clean:
	go clean -i ./...
	rm -vf \
	"./${PROJECT_NAME}" \
	 ./coverage.* \
	  ./cpd.*

dependencies: dep
	@go mod download

run:
	go run main.go

build: ## Build the binary file
	go build -i -v $(PKG_LIST)

build-out: ## Build the binary file
	go build -i -v $(PKG_LIST)  -o $(PROJECT_NAME)_$(VERSION)

install: build ## Build the binary file
	go install ${LDFLAGS}

lint: ## exevute lint
	@golangci-lint run ./...

fmt: ## formmat the files
	@go fmt ${PKG_LIST}

cpd:  ## cpd
	dupl -t 200 -html >cpd.html

test: ## execute test
	@echo "go test ${PKG_LIST}"
	go test -i ${PKG_LIST} || exit 1
	echo ${PKG_LIST} | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

race:  ## Run data race detector
	@go test -race -short ${PKG_LIST}

bench:  ## run benchmarks
	go test -bench ${PKG_LIST}

msan:  ## Run memory sanitizer
	@go test -msan -short ${PKG_LIST}

vet: ## Vet examines Go source code and reports suspicious constructs, such as Printf calls whose arguments do not align with the format string
	@echo "go vet ."
	@go vet ${PKG_LIST} ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

security: ## Execute go sec security step
	gosec -tests ./...


misspell :## One way of improving the accuracy of your writing is to spell things right.
	@misspell -locale US  .

coverage: ## Generate global code coverage report
	./scripts/coverage.sh;

dep: ## Get the dependencies
	go get -v -u github.com/mattn/goveralls && \
	go get -v -u github.com/mibk/dupl && \
	go get -v -u github.com/client9/misspell/cmd/misspell && \
    go get github.com/securego/gosec/cmd/gosec

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
