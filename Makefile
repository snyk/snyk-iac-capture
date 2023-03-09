# Copyright 2023 Snyk Ltd.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

.PHONY: build tag release test lint license-update

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

license-update:
	@echo "==> Updating license information..."
	@rm -rf licenses
	@go-licenses save "github.com/snyk/snyk-iac-capture" --save_path="licenses" --ignore "github.com/snyk/snyk-iac-capture"
