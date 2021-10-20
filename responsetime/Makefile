.PHONY: all
all: build

GOOS=linux
GOARCH=amd64

.PHONY: build
build:
	GOOS=${GOOS} GOARCH=${GOARCH} go build -o bin/server ./server
	GOOS=${GOOS} GOARCH=${GOARCH} go build -o bin/client  ./client

PRIVATE_KEY=compare-http.pem
USER=ec2-user
SERVER_IP_ADDRESS=
CLIENT_IP_ADDRESS=
.PHONY: deploy
deploy:
	scp -ri ~/.ssh/${PRIVATE_KEY} ./bin ${USER}@${SERVER_IP_ADDRESS}:~/compare-http/responsetime
	scp -ri ~/.ssh/${PRIVATE_KEY} ./bin ${USER}@${CLIENT_IP_ADDRESS}:~/compare-http/responsetime

.PHONY: asset
asset:
	mkdir -p server/assets
	for i in 1 10 100 1000 ; do\
  		dd if=/dev/zero of=server/assets/$${i}M.txt bs=1000000 count=$$i ;\
	done
