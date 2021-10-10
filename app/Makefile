.PHONY: all
all: build

.PHONY: certificate
certificate:
	# create CA certificate
	openssl genrsa -out tls/ca.key 2048
	openssl req -new -sha256 -key tls/ca.key -out tls/ca.csr -config tls/openssl.cnf
	openssl x509 -in tls/ca.csr -days 365 -req -signkey tls/ca.key -sha256 -out tls/ca.crt -extfile tls/openssl.cnf -extensions CA

	# create server certificate
	openssl genrsa -out tls/server.key 2048
	openssl req -new -nodes -sha256 -key tls/server.key -out tls/server.csr -config tls/openssl.cnf
	openssl x509 -req -days 365 -in tls/server.csr -sha256 -out tls/server.crt -CA tls/ca.crt -CAkey tls/ca.key -CAcreateserial -extfile tls/openssl.cnf -extensions Server

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
	scp -ri ~/.ssh/${PRIVATE_KEY} ./* ${USER}@${SERVER_IP_ADDRESS}:~/compare-http
	scp -ri ~/.ssh/${PRIVATE_KEY} ./* ${USER}@${CLIENT_IP_ADDRESS}:~/compare-http
