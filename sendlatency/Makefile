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
	scp -ri ~/.ssh/${PRIVATE_KEY} ./bin ${USER}@${SERVER_IP_ADDRESS}:~/compare-http/sendlatency
	scp -ri ~/.ssh/${PRIVATE_KEY} ./bin ${USER}@${CLIENT_IP_ADDRESS}:~/compare-http/sendlatency

.PHONY: asset
asset:
	mkdir -p client/assets
	# KB and MB assets
	for i in `seq 0 3` ; do\
  		dd if=/dev/zero of=client/assets/$$((10**i))K.txt bs=$$((10**3)) count=$$((10**i)) ;\
  		dd if=/dev/zero of=client/assets/$$((10**i))M.txt bs=$$((10**6)) count=$$((10**i)) ;\
	done
