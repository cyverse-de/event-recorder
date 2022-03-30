FROM golang:1.17-alpine

RUN apk add --no-cache make
RUN apk add --no-cache git
RUN go get -u github.com/jstemmer/go-junit-report

ENV CGO_ENABLED=0

WORKDIR /src/event-recorder
COPY . .
RUN go test ./... && \
    go build .

FROM scratch

WORKDIR /app

COPY --from=0 /src/event-recorder/event-recorder /bin/event-recorder

ENTRYPOINT ["event-recorder"]
