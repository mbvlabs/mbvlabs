ARG GO_VERSION=1.26.4
ARG NODE_VERSION=24

FROM node:${NODE_VERSION}-bookworm AS assets-builder

WORKDIR /usr/src/app

COPY package.json package-lock.json ./
RUN npm ci

COPY css ./css
COPY resources ./resources
COPY views ./views
COPY vite.config.ts tsconfig.json components.json ./

RUN ./node_modules/.bin/vite build

FROM golang:${GO_VERSION}-bookworm AS go-builder

WORKDIR /usr/src/app

COPY . .
COPY --from=assets-builder /usr/src/app/assets/dist ./assets/dist

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /usr/local/bin/run-app ./cmd/app

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=go-builder /usr/local/bin/run-app /usr/local/bin/run-app

WORKDIR /app

COPY --from=go-builder /usr/src/app/views/root.go.html ./views/root.go.html

EXPOSE 8080

CMD ["run-app"]
