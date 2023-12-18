FROM golang:1.21.5-alpine3.19 as build

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build
COPY . .
RUN go build -o httpServer http_server.go

FROM alpine:3.19
COPY --from=build /build/httpServer /bin/httpServer
ENTRYPOINT ["/bin/httpServer"]
