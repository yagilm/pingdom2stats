VERSION=$(shell git describe --tags)
DOCKER_REPO=""
pingdom2mysql: *.go
	go build -ldflags "-X main.version=$(VERSION)"

version:
	@echo $(VERSION)

run: pingdom2mysql
	./pingdom2mysql --help

pingdom2mysql-osx: *.go
	GOOS=darwin GOARCH=amd64 go build -o $@ -ldflags "-X main.VERSION=$(VERSION)"

pingdom2mysql-armv6: *.go
	GOOS=linux GOARCH=arm GOARM=6 go build -o $@ -ldflags "-X main.VERSION=$(VERSION)"

pingdom2mysql-x86: *.go
	GOOS=linux GOARCH=386 go build -o $@ -ldflags "-X main.VERSION=$(VERSION)"

pingdom2mysql-amd64: *.go
	GOOS=linux GOARCH=amd64 go build -o $@ -ldflags "-X main.VERSION=$(VERSION)"

pingdom2mysql-docker:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o $@ -ldflags "-X main.VERSION=$(VERSION)"
	docker build --tag $(DOCKER_REPO)pingdom2mysql:$(VERSION) .
	docker tag $(DOCKER_REPO)pingdom2mysql:$(VERSION) $(DOCKER_REPO)pingdom2mysql:latest

build-all: pingdom2mysql-armv6 pingdom2mysql-osx pingdom2mysql-x86 pingdom2mysql-amd64
# pingdom2mysql-docker


.PHONY: version run build-all pingdom2mysql-docker
