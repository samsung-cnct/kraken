NAME      := k2cli
VERSION   := 1.0.6
TYPE      := stable
COMMIT    := $(shell git rev-parse HEAD)
godep=GOPATH=$(shell godep path):${GOPATH}

build:
	@godep go build -ldflags "-X github.com/samsung-cnct/k2cli/cmd.K2CliMajorMinorPatch=$(VERSION) \
		-X github.com/samsung-cnct/k2cli/cmd.K2CliType=$(TYPE) \
		-X github.com/samsung-cnct/k2cli/cmd.K2CliGitCommit=$(COMMIT)"

compile:
	@rm -rf build/
	@$(GODEP) gox -ldflags "-X github.com/samsung-cnct/k2cli/cmd.K2CliMajorMinorPatch=$(VERSION) \
									-X github.com/samsung-cnct/k2cli/cmd.K2CliType=$(TYPE) \
									-X github.com/samsung-cnct/k2cli/cmd.K2CliGitCommit=$(COMMIT)" \
	-osarch="linux/386" \
	-osarch="linux/amd64" \
	-osarch="darwin/amd64" \
	-output "build/{{.Dir}}_$(VERSION)_{{.OS}}_{{.Arch}}/$(NAME)" \
	./...

install:
	@godep go install -ldflags "-X github.com/samsung-cnct/k2cli/cmd.K2CliMajorMinorPatch=$(VERSION) \
									-X github.com/samsung-cnct/k2cli/cmd.K2CliType=$(TYPE) \
									-X github.com/samsung-cnct/k2cli/cmd.K2CliGitCommit=$(COMMIT)"

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
	github-release samsung-cnct/$(NAME) $(VERSION) "$$(git rev-parse --abbrev-ref HEAD)" "**Changelog**<br/>$$changelog" 'dist/*'; \
	git pull

.PHONY: build compile install deps dist release
