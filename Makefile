BUILD = `date +%FT%T%z`
HASH = `git rev-parse HEAD`
BRANCH := $(shell git name-rev --name-only HEAD)
BUILD_VERSION := $(subst tags/,,$(BRANCH))

LDFLAGS=-ldflags " -s -X main.BuildTime=${BUILD} -X main.Version=${BUILD_VERSION} -X main.GitCommit=${HASH}"

build: clean
	cp ./configs/config.yaml.dist dist/configs/config.yaml.dist
	cp ./README.md dist/README.md

	go build ${LDFLAGS} -o ./dist/bin/superchutou ./cmd/*.go

run:
	go run --race ./cmd/*.go -c ./configs/config.yaml

clean:
	rm -rf dist
	mkdir dist
	mkdir dist/bin
	mkdir dist/logs
	mkdir dist/configs
