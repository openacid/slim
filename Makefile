# `grep -v` does not work on travis. No time to find out why -- xp 2019 Feb 22
PKGS := $(shell go list ./... | grep -v "^github.com/openacid/slim/\(vendor\|prototype\|iohelper\|serialize\|version\)")

# PKGS := github.com/openacid/slim/array \
#         github.com/openacid/slim/bit \
#         github.com/openacid/slim/trie \

SRCDIRS := $(shell go list -f '{{.Dir}}' $(PKGS))

# gofmt check vendor dir. we need to skip vendor manually
GOFILES := $(shell find $(SRCDIRS) -not -path "*/vendor/*" -name "*.go")
GO := go

check: test vet gofmt misspell unconvert staticcheck ineffassign unparam

travis: vet gofmt misspell unconvert ineffassign unparam test

test-short:
	# fail fast with severe bugs
	$(GO) test -short      $(PKGS)

test:
	# $(GO) test -tags debug $(PKGS)
	# test release version and generate coverage data for task `coveralls`.
	$(GO) test -covermode=count -coverprofile=coverage.out $(PKGS)

lint: vet gofmt misspell unconvert ineffassign unparam

vet:
	$(GO) vet $(PKGS)

staticcheck:
	$(GO) install honnef.co/go/tools/cmd/staticcheck@latest
	# ST1016: methods on the same type should have the same receiver name
	#         .pb.go have this issue.
	staticcheck -checks all,-ST1016 $(PKGS)

misspell:
	$(GO) install github.com/client9/misspell/cmd/misspell@latest
	find $(SRCDIRS) -name '*.go' -or -name '*.md' | grep -v "\bvendor/" | xargs misspell \
		-locale US \
		-error
	misspell \
		-locale US \
		-error \
		*.md *.go

unconvert:
	$(GO) install github.com/mdempsky/unconvert@latest
	unconvert -v $(PKGS)

ineffassign:
	$(GO) install github.com/gordonklaus/ineffassign@latest
	# With ineffassign v0.0.0-20210104184537-8eed68eb605f
	# The following form now complains error: "named files must all be in one directory"
	#    find $(SRCDIRS) -name '*.go' | grep -v "\bvendor/" | xargs ineffassign
	# Seems this works:
	ineffassign $(PKGS)

pedantic: check errcheck

unparam:
	$(GO) install mvdan.cc/unparam@latest
	unparam ./...

errcheck:
	$(GO) install github.com/kisielk/errcheck@latest
	errcheck $(PKGS)

gofmt:
	@echo Checking code is gofmted
	@test -z "$(shell gofmt -s -l -d -e $(GOFILES) | tee /dev/stderr)"

ben: test
	$(GO) test ./... -run=none -bench=. -benchmem

gen:
	$(GO) generate ./...

readme:
	python ./scripts/build_md.py
	# brew install nodejs
	# npm install -g doctoc
	doctoc --title '' --github README.md

fix:
	gofmt -s -w $(GOFILES)
	unconvert -v -apply $(PKGS)

# local coverage
coverage:
	$(GO) test -covermode=count -coverprofile=coverage.out $(PKGS)
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

coveralls:
	$(GO) install golang.org/x/tools/cmd/cover@latest
	$(GO) install github.com/mattn/goveralls@latest
	goveralls -ignore='**/*.pb.go' -coverprofile=coverage.out -service=travis-ci
	# -repotoken $$COVERALLS_TOKEN
