VERSION=$(shell git describe --tags)

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

build-all: pingdom2mysql-armv6 pingdom2mysql-osx pingdom2mysql-x86 pingdom2mysql-amd64

.PHONY: version run build-all
