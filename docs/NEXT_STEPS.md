# First development milestones

## Together: agree before creating business tables

Review Section 11 of the development plan: assessment ID discovery, assessment
readiness, opportunity context in employer discovery, matching evidence weights,
scoring rules, ownership relations and the common error envelope.

Adopt the branch names in this repository README; these supersede the longer
developer-branch workflow in the earlier planning document for this initialization.

## Developer 1: first pull request

Add reviewed auth/profile SQL migrations and the migration runner; implement the
separate registrations, verification, login, refresh and logout. Use secure password
hashing, UUIDv7 IDs, short-lived access tokens and session revocation. Authentication
context must carry server-verified user/role ownership. Logout must invalidate both
access and refresh tokens as the contract requires. Then implement profile creation
and seeded challenge reads. Test cross-role and cross-candidate access denial.

## Developer 2: first pull request

Agree the challenge fixture, private tests, fixed rubric checks and follow-up with
Developer 1. Add submission/test/follow-up/job migrations, then submission persistence
and owner-only polling. Build a durable worker and isolated runner adapter before
enabling candidate code execution. Keep completed execution distinct from completed
assessment. Add trusted result parsing, retry safety and bounded resource limits.

## Next integration

Add fixed follow-up -> structured AI observations -> validated evidence -> backend
rubric and integrity rules -> passport. Developer 1 consumes a safe passport/evidence
projection for matching and employer views, then implements invites and reveal.

The worker must not execute arbitrary code locally as a temporary shortcut. A stub
may return clearly labelled test fixtures in tests only; it must not produce a
real-looking assessment for an actual candidate.

## Assessment branch status

Completed: pgxpool connection layer, direct-URL migration runner, migrations 000004
and 000005, submission state rules, trusted test-result validation, isolated runner
interface, AI-observation validation, fixed assessment rules and passport projection.

Still needed: a real isolated Node.js execution implementation, persistence repositories
and HTTP handlers after authentication middleware and challenge migrations are merged,
an LLM provider credential and adapter, and the end-to-end candidate/employer flow.
