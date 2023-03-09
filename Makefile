.PHONY: build tag release test lint

GOLANG_CI_LINT_VERSION=v1.50.1

build:
	./scripts/build.sh

tag:
	./scripts/tag.sh

release:
	./scripts/release.sh

test:
	go test ./...

lint:
	@which golangci-lint > /dev/null 2>&1 || (curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | bash -s -- -b "$(go env GOPATH)/bin" "$(GOLANG_CI_LINT_VERSION)")
	golangci-lint run -v --timeout=10m

.PHONY: license-update
license-update:
	@echo "==> Updating license information..."
	@rm -rf licenses
	@go-licenses save "github.com/snyk/snyk-iac-capture" --save_path="licenses" --ignore "github.com/snyk/snyk-iac-capture"
