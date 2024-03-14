build:
	echo "building httpserver binary..."
	mkdir -p bin/amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/amd64 cmd/http_server.go


build-image:
	docker build -t core.harbor.domain/http-server/http_server:v10 .
	docker push core.harbor.domain/http-server/http_server:v10


create-config:
	kubectl create cm http-server-config --from-file=resources/config/config.yaml --dry-run=client -o yaml > resources/deploy/configmap.yaml


create-certs:
	openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout ./resources/certs/tls.key -out resources/certs/tls.crt -subj "/CN=cncamp.com/O=cncamp" -addext "subjectAltName = DNS:cncamp.com"
	kubectl create secret tls http-server-tls --cert=./resources/certs/tls.crt --key=./resources/certs/tls.key