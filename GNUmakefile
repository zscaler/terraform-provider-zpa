SWEEP?=global
TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep "zpa/")
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=zpa
GOFMT:=gofumpt
TFPROVIDERLINT=tfproviderlint
STATICCHECK=staticcheck
TF_PLUGIN_DIR=~/.terraform.d/plugins
ZPA_PROVIDER_NAMESPACE=zscaler.com/zpa/zpa

# Expression to match against tests
# go test -run <filter>
# e.g. Iden will run all TestAccIdentity tests
ifdef TEST_FILTER
	TEST_FILTER := -run $(TEST_FILTER)
endif

TESTARGS?=-test.v

default: build

dep: # Download required dependencies
	go mod tidy

build: fmtcheck
	go install

clean:
	go clean -cache -testcache ./...

clean-all:
	go clean -cache -testcache -modcache ./...

sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test $(TEST) -sweep=$(SWEEP) $(SWEEPARGS)

test:
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) $(TEST_FILTER) -timeout=30s -parallel=10

testacc:
	TF_ACC=1 go test $(TEST) $(TESTARGS) $(TEST_FILTER) -timeout 120m

build13: GOOS=$(shell go env GOOS)
build13: GOARCH=$(shell go env GOARCH)
ifeq ($(OS),Windows_NT)  # is Windows_NT on XP, 2000, 7, Vista, 10...
build13: DESTINATION=$(APPDATA)/terraform.d/plugins/$(ZPA_PROVIDER_NAMESPACE)/3.0.2/$(GOOS)_$(GOARCH)
else
build13: DESTINATION=$(HOME)/.terraform.d/plugins/$(ZPA_PROVIDER_NAMESPACE)/3.0.2/$(GOOS)_$(GOARCH)
endif
build13: fmtcheck
	go mod tidy && go mod vendor
	@echo "==> Installing plugin to $(DESTINATION)"
	@mkdir -p $(DESTINATION)
	go build -o $(DESTINATION)/terraform-provider-zpa_v3.0.2

lint: vendor
	@echo "✓ Linting source code with https://staticcheck.io/ ..."
	@go run honnef.co/go/tools/cmd/staticcheck@v0.4.0 ./...

coverage: test
	@echo "✓ Opening coverage for unit tests ..."
	@go tool cover -html=coverage.txt

vet:
	@echo "==> Checking source code against go vet and staticcheck"
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

imports:
	goimports -w $(GOFMT_FILES)

fmt:
	@echo "✓ Formatting source code with goimports ..."
	@go run golang.org/x/tools/cmd/goimports@latest -w $(shell find . -type f -name '*.go' -not -path "./vendor/*")
	@echo "✓ Formatting source code with gofmt ..."
	@gofmt -w $(shell find . -type f -name '*.go' -not -path "./vendor/*")

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

fmt-docs:
	@echo "✓ Formatting code samples in documentation"
	@terrafmt fmt -p '*.md' .


errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

tools:
	@echo "==> installing required tooling..."
	@sh "$(CURDIR)/scripts/gogetcookie.sh"
	GO111MODULE=off go get -u github.com/client9/misspell/cmd/misspell

vendor-status:
	@govendor status

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-lint:
	@echo "==> Checking website against linters..."
	@misspell -error -source=text website/

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.PHONY: build test testacc vet fmt fmtcheck errcheck tools vendor-status test-compile website-lint website website-test