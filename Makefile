.PHONY: clean build-linux-windows build-macos release

clean: 
	rm -rf dist/ dist-linux-win/ dist-macos/

# Build Linux & Windows versions in Docker
build-linux-windows:
	docker build -f Dockerfile.build -t goreleaser-cross-x11 .
	docker run --rm \
	  -v "$$PWD:/app" -w /app goreleaser-cross-x11 \
	  build --config .goreleaser.yaml --id linux --id windows --snapshot
	mv dist/ dist-linux-win/

# Build Mac versions natively (must be ran on MacOS)
build-macos:
	goreleaser build --config .goreleaser.yaml --id macos --snapshot
	mv dist/ dist-macos/

# Merge dist folders into one
merge:
	mkdir -p dist/
	rsync -a --exclude=artifacts.json --exclude=config.yaml dist-linux-win/ dist/
	rsync -a --exclude=artifacts.json --exclude=config.yaml dist-macos/ dist/
	jq -s 'add | unique_by(.name, .path)' \
	  dist-linux-win/artifacts.json dist-macos/artifacts.json > dist/artifacts.json
	yq eval ' .builds as $$b1 | (load("dist-macos/config.yaml") | .builds) as $$b2 | .builds = ($$b1 + $$b2)' dist-linux-win/config.yaml > dist/config.yaml

# Release the built artifacts
release: clean build-linux build-native
	goreleaser release --config .goreleaser.yaml --skip=build
