FROM golang:1.24

ENV CGO_ENABLED=0

WORKDIR /src/event-recorder
COPY . .
RUN go test ./... && \
    ls -l && \
    go build .

FROM debian:stable-slim

WORKDIR /app

COPY --from=0 /src/event-recorder/event-recorder /bin/event-recorder

ENTRYPOINT ["event-recorder"]
