FROM golang:1.18.4-alpine AS builder
RUN apk add --no-cache git
RUN go install -v github.com/Nihility981/httpx_removebanner/cmd/httpx@latest

FROM alpine:3.16.1
RUN apk -U upgrade --no-cache \
    && apk add --no-cache bind-tools ca-certificates
COPY --from=builder /go/bin/httpx /usr/local/bin/

ENTRYPOINT ["httpx"]
