.PHONY: clean build-linux-windows build-mac

clean: 
	rm -rf dist/

# Build Linux & Windows versions in Docker
build-linux-windows:
	docker build -f Dockerfile.build -t goreleaser-cross-x11 .
	docker run --rm \
	  -v "$$PWD:/app" -w /app goreleaser-cross-x11 \
	  build --config .goreleaser.yaml --id linux --id windows --snapshot

# Build Mac versions natively (must be ran on a mac)
build-mac:
	  goreleaser build --config .goreleaser.yaml --id macos --snapshot

release: clean build-linux build-native
	goreleaser release --config .goreleaser.yml --skip=build
