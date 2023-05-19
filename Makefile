BIN="./bin"
SRC=$(shell find . -name "*.go")

ifeq (, $(shell which golangci-lint))
$(warning "could not find golangci-lint in $(PATH), run: curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh")
endif

.PHONY: pre_clean lint test

all: pre_clean lint test darwin linux linux-arm windows

fmt:
	$(info ******************** checking formatting ********************)
	@test -z $(shell gofmt -l $(SRC)) || (gofmt -d $(SRC); exit 1)

lint:
	$(info ******************** running lint tools ********************)
	golangci-lint run -v

test:
	go test -covermode=atomic -coverprofile=coverage.out ./...

pre_clean:
	rm -rf bin

darwin:
	VERSION=$$(cat VERSION)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-X github.com/ditschedev/swag-ts/internal/config.Version=$$(cat VERSION) -X github.com/ditschedev/swag-ts/internal/config.Date=`TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ'`" -o ./bin/swag-ts-darwin-x86_64
	chmod +x ./bin/swag-ts-darwin-x86_64

linux:
	VERSION=$$(cat VERSION)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X github.com/ditschedev/swag-ts/internal/config.Version=$$(cat VERSION) -X github.com/ditschedev/swag-ts/internal/config.Date=`TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ'`" -o ./bin/swag-ts-linux-x86_64
	chmod +x ./bin/swag-ts-linux-x86_64

linux-arm:
	VERSION=$$(cat VERSION)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-X github.com/ditschedev/swag-ts/internal/config.Version=$$(cat VERSION) -X github.com/ditschedev/swag-ts/internal/config.Date=`TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ'`" -o ./bin/swag-ts-linux-arm64
	chmod +x ./bin/swag-ts-linux-arm64

windows:
	VERSION=$$(cat VERSION)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-X github.com/ditschedev/swag-ts/internal/config.Version=$$(cat VERSION) -X github.com/ditschedev/swag-ts/internal/config.Date=`TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ'`" -o ./bin/swag-ts-windows-x86_64
	chmod +x ./bin/swag-ts-windows-x86_64