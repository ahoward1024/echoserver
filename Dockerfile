# echoserver production Docker image
# Based on https://github.com/chemidy/smallest-secured-golang-docker-image
ARG VERSION=0.1.0

# -------- Build image -------- #
# golang:1.17.7-alpine3.15
FROM golang@sha256:4e84e7209cf9d803fdb936260e37811a09c2b5107a1cc88db1bf6f1c9f3be9eb as builder
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates
ENV USER=appuser
ENV UID=10001
# See https://stackoverflow.com/a/55757473/12429735
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"
WORKDIR $GOPATH/src/echoserver/echoserver/
COPY go.mod .
ENV GO111MODULE=on
RUN go mod download
RUN go mod verify
COPY src .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -o /go/bin/echoserver .


# -------- Runtime image -------- #
FROM scratch
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /go/bin/echoserver /go/bin/echoserver
USER appuser:appuser
ENTRYPOINT ["/go/bin/echoserver"]
