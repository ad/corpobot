FROM golang:alpine as builder
RUN apk update && apk add --no-cache git ca-certificates tzdata gcc musl-dev && update-ca-certificates

ENV USER=appuser
ENV UID=10001

RUN adduser --disabled-password --gecos "" --home "/nonexistent" --shell "/sbin/nologin" --no-create-home --uid "${UID}" "${USER}"
WORKDIR $GOPATH/src/app/
COPY . .

RUN GOOS=linux GOARCH=amd64 go build -mod=vendor -ldflags='-w -s -extldflags "-static"' -a -o /go/bin/corpobot .

FROM scratch
EXPOSE 8080
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /go/bin/corpobot /go/bin/corpobot

USER appuser:appuser

ENTRYPOINT ["/go/bin/corpobot"]