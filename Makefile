# ClawFirm Makefile.
# Convenience targets for build / test / local-up / image build.

SHELL          := /usr/bin/env bash
GO             ?= go
GOFLAGS        ?= -trimpath
LDFLAGS_VERSION = -X github.com/jimmershere/clawfirm/internal/version.Version=$(shell git describe --tags --dirty --always 2>/dev/null || echo dev)
LDFLAGS_COMMIT  = -X github.com/jimmershere/clawfirm/internal/version.Commit=$(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
LDFLAGS_DATE    = -X github.com/jimmershere/clawfirm/internal/version.Date=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS         = -s -w $(LDFLAGS_VERSION) $(LDFLAGS_COMMIT) $(LDFLAGS_DATE)

.PHONY: help build install test lint clean solo-up solo-down smb-up images

help:               ## show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-12s %s\n", $$1, $$2}'

build:              ## build the clawfirm CLI binary into ./bin/
	mkdir -p bin
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o bin/clawfirm ./cmd/clawfirm

install: build      ## install to /usr/local/bin (may require sudo)
	install -m 0755 bin/clawfirm /usr/local/bin/clawfirm

test:               ## run all Go tests
	$(GO) test ./... -race -count=1

lint:               ## run gofmt + go vet + golangci-lint (if installed)
	gofmt -l . | tee /dev/stderr | (! read)
	$(GO) vet ./...
	command -v golangci-lint >/dev/null && golangci-lint run ./... || echo "golangci-lint not installed; skipping"

clean:              ## remove build artifacts
	rm -rf bin dist

solo-up:            ## bring up the Solo tier compose stack
	cd deploy/solo && docker compose up -d

solo-down:          ## stop the Solo tier compose stack
	cd deploy/solo && docker compose down

smb-up:             ## apply the SMB Kustomize overlay to the current kube context
	kubectl apply -k deploy/smb/k3s

images:             ## build all four ClawFirm-original service images locally
	docker build -t clawfirm/mcp-gateway:dev    services/mcp-gateway
	docker build -t clawfirm/memory-service:dev services/memory-service
	docker build -t clawfirm/approval-shim:dev  services/approval-shim
	docker build -t clawfirm/policy-judge:dev   services/policy-judge
