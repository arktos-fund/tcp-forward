.PHONY: build-all
# Build all images
build-all:
	@goreleaser release --snapshot --skip-publish --rm-dist -f .goreleaser.yaml

.PHONY: build/%	
# Build dev
build/%:
	@goreleaser release --snapshot --skip-publish --rm-dist -f .goreleaser.yaml