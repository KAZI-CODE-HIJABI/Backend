# KAZI Backend

Go backend for **KAZI — Prove what you can do.**

## Current milestone: runnable foundation

Implemented: Gin API startup, validated environment configuration, Zap logging,
pgxpool/PostgreSQL connection pooling, graceful shutdown, liveness/readiness checks,
Docker Compose, unit tests and CI with PostgreSQL integration testing.

Business endpoints, authentication, schema migrations, jobs, candidate-code execution
and AI assessment are **not implemented yet**. Domain folders reserve ownership and
contain package documentation, not pretend-success handlers. The initial API returns
404 for unimplemented business routes.

## Start with Docker

Requires Docker with Compose. Copy `.env.example` to `.env`:

```powershell
Copy-Item .env.example .env
docker compose up --build -d
Invoke-RestMethod http://localhost:8080/healthz
Invoke-RestMethod http://localhost:8080/readyz
```

On macOS/Linux use `cp .env.example .env` and `curl` for the health requests.
`/healthz` returns 200 when the API is alive. `/readyz` returns 200 only when PostgreSQL
responds, otherwise 503 with no database error details. These infrastructure routes
are outside the business contract prefix `/api/v1`.

```sh
docker compose logs -f api
docker compose down
```

`down` preserves the local database volume. This Compose configuration is for local
development, with loopback-only published ports and local-only credentials. It is
not a production deployment configuration.

## Run Go locally

Install Go 1.26 or newer and start just the database:

```sh
docker compose up -d db
go mod download
go run ./cmd/api
```

The API loads `.env`; actual environment variables override it. The host database
address is `localhost:5432`; inside Compose it is `db:5432`. Do not run the host API
and Compose API on port 8080 at the same time.

```sh
go fmt ./...
go vet ./...
go test ./...
go build -o bin/api ./cmd/api
```

## Database migrations

Use the direct database URL in `DATABASE_MIGRATION_URL`. Then run:

```sh
go run ./cmd/migrate up
```

or `make migrate-up`. Migrations `000004` and `000005` belong to the assessment
branch and require the authentication/profile and challenge migrations planned as
`000001` and `000002`. Do not apply them until those prerequisite migrations exist.

Database integration tests skip unless `TEST_DATABASE_URL` is set. To include them
in PowerShell after starting the database:

```powershell
$env:TEST_DATABASE_URL = 'postgres://kazi:kazi_local_only@localhost:5432/kazi?sslmode=disable'
go test ./...
```

CI runs the integration test against a real PostgreSQL service and runs Go's race
detector. Windows race testing additionally requires a compatible C compiler.

## Team branches

- `main`: stable shared foundation/demo.
- `develop`: shared integration branch.
- `feature/core-backend`: Developer 1's starting branch.
- `feature/assessment-ai`: Developer 2's starting branch (AI owner).

Create focused feature branches from `develop` as work grows, submit pull requests
to `develop`, and synchronize frequently. Promote tested integration changes to
`main`. Never force-push over your teammate's work.

## Ownership and next work

See [the complete development plan](docs/KAZI_Backend_Development_Plan.md),
[the API contract](docs/KAZI_API_Contract_v2.md), and [the first milestones](docs/NEXT_STEPS.md).

Developer 1 owns core setup, auth, candidate/employer profiles, challenges,
opportunities, matching, discovery, invites and scoped identity reveal.

Developer 2 owns submissions, isolated runner, tests, fixed follow-up, AI,
assessment, evidence, integrity and Skills Passport.

Use Gin's validator-backed binding for request validation. Gin, pgxpool, Viper and
Zap are wired now. JWT, UUIDv7 and hot-reload tooling should be
chosen and pinned when their modules are introduced; unused packages are not added
just to fill out a stack list.

## Non-negotiable product boundaries

- Candidate code never runs on the application server. This Dockerfile runs only
  the API; it is not the candidate sandbox.
- AI supplies qualitative evidence. Backend rules set final score and integrity.
- Employers see anonymous projections until explicit per-opportunity reveal.
- Invite acceptance does not reveal identity.
- The challenge follow-up is fixed and required before final assessment.
- Backend-first milestone: register -> verify email -> login -> profile -> challenge
  -> submit -> isolated tests -> follow-up -> assessment -> passport.

## References

- [Go downloads](https://go.dev/dl/)
- [Gin quickstart](https://gin-gonic.com/en/docs/quickstart/)
- [Go REST API tutorial](https://go.dev/doc/tutorial/web-service-gin)
