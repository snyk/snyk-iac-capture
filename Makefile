.PHONY: build release test

build:
	./scripts/build.sh

tag:
	./scripts/tag.sh

release:
	./scripts/release.sh

test:
	go test ./...