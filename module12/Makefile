build:
	echo "building httpserver binary..."
	mkdir -p bin/amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/amd64 cmd/http_server.go


build-image:
	docker build -t core.harbor.domain/http-server/http_server:v12 .
	docker push core.harbor.domain/http-server/http_server:v12


create-config:
	kubectl create cm http-server-service0-config --from-file=resources/config/config-service0.yaml --dry-run=client -o yaml > resources/deploy/service0-configmap.yaml
	kubectl create cm http-server-service1-config --from-file=resources/config/config-service1.yaml --dry-run=client -o yaml > resources/deploy/service1-configmap.yaml
	kubectl create cm http-server-service2-config --from-file=resources/config/config-service2.yaml --dry-run=client -o yaml > resources/deploy/service2-configmap.yaml


create-certs:
	openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout ./resources/certs/tls.key -out resources/certs/tls.crt -subj "/CN=cncamp.com/O=cncamp" -addext "subjectAltName = DNS:cncamp.com"
	kubectl create secret tls http-server-tls --cert=./resources/certs/tls.crt --key=./resources/certs/tls.key


create-istio-certs:
	openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout ./resources/certs/cncamp.io.key -out resources/certs/cncamp.io.crt -subj "/CN=*.cncamp.io/O=cncamp" -addext "subjectAltName = DNS:cncamp.io"
	kubectl -n istio-system create secret tls cncamp-io-tls --cert=./resources/certs/cncamp.io.crt --key=./resources/certs/cncamp.io.key


install:
	kubectl apply -f ./resources/deploy/
	helm install service0 ./resources/deploy/chart/httpserver/ -f ./resources/deploy/chart/values-service0.yaml
	helm install service1 ./resources/deploy/chart/httpserver/ -f ./resources/deploy/chart/values-service1.yaml
	helm install service2 ./resources/deploy/chart/httpserver/ -f ./resources/deploy/chart/values-service2.yaml


uninstall:
	helm uninstall service0
	helm uninstall service1
	helm uninstall service2


reinstall: uninstall install