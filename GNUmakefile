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

docs:
	go generate

build: fmtcheck
	go install

clean:
	go clean -cache -testcache ./...

clean-all:
	go clean -cache -testcache -modcache ./...

sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test $(TEST) -sweep=$(SWEEP) $(SWEEPARGS)

# Unit tests: offline, run against an in-memory mock of the ZPA API (see
# zpa/unit_fixture_test.go). No credentials, no tenant changes. Tests that need
# a tenant (TestAcc*, the forced sweeper) are not run. TF_ACC is cleared here so
# an exported TF_ACC never turns a unit run into API calls.

test:
	@$(MAKE) test\:unit

test\:unit:
	@echo "$(COLOR_ZSCALER)Running zpa unit tests (no API calls)...$(COLOR_NONE)"
	TF_ACC= go test ./zpa -count=1 -race $(TESTARGS) -skip '^TestAcc|^TestRunForcedSweeper$$' \
		-cover -coverprofile=unit-coverage.out -covermode=atomic -timeout 10m
	go tool cover -func unit-coverage.out | grep total:

# Acceptance tests: real API calls against the ZPA tenant configured in the
# environment (ZSCALER_CLIENT_ID, ZSCALER_CLIENT_SECRET, ZSCALER_VANITY_DOMAIN,
# ZPA_CUSTOMER_ID, ...). Each test creates and deletes tf-acc-test-* resources.
test\:acceptance:
	@echo "$(COLOR_ZSCALER)Running zpa acceptance tests (real API calls; creates and deletes resources in the tenant)...$(COLOR_NONE)"
	TF_ACC=1 go test -v -race -cover -coverprofile=zpacoverage.out -covermode=atomic ./zpa -parallel 1 -timeout 120m $(TESTARGS)
	go tool cover -func zpacoverage.out | grep total:

testacc:
	TF_ACC=1 go test $(TEST) $(TESTARGS) $(TEST_FILTER) -timeout 120m

test\:integration\:zpa:
	@echo "$(COLOR_ZSCALER)Running zpa integration tests...$(COLOR_NONE)"
	go test -v -race -cover -coverprofile=zpacoverage.out -covermode=atomic ./zpa -parallel 1 -timeout 120m
	go tool cover -html=zpacoverage.out -o zpacoverage.html
	go tool cover -func zpacoverage.out | grep total:

build13: GOOS=$(shell go env GOOS)
build13: GOARCH=$(shell go env GOARCH)
ifeq ($(OS),Windows_NT)  # is Windows_NT on XP, 2000, 7, Vista, 10...
build13: DESTINATION=$(APPDATA)/terraform.d/plugins/$(ZPA_PROVIDER_NAMESPACE)/4.4.15/$(GOOS)_$(GOARCH)
else
build13: DESTINATION=$(HOME)/.terraform.d/plugins/$(ZPA_PROVIDER_NAMESPACE)/4.4.15/$(GOOS)_$(GOARCH)
endif
build13: fmtcheck
	@echo "==> Installing plugin to $(DESTINATION)"
	@mkdir -p $(DESTINATION)
	go build -o $(DESTINATION)/terraform-provider-zpa_v4.4.15

vet:
	@echo "==> Checking source code against go vet and staticcheck"
	@go vet ./...
	@staticcheck ./...

imports:
	goimports -w $(GOFMT_FILES)

fmt: tools # Format the code
	@echo "formatting the code with $(GOFMT)..."
	@$(GOFMT) -l -w .

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

fmt-docs:
	@echo "✓ Formatting code samples in documentation"
	@terrafmt fmt -p '*.md' .

vendor-status:
	@govendor status

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

lint:
	@echo "==> Checking source code against linters..."
	@$(TFPROVIDERLINT) \
		-c 1 \
		-AT001 \
    	-R004 \
		-S001 \
		-S002 \
		-S003 \
		-S004 \
		-S005 \
		-S007 \
		-S008 \
		-S009 \
		-S010 \
		-S011 \
		-S012 \
		-S013 \
		-S014 \
		-S015 \
		-S016 \
		-S017 \
		-S019 \
		./$(PKG_NAME)

tools:
	@which $(GOFMT) || go install mvdan.cc/gofumpt@v0.9.2
	@which $(TFPROVIDERLINT) || go install github.com/bflad/tfproviderlint/cmd/tfproviderlint@latest
	@which $(STATICCHECK) || go install honnef.co/go/tools/cmd/staticcheck@latest

tools-update:
	@go install mvdan.cc/gofumpt@v0.9.2
	@go install github.com/bflad/tfproviderlint/cmd/tfproviderlint@latest
	@go install honnef.co/go/tools/cmd/staticcheck@latest

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.PHONY: build test testacc vet fmt fmtcheck errcheck tools vendor-status test-compile website-lint website website-test

