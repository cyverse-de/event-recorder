FROM golang:1.25

ENV CGO_ENABLED=0

WORKDIR /src/event-recorder
COPY . .
RUN apt-get update && \
    apt install just && \
    just test && \
    just build

FROM gcr.io/distroless/static-debian13

WORKDIR /app

COPY --from=0 /src/event-recorder/event-recorder /bin/event-recorder

ENTRYPOINT ["event-recorder"]
