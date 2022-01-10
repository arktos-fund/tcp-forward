.PHONY: build
# Build dev
build:
	@goreleaser release --snapshot --skip-publish --rm-dist -f .goreleaser.yaml