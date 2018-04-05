VERSION=$(shell git describe --tags)
DOCKER_REPO=""
pingdom2stats: *.go
	go build -ldflags "-X main.version=$(VERSION)"

version:
	@echo $(VERSION)

run: pingdom2stats
	./pingdom2stats --help

pingdom2stats-osx: *.go
	GOOS=darwin GOARCH=amd64 go build -o $@ -ldflags "-X main.VERSION=$(VERSION)"

pingdom2stats-armv6: *.go
	GOOS=linux GOARCH=arm GOARM=6 go build -o $@ -ldflags "-X main.VERSION=$(VERSION)"

pingdom2stats-x86: *.go
	GOOS=linux GOARCH=386 go build -o $@ -ldflags "-X main.VERSION=$(VERSION)"

pingdom2stats-amd64: *.go
	GOOS=linux GOARCH=amd64 go build -o $@ -ldflags "-X main.VERSION=$(VERSION)"

pingdom2stats-docker:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o $@ -ldflags "-X main.VERSION=$(VERSION)"
	docker build --tag $(DOCKER_REPO)pingdom2stats:$(VERSION) .
	docker tag $(DOCKER_REPO)pingdom2stats:$(VERSION) $(DOCKER_REPO)pingdom2stats:latest

build-all: pingdom2stats-armv6 pingdom2stats-osx pingdom2stats-x86 pingdom2stats-amd64
# pingdom2stats-docker


.PHONY: version run build-all pingdom2stats-docker
