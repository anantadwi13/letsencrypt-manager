FROM golang:1.17 AS builder

WORKDIR /go/src/letsencrypt
COPY . .
RUN go mod tidy
RUN GOOS=linux CGO_ENABLED=0 go build -o service ./cmd/service/


FROM certbot/certbot:v1.18.0

WORKDIR /root
COPY --from=builder /go/src/letsencrypt/service .

VOLUME ["/etc/letsencrypt"]

EXPOSE 80/tcp 5555/tcp

ENTRYPOINT []
CMD ./service
