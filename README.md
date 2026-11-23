# MBV Labs

The MBV Labs site on Andurel V2. The public site and private admin use React
and Inertia. Public routes opt in explicitly to server-side rendering and
preserve the existing URLs, metadata, RSS and Markdown responses. API
contracts, content, and background jobs use PostgreSQL and River.

React controls use shadcn/ui's Radix Sera style. The project configuration
is in `components.json`, and the registry components are checked in under
`resources/js/Components/ui`. Add components with `pnpm dlx shadcn@4.23.0 add <component>`.
To exercise the UI against a running, seeded preview, install Chromium with
`pnpm exec playwright install chromium`, then run
`MBVLABS_PREVIEW_URL=http://localhost:8080 pnpm exec playwright test tests/seeded-preview.spec.ts`.

## Development

Requirements: Go 1.27.3, Node 22+, pnpm 12, PostgreSQL, and
`andurel` v2.0.0-rc.1 (`go install github.com/mbvlabs/andurel/v2@v2.0.0-rc.1`).

```sh
pnpm install --frozen-lockfile
cp .env.example .env
# Set HOST_PRIMARY and the database, auth, and email settings for this machine.
go run ./cmd/migrate
andurel run
```

Run `go run ./cmd/queue` for background jobs if your development runner isn't
starting the worker. On a **disposable development database only**,
`go run ./cmd/seeds` creates verified test accounts including
`admin@example.com` with password `password123`. The sign-in page is
`/users/sign-in`, and the private dashboard is `/admin`. There is no public
registration.

```sh
pnpm typecheck
pnpm build
go vet ./...
go test ./... # Database-backed package suites require Docker.
```

The V2 generators are `andurel sync routes`, `andurel sync payloads`,
`andurel sync queries`, and `andurel sync emails`.

## Deploying over V1

1. Snapshot the PostgreSQL database and rehearse against a restorable copy.
   Keep the existing Goose history; do not set up a new, empty database for
   the production cutover. Versions 1–15 follow the legacy application;
   versions 16 and 17 upgrade River and add V2 auth columns.
2. Build the image (`docker build -t mbvlabs .`). Keep the same public
   hostname, pepper, token signing key, database, and integration
   credentials. `HOST_PRIMARY` replaces `DOMAIN` (the old name is accepted
   as a fallback). Existing users keep their passwords, but must sign in again
   because V2 uses a different session format.
3. Stop the V1 web and worker processes before migrating. Run `run-migrate`
   once using the new image and the existing database settings. Do not run
   the development seeder. To roll back, stop V2 and restore the rehearsed
   pre-migration snapshot before restarting V1.
4. Start `run-app`, `run-queue`, and `run-ssr` as separate processes from
   the same image. Set `ENVIRONMENT=production` and `PROTOCOL=https` for
   the public origin. Keep SSR private: if it runs in another container on
   the same private network, set `INERTIA_SSR_LISTEN=http://0.0.0.0:13714`
   for `run-ssr` and `INERTIA_SSR_URL=http://<ssr-service>:13714` for
   `run-app`; do not publish the SSR port. On a shared host/network namespace
   the default localhost address works for both processes.
5. Check `/api/ready`, the public pages, sign-in, contact, and queue activity
   before switching traffic.

Amazon SES remains the production email default. Supply the existing
`AWS_REGION`, `AWS_SES_ACCESS_KEY_ID`, `AWS_SES_SECRET_ACCESS_KEY`,
`AWS_SES_CONFIGURATION_SET`, and `DEFAULT_SENDER_SIGNATURE` as appropriate.
The article, diary, and search integrations use `API_BASIC_AUTH_USERNAME`
and `API_BASIC_AUTH_PASSWORD`; search requests also need `SERPER_API_KEY`
and `FIRECRAWL_API_KEY`. Keep all credentials in deployment configuration.

The migration chain targets an empty database or the real V1 schema, not
a database already migrated by the discarded, empty V2 scaffold.
