.PHONY: clean build release

clean: 
	rm -rf dist/

build:
	docker build -t pleasant-cli .

release: clean
	goreleaser release
