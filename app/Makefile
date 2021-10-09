.PHONY: all
all: build

GOOS=linux
GOARCH=amd64

.PHONY: build
build:
	GOOS=${GOOS} GOARCH=${GOARCH} go build -o bin/server ./server
	GOOS=${GOOS} GOARCH=${GOARCH} go build -o bin/client  ./client
