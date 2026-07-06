ARG GO_VERSION=1.26.4
ARG ANDUREL_VERSION=v1.0.0-beta.5

FROM golang:${GO_VERSION}-bookworm AS builder
ARG ANDUREL_VERSION

WORKDIR /usr/src/app

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    nodejs \
    npm \
    && rm -rf /var/lib/apt/lists/*

RUN go install github.com/mbvlabs/andurel@${ANDUREL_VERSION}

COPY . .

RUN andurel build

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /usr/src/app/mbvlabs /usr/local/bin/run-app

WORKDIR /app

EXPOSE 8080

CMD ["run-app"]
