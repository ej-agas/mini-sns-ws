SOURCES := $(wildcard *.go cmd/*/*.go)

VERSION=$(shell git describe --tags --long --dirty 2>/dev/null)

ifeq ($(VERSION),)
	VERSION=UNKNOWN
endif

app: $(SOURCES)
	go build -ldflags "-X main.version=${VERSION}" -o $@ ./cmd/app

.PHONY: lint
lint:
	golangci-lint run

docker: $(SOURCES) build/Dockerfile
	docker build -t mini-sns-ws:latest . -f build/Dockerfile --build-arg VERSION=$(VERSION)

.PHONY: publish
publish: make docker
	docker tag  sort-anim:latest ej-agas/mini-sns-ws:$(VERSION)
	docker push ej-agas/mini-sns-ws:$(VERSION)