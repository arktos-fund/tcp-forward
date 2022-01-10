.PHONY: build
# Build dev
build:
	@goreleaser release --snapshot --skip-publish --rm-dist -f .goreleaser.yaml

.PHONY: release	
# Build release
release:
	@goreleaser release --rm-dist -f .goreleaser.yaml