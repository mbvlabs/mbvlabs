ARG NODE_VERSION=24
ARG GO_VERSION=1.27.3

FROM node:${NODE_VERSION}-bookworm-slim AS frontend
WORKDIR /app
RUN npm install --global pnpm@12.9.1
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile
COPY . .
RUN pnpm build

FROM golang:${GO_VERSION}-bookworm AS backend
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /app/assets/dist ./assets/dist
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/ ./cmd/app ./cmd/queue ./cmd/ssr ./cmd/migrate

# One image, separate commands for the V2 web, River, SSR and migration processes.
FROM node:${NODE_VERSION}-bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates tzdata \
    && rm -rf /var/lib/apt/lists/*
COPY --from=backend /out/app /usr/local/bin/run-app
COPY --from=backend /out/queue /usr/local/bin/run-queue
COPY --from=backend /out/ssr /usr/local/bin/run-ssr
COPY --from=backend /out/migrate /usr/local/bin/run-migrate
WORKDIR /app
USER node
ENV HOST=0.0.0.0
EXPOSE 8080
CMD ["run-app"]
