.PHONY: all
all: build

GOOS=linux
GOARCH=amd64

.PHONY: build
build:
	GOOS=${GOOS} GOARCH=${GOARCH} go build -o bin/server ./server
	go build -o bin/client ./client

.PHONY: clean
clean:
	rm -f bin/*
	rm -f out/*

PRIVATE_KEY=compare-http.pem
USER=ec2-user
SERVER_IP_ADDRESS=

.PHONY: deploy
deploy:
	scp -ri ~/.ssh/${PRIVATE_KEY} ./bin ${USER}@${SERVER_IP_ADDRESS}:~/compare-http/continuoussend
