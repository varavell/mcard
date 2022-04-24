.PHONY: clean build fmt vendor

build:
	go build main.go

fmt:
	go fmt

vendor:
	go mod vendor

clean:
	rm -rf vendor