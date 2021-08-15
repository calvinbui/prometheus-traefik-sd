FROM golang:1.16.7-alpine3.14 as builder
WORKDIR /build
COPY . .
RUN \
    CGO_ENABLED=0 \
    go build -o main cmd/main.go

FROM scratch
COPY --from=builder /build/main /main
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
ENTRYPOINT ["/main"]
