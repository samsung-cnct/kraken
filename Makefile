NAME      := kraken
VERSION   := 1.0.8
KLIB_VER  := latest
TYPE      := stable
COMMIT    := $(shell git rev-parse HEAD)
REL_BRANCH := "$$(git rev-parse --abbrev-ref HEAD)"
PROJECT_PATH := "/go/src/github.com/samsung-cnct/kraken/"
GO_CONTAINER := "quay.io/samsung_cnct/kraken-gobuild"
GO_VERSION := 1.8.3
GO_DOCKER_SRC_VOLUME := ${GOPATH}/src:/go/src
GO_DOCKER_BIN_VOLUME := ${GOPATH}/bin:/go/bin
GOOS?=linux

build:
	@docker run -v $(GO_DOCKER_SRC_VOLUME) -e GOOS=$(GOOS) $(GO_CONTAINER):$(GO_VERSION) sh -c \
	'cd $(PROJECT_PATH) && \
	 go build -ldflags "\
	    -X github.com/samsung-cnct/kraken/cmd.KrakenMajorMinorPatch=$(VERSION) \
		-X github.com/samsung-cnct/kraken/cmd.KrakenType=$(TYPE) \
		-X github.com/samsung-cnct/kraken/cmd.KrakenGitCommit=$(COMMIT) \
		-X github.com/samsung-cnct/kraken/cmd.KrakenlibTag=$(KLIB_VER)"'

compile:
	@rm -rf build/
	@docker run -v $(GO_DOCKER_SRC_VOLUME) $(GO_CONTAINER):$(GO_VERSION) sh -c \
	'cd $(PROJECT_PATH) && \
	 gox -ldflags "\
            -X github.com/samsung-cnct/kraken/cmd.KrakenMajorMinorPatch=$(VERSION) \
            -X github.com/samsung-cnct/kraken/cmd.KrakenType=$(TYPE) \
            -X github.com/samsung-cnct/kraken/cmd.KrakenGitCommit=$(COMMIT) \
            -X github.com/samsung-cnct/kraken/cmd.KrakenlibTag=$(KLIB_VER)" \
        -osarch="linux/386" \
        -osarch="linux/amd64" \
        -osarch="darwin/amd64" \
        -output "build/{{.Dir}}_$(VERSION)_{{.OS}}_{{.Arch}}/$(NAME)" \
	 ./...'

install:
	@docker run -v $(GO_DOCKER_SRC_VOLUME) -v $(GO_DOCKER_BIN_VOLUME) $(GO_CONTAINER):$(GO_VERSION) sh -c \
	'cd $(PROJECT_PATH) && \
	 go install -ldflags "\
	    -X github.com/samsung-cnct/kraken/cmd.KrakenMajorMinorPatch=$(VERSION) \
        -X github.com/samsung-cnct/kraken/cmd.KrakenType=$(TYPE) \
        -X github.com/samsung-cnct/kraken/cmd.KrakenGitCommit=$(COMMIT)"'

deps:
	go get github.com/mitchellh/gox

dist: compile
	$(eval FILES := $(shell ls build))
	@rm -rf dist && mkdir dist
	@for f in $(FILES); do \
		(cd $(shell pwd)/build/$$f && tar -cvzf ../../dist/$$f.tar.gz *); \
		(cd $(shell pwd)/dist && shasum -a 512 $$f.tar.gz > $$f.sha512); \
		echo $$f; \
	done

release: dist
	@latest_tag=$$(git describe --tags `git rev-list --tags --max-count=1`); \
	comparison="$$latest_tag..HEAD"; \
	if [ -z "$$latest_tag" ]; then comparison=""; fi; \
	changelog=$$(git log $$comparison --oneline --no-merges --reverse); \
	github-release samsung-cnct/$(NAME) $(VERSION) $(REL_BRANCH) "**Changelog**<br/>$$changelog" 'dist/*'; \


verify:
	./bin/verify.sh; \

.PHONY: build compile install deps dist release
