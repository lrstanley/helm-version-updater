.DEFAULT_GOAL := build-all

export PROJECT := "helm-version-updater"
export PACKAGE := "github.com/lrstanley/helm-version-updater"

license:
	curl -sL https://liam.sh/-/gh/g/license-header.sh | bash -s

prepare: clean go-prepare
	@echo

build-all: prepare go-build
	@echo

clean:
	/bin/rm -rfv ${PROJECT}

go-fetch:
	go mod download
	go mod tidy

go-upgrade-deps:
	go get -u ./...
	go mod tidy

go-upgrade-deps-patch:
	go get -u=patch ./...
	go mod tidy

go-prepare: go-fetch
	go generate -x ./...

go-dlv: go-prepare
	dlv debug \
		--headless --listen=:2345 \
		--api-version=2 --log \
		--allow-non-terminal-interactive \
		${PACKAGE} -- --debug

go-debug: go-prepare
	go run ${PACKAGE} --debug

go-build: go-prepare
	CGO_ENABLED=0 \
	go build \
		-ldflags '-d -s -w -extldflags=-static' \
		-tags=netgo,osusergo,static_build \
		-installsuffix netgo \
		-trimpath \
		-o ${PROJECT} \
		${PACKAGE}
