.DEFAULT_GOAL := help
.PHONY: help
help: ## Display make tasks
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

NAME        := kraken
VERSION     := 1.2.3
KLIB_VER    ?= latest
TYPE        := stable
COMMIT      := $(shell git rev-parse HEAD)
REL_BRANCH  := "$$(git rev-parse --abbrev-ref HEAD)"
LDFLAGS     := -X github.com/samsung-cnct/kraken/cmd.KrakenMajorMinorPatch=$(VERSION) \
			   -X github.com/samsung-cnct/kraken/cmd.KrakenType=$(TYPE) \
			   -X github.com/samsung-cnct/kraken/cmd.KrakenGitCommit=$(COMMIT) \
			   -X github.com/samsung-cnct/kraken/cmd.KrakenlibTag=$(KLIB_VER)

.PHONY: bootstrap
bootstrap: setup ## get tools needed for local project development work
	go get github.com/jteeuwen/go-bindata/...

.PHONY: setup
setup: ## get tools needed for vet, test, build, and other CI/CD tasks
	go get github.com/golang/lint/golint
	go get honnef.co/go/tools/cmd/gosimple

.PHONY: vet
vet: ## validate code and configuration
	go vet main.go
	go vet ./cmd/
	go fmt main.go
	go fmt ./cmd/
	golint -set_exit_status main.go
	golint -set_exit_status ./cmd/
	gosimple main.go
	gosimple ./cmd/

.PHONY: unit-test
unit-test: ## run unit tests
	go test -v -race ./...

.PHONY: accpt-test-aws
accpt-test-aws: HOME = ${PWD}
accpt-test-aws: ## run acceptance tests for AWS (set CI_JOB_ID for local testing)
	hack/accpt_test aws

.PHONY: accpt-test-gke
accpt-test-gke: ## run acceptance tests for GKE (set CI_JOB_ID for local testing)
	hack/accpt_test gke

.PHONY: build
build: ## build the golang executable for the target archtectures
	-rm -rf build dist && mkdir build && mkdir dist
	env CGO_ENABLED=0 GOARCH="amd64" GOOS="${target}" go build -o "./build/$(NAME)-$(VERSION)-${target}-amd64" --ldflags '$(LDFLAGS)'; \
	env CGO_ENABLED=0 GOARCH="amd64" GOOS="${target}" tar -czf "./dist/$(NAME)-$(VERSION)-${target}-amd64.tgz" "./build/$(NAME)-$(VERSION)-${target}-amd64"; \
	env CGO_ENABLED=0 GOARCH="amd64" GOOS="${target}" shasum -a 512 "./build/$(NAME)-$(VERSION)-${target}-amd64" > "./dist/$(NAME)-$(VERSION)-${target}-amd64.sha512"; \

.PHONY: local_build
local_build: ## build for your machine
	CGO_ENABLED=0 go build .

.PHONY: clean
clean: ## Cleanup after make compile
	-rm -rf build dist

.PHONY: release
release: build ## Create a GitHub release
	@latest_tag=$$(git describe --tags `git rev-list --tags --max-count=1`); \
	comparison="$$latest_tag..HEAD"; \
	if [ -z "$$latest_tag" ]; then comparison=""; fi; \
	changelog=$$(git log $$comparison --oneline --no-merges --reverse); \
	github-release samsung-cnct/$(NAME) $(VERSION) $(REL_BRANCH) "**Changelog**<br/>$$changelog" 'dist/*'; \

.PHONY: regenerate-bindata
regenerate-bindata: ## Regnerate cmd/bindata.go after changes in ./data/
	go-bindata data/
	sed s/package\ main/package\ cmd/ < bindata.go > cmd/bindata.go
	rm bindata.go
