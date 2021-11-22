TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep "zpa/")
PKG_NAME=zpa
TF_PLUGIN_DIR=~/.terraform.d/plugins
ZPA_PROVIDER_NAMESPACE=zscaler.com/zpa/zpa

default: build

build: fmtcheck
	go install


build13: GOOS=$(shell go env GOOS)
build13: GOARCH=$(shell go env GOARCH)
ifeq ($(OS),Windows_NT)  # is Windows_NT on XP, 2000, 7, Vista, 10...
build13: DESTINATION=$(APPDATA)/terraform.d/plugins/$(ZPA_PROVIDER_NAMESPACE)/2.0.3/$(GOOS)_$(GOARCH)
else
build13: DESTINATION=$(HOME)/.terraform.d/plugins/$(ZPA_PROVIDER_NAMESPACE)/2.0.3/$(GOOS)_$(GOARCH)
endif
build13: fmtcheck
	@go mod tidy && go mod vendor
	@echo "==> Installing plugin to $(DESTINATION)"
	@mkdir -p $(DESTINATION)
	go build -o $(DESTINATION)/terraform-provider-zpa_v2.0.3

test: fmtcheck
	go test $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=600s -parallel=4

testacc: fmtcheck
	TF_ACC=true go test $(TEST) -v $(TESTARGS) -timeout 600m

vet:
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
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

tools:
	go get -u github.com/kardianos/govendor
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install

vendor-status:
	@govendor status

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

.PHONY: build test testacc vet fmt fmtcheck errcheck tools vendor-status test-compile website-lint website website-test
