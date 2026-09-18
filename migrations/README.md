# SQL migrations

Developer 1 coordinates migration ordering; both developers own their domain tables.
Use numbered SQL migrations with matching up/down files. The repository uses
`golang-migrate`; do not create tables at application startup or rewrite a migration
after it has been applied.

The foundation does not create business tables. Agree the schema in the development
plan before adding the first migration. Run migrations with a direct database URL;
the application itself uses a pooled `DATABASE_URL`.

Include opportunities and opportunity_requirements (missing from the earlier quick
table list). Keep identity_reveals scoped by candidate_id and opportunity_id.
Enforce unique active submissions, unique follow-ups and worker retry idempotency
in the database. A global candidate identity_revealed flag cannot authorize disclosure.
