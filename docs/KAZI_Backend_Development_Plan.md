# KAZI Backend Development Plan

**Version:** 1.0 — Two-developer Go backend plan, Hackathon MVP  
**Developer 1:** Core Backend & Platform — teammate  
**Developer 2:** Submissions, Assessment & AI — you  
**Product:** KAZI — Prove what you can do.  
**Principle:** Let the work speak first.

## 1. Purpose and source of truth

Build the backend together before starting substantial frontend implementation. Developer 1 owns the platform around the assessment; Developer 2 owns the submission-to-evidence pipeline, including all AI integration.

This document consolidates the two-developer plan from **“Branch · Write winning hackathon idea”** and checks it against the uploaded **KAZI API Contract Specification v0.2 (Draft, Hackathon MVP)**, originally provided as `/mnt/data/KAZI_API_Contract_v2.md`. References such as “Contract §5.2” refer to sections in that file.

The contract is authoritative for existing routes, payloads, enums, authorization and response codes. Project structure, database tables, internal worker design and operational practices below are implementation recommendations from the plan, not additional contract requirements. Draft-contract gaps are listed in Section 11; agree and document their resolution before implementing dependent behavior.

### MVP scope

- One active domain: `SOFTWARE_ENGINEERING`.
- One complete practical challenge journey: choose, do, explain, answer a fixed follow-up, receive evidence, build a passport, match and connect.
- Go backend with PostgreSQL and REST endpoints under `/api/v1`.
- One candidate execution language initially: JavaScript/Node.js. Go is the backend language; it need not be the candidate challenge language.
- Seeded challenges and opportunities. No employer job-creation API or full job marketplace.
- Authenticated employer access to anonymous candidate evidence.
- Candidate-controlled identity reveal for each opportunity.

**AI contributes qualitative evidence but does not directly set the final score or integrity status. KAZI's assessment engine owns those decisions.** (Contract §6.1.)

## 2. Architecture and backend stack

| Component | Plan |
|---|---|
| Backend language | Go |
| API | REST, `/api/v1`, JSON |
| Database | PostgreSQL |
| Authentication | JWT access tokens plus refresh tokens; separate candidate/employer registration |
| Application structure | Modular monolith in one repository |
| Async processing | Durable jobs consumed by a Go worker |
| Candidate code execution | Isolated runner outside the main application server |
| AI integration | Server-side LLM API calls orchestrated from Go |
| Demo data | Seeded challenge, rubric, tests, employers and opportunities |

```text
Future frontend / API client
            |
            v
Go REST API
  Developer 1: Auth, Candidate, Employer, Challenge,
               Opportunity, Matching, Discovery, Invite, Identity
  Developer 2: Submission, Follow-up, Assessment, Evidence, Passport
            |
            v
PostgreSQL: application records + durable jobs
            |
            v
Go worker (Developer 2)
  Run submission -> isolated execution service -> trusted test results
  Assess completed follow-up -> AI qualitative analysis
                            -> fixed rubric and integrity rules
                            -> assessment + evidence + passport
                                                     |
                                                     v
                              Developer 1: matching and employer view
```

Use one application codebase with clear internal modules. A separate worker process and isolated runner do not require splitting the platform into microservices. Candidate code must never run inside the API process or on the main application server (Contract §5.1). The AI integration can remain in Go; calling an LLM does not require a Python service.

## 3. Recommended Go project structure

```text
kazi-backend/
  cmd/
    api/main.go
    worker/main.go
  internal/
    auth/                  # Developer 1
    candidate/             # Developer 1
    employer/              # Developer 1
    challenge/             # Developer 1; challenge content agreed together
    opportunity/           # Developer 1
    matching/              # Developer 1; consumes passport evidence
    discovery/             # Developer 1; anonymous employer projection
    invite/                # Developer 1
    identity/              # Developer 1
    submission/            # Developer 2
    runner/                # Developer 2; isolated execution client
    followup/              # Developer 2
    ai/                    # Developer 2; provider, prompts, validation
    assessment/            # Developer 2; rubric and scoring rules
    evidence/              # Developer 2
    integrity/             # Developer 2
    passport/              # Developer 2
    jobs/                  # Developer 2; durable async processing
    database/              # Developer 1 leads shared infrastructure
    middleware/            # Developer 1
    config/                # Developer 1; both add module settings
    response/              # Developer 1; shared API errors
  migrations/              # Shared; coordinated numbering/order
  seeds/                   # Shared demo fixtures
  sandbox/                 # Developer 2; runner image and trusted harness
  tests/integration/       # Shared API journey and boundary checks
  docs/
    KAZI_API_Contract_v2.md
    assessment-rules.md
    integration-decisions.md
  .env.example
  go.mod
  go.sum
  README.md
```

Within a domain, keep HTTP handlers, business logic, repository operations and domain types together, for example `handler.go`, `service.go`, `repository.go` and `types.go`. Share only the types needed at module boundaries; avoid turning a global `models` package into a dependency for everything. Add `pkg/` only for code that actually needs reuse outside the application.

## 4. Developer 1 — Core Backend & Platform

### 4.1 Foundation

Own Go setup, configuration, database connection, migration tooling, routing, middleware, validation, structured logs, error responses and local startup documentation. Coordinate shared identifiers, timestamps, API conventions and configuration names with Developer 2.

Provide authentication context that downstream modules can trust: user ID, role and resolved candidate/employer ownership. Submission ownership must come from this context, not a candidate ID supplied in the request body.

### 4.2 Authentication

Implement separate candidate and employer registration, one shared login endpoint, verification codes, password hashing, email uniqueness, access-token generation, refresh-token lifecycle, logout and role authorization (Contract §1).

- Candidate registration collects email/password without requiring a name.
- Employer registration collects company identity; enforce the contract's mandatory fields without treating every example field as required.
- Login rejects unverified and suspended accounts.
- Logout invalidates both access and refresh tokens; design revocation/session checking accordingly.
- Preserve the contract's response shapes, including `expires_in` values of 600 for registration verification and 900 for access tokens.

### 4.3 Candidate profiles

Create one candidate profile per user and a unique anonymous identifier such as `KZ-1042`. Implement self-profile reads/updates and invitation list/respond operations (Contract §2).

Names and country are optional for completing challenges and receiving a passport. The candidate's own view may contain their identity; that object must never be reused as an employer response. Duplicate profile creation returns `409`.

### 4.4 Employer profiles

Implement employer self-profile reads and permitted updates, including company/contact information and profile completion fields (Contract §3). Store verification status using `PENDING`, `VERIFIED` or `UNVERIFIED`. Formal employer verification is not an authorization requirement in this MVP.

### 4.5 Challenges

Implement public challenge list/detail APIs, category filtering and seeded records (Contract §4). Store title, description, difficulty, `repo_url`, category, rubric and `followup_question`.

Developer 2 supplies the runnable challenge fixture, trusted tests and assessment requirements. Agree on challenge/test/rubric versions together so an existing submission is not silently evaluated against changed material.

### 4.6 Opportunities and matching

Seed opportunities and their requirements; implement public listing and authenticated candidate self-match (Contract §8). Consume Developer 2's completed passport/evidence, calculate rule-based matches and return `match_score`, `matched_requirements` and `gaps`. Return `404` when the opportunity or candidate passport does not exist.

Agree on canonical skill names, evidence weights and the definition of a matched requirement. Do not simply copy the assessment score into `match_score`; the two measure different things. See Section 11 for the draft formula ambiguity.

### 4.7 Employer candidate discovery

Own the `/employer/*` anonymity boundary (Contract §9). Every request requires an authenticated `EMPLOYER`. The MVP contract defines candidate evidence lookup by candidate ID, not a candidate search/list API; use a seeded demo candidate reference rather than silently adding browsing routes.

Select only employer-safe fields at the query/projection layer. Initially return anonymous ID, category, score, skills, evidence summary and `identity_revealed: false`. Do not return name, email, phone, school or raw candidate data. After authorized reveal, the contract adds **name and country only**, not all private fields.

### 4.8 Invites and identity reveal

Create invites tied to candidate/opportunity pairs; list them with job/company/location context; allow the recipient to accept or decline (Contract §§2.4–2.5, 9.2).

- Invite statuses: `PENDING`, `ACCEPTED`, `DECLINED`.
- Duplicate invitations and repeated responses return `409` as specified.
- Sending or accepting an invite does not reveal identity.
- Only the candidate can call the explicit reveal action.
- Reveal requires the candidate's `name` and `country` and applies to one opportunity (Contract §10.1).
- Store reveal grants by candidate/opportunity, never as a global authorization boolean.
- Verify the employer's relationship to the opportunity before exposing identity. Revealing for one opportunity must not expose identity for another.

## 5. Developer 2 — Submissions, Assessment & AI

**Developer 2 is you. You own all AI integration and the full assessment pipeline, including its HTTP endpoints.** Developer 1 supports shared infrastructure and route wiring, while final module responsibility remains with you.

### 5.1 Submission engine

Implement submission creation, owner-only polling and workflow persistence (Contract §§5.1–5.2).

```json
{
  "challenge_id": "019171b3-4fa4-711b-8e7c-a49688bc8a96",
  "code": "string",
  "explanation": "string"
}
```

Return `202 Accepted` with `submission_id` and `status: QUEUED`. Validate required fields and challenge existence, and enforce one active submission per candidate/challenge with a database-backed concurrency constraint. Store the submission and its work item reliably before acknowledging it.

```text
QUEUED -> RUNNING -> COMPLETED
                 -> FAILED
```

`COMPLETED` here means execution has completed: Contract §5.2 explicitly shows a completed submission with `followup_status: PENDING`. It does not mean the final assessment or passport already exists. Keep assessment processing state separate internally. Failed candidate assertions can still produce a completed execution with failed tests; distinguish that from runner/infrastructure failure.

### 5.2 Isolated code runner

Own dispatch, execution limits, result collection and cleanup:

1. Load the submission and pinned challenge/test version.
2. Create a fresh isolated environment outside the application server.
3. Supply the candidate code and trusted test harness.
4. Run under bounded CPU, memory, process, output and wall-clock limits.
5. Collect structured results, execution time and diagnostic output.
6. Destroy the environment and persist execution outcome.

Do not provide database credentials, AI credentials, host filesystem mounts or container-control sockets to candidate code. Disable outbound network access unless an explicitly designed challenge requires it. Use a non-privileged execution environment. Treat candidate stdout as untrusted; it must not be able to forge the trusted test result record. A container name alone is not an isolation design.

Support JavaScript/Node.js first. Prebuild dependencies and keep the harness reproducible so execution does not depend on arbitrary installation during each run.

### 5.3 Automated tests and result processing

Prepare challenge-specific tests covering the intentional bug, expected behavior and edge cases. Parse results into the contract shape:

```json
{
  "passed": 8,
  "failed": 2,
  "total": 10,
  "execution_time": 1.82
}
```

Retain internal per-test evidence so scoring can distinguish competencies. Agree on execution-time units and precision. Protect hidden tests from public challenge responses and unnecessary AI exposure.

### 5.4 Fixed follow-up

Implement `POST /submissions/{submission_id}/followup` (Contract §5.3). Use the fixed question from the challenge, storing its version/snapshot, answer, submission ID and timestamps.

Accept an answer only from the owner, after submission execution is `COMPLETED`, and only once. Invalid state or a repeated answer returns `409`. Store `followup_status: COMPLETED` and enqueue assessment work. Do not add adaptive or AI-generated follow-up questions for this MVP.

### 5.5 AI, assessment, evidence and integrity

Own provider integration, prompt versioning, response validation, objective-evidence processing, scoring rules, evidence explanations, integrity rules and owner-only assessment reads. The detailed pipeline and guardrails are in Section 8.

### 5.6 Skills Passport

Own passport generation and `GET /passport/me` (Contract §7.1). Populate:

- `candidate_id`, `anonymous_id`, category and `overall_score`.
- Skills with contract-approved evidence statuses.
- `evidence_summary`: `challenges_completed`, `tests_passed`, `explanation_submitted`, `followup_completed`.

Return `404` when no completed assessment exists. Publish only completed, consistent assessment data. For the initial one-challenge demo, use that completed assessment; document the selection/aggregation policy before supporting multiple attempts or challenges. Provide Developer 1 a stable internal evidence/passport read interface without granting employer access to candidate-only endpoints.

## 6. Complete API ownership

All paths below are relative to **`/api/v1`**. Public endpoints need no bearer token; all other rows require the appropriate authenticated role and ownership checks.

| Contract | Method and route | Owner | Access / main behavior |
|---|---|---|---|
| §1.1 | `POST /auth/register/candidate` | Developer 1 | Public; separate candidate registration, `201` |
| §1.2 | `POST /auth/register/employer` | Developer 1 | Public; separate employer registration, `201` |
| §1.3 | `POST /auth/verify-email` | Developer 1 | Public; validate code, `200` |
| §1.4 | `POST /auth/login` | Developer 1 | Public; return role and tokens, `200` |
| §1.5 | `POST /auth/refresh` | Developer 1 | Public endpoint; valid refresh token required, `200` |
| §1.6 | `POST /auth/logout` | Developer 1 | Authenticated user; revoke session, `204` |
| §2.1 | `POST /candidates` | Developer 1 | Candidate; create own profile, `201` |
| §2.2 | `GET /candidates/me` | Developer 1 | Candidate owner, `200` |
| §2.3 | `PATCH /candidates/me` | Developer 1 | Candidate owner, `200` |
| §2.4 | `GET /candidates/me/invites` | Developer 1 | Candidate owner; optional status filter, `200` |
| §2.5 | `PATCH /candidates/me/invites/{invite_id}` | Developer 1 | Invite recipient; accept/decline, `200` |
| §3.1 | `GET /employers/me` | Developer 1 | Employer owner, `200` |
| §3.2 | `PATCH /employers/me` | Developer 1 | Employer owner, `200` |
| §4.1 | `GET /challenges` | Developer 1 | Public; optional category filter, `200` |
| §4.2 | `GET /challenges/{challenge_id}` | Developer 1 | Public, `200` |
| §5.1 | `POST /submissions` | Developer 2 | Candidate; enqueue execution, `202` |
| §5.2 | `GET /submissions/{submission_id}` | Developer 2 | Submission owner; polling, `200` |
| §5.3 | `POST /submissions/{submission_id}/followup` | Developer 2 | Submission owner; fixed answer, `200` |
| §6.1 | `GET /assessments/{assessment_id}` | Developer 2 | Related candidate owner, `200` |
| §7.1 | `GET /passport/me` | Developer 2 | Candidate owner, `200` |
| §8.1 | `GET /opportunities` | Developer 1 | Public seeded list, `200` |
| §8.2 | `GET /opportunities/{opportunity_id}/match` | Developer 1 | Candidate self-match, `200` |
| §9.1 | `GET /employer/candidates/{candidate_id}` | Developer 1 | Employer; anonymous/scoped revealed view, `200` |
| §9.2 | `POST /employer/candidates/{candidate_id}/invite` | Developer 1 | Employer; opportunity-specific invite, `201` |
| §10.1 | `POST /identity/reveal` | Developer 1 | Candidate owner; opportunity-specific grant, `200` |

Preserve the distinction between plural `/employers/me` and singular `/employer/candidates/...`. There is no public assessment-creation endpoint in v0.2; the worker creates assessments internally.

Use the contract's per-route errors. Shared conventions are `400` validation, `401` invalid/missing authentication, `403` forbidden, `404` unavailable resource, `409` conflict and `500` server error. Agree on an error-body schema because the draft lists errors without defining one common JSON envelope.

## 7. Shared database and integration responsibilities

### 7.1 Schema ownership

Both developers agree on the schema before parallel implementation. Developer 1 leads database setup and migration coordination. Each developer authors and maintains the tables/repositories for their domain; Developer 2 does not wait for Developer 1 to build all assessment storage.

| Suggested tables | Primary owner | Responsibility |
|---|---|---|
| `users`, `refresh_tokens`, `verification_codes` | Developer 1 | Authentication identity and sessions |
| `candidate_profiles`, `employer_profiles` | Developer 1 | Separate role-specific profiles |
| `challenges` | Developer 1, shared design | Metadata, rubric and fixed question versions |
| `submissions`, `test_results`, `followups` | Developer 2 | Work, execution evidence and explanation follow-up |
| `assessment_jobs` or shared `jobs` | Developer 2 | Durable queue, retries and worker state |
| `assessments`, `assessment_evidence` | Developer 2 | Final scoring and traceable evidence |
| `ai_evaluations`, `integrity_signals` | Developer 2 | Internal analysis, provenance and rule inputs |
| `skills`, `candidate_skills` | Developer 2, shared vocabulary | Passport evidence and matching inputs |
| `opportunities`, `opportunity_requirements` | Developer 1 | Seeded employer-linked opportunity data |
| `matches` | Developer 1 | Optional cached/recorded match output |
| `invites`, `identity_reveals` | Developer 1 | Invitations and candidate/opportunity grants |

A shared authentication `users` table is compatible with separate candidate and employer entities. Keep their profile fields in distinct models/tables. A passport may be a projection over assessments and skills rather than a separate table; match storage is optional if calculated on read.

### 7.2 Rules to agree together

- UUIDv7 identifiers where specified by the contract; anonymous IDs are a separate display identifier. Use UTC timestamps.
- Foreign keys and ownership relationships: user to profile, candidate/challenge to submission, submission to assessment, employer to opportunity, candidate/opportunity to invite and reveal.
- Uniqueness for profile per user, anonymous ID, active submission per candidate/challenge, follow-up per submission, invite per candidate/opportunity and reveal per candidate/opportunity.
- Idempotent worker processing: retries must not duplicate assessments, evidence or passport updates.
- Atomic state transitions: persist assessment, evidence and updated passport consistently; commit follow-up and its job together.
- Keep private identity out of assessment/AI input and employer projections. Restrict diagnostic logs containing code or candidate text.
- Keep migrations additive and coordinated. Do not edit a migration already applied by the teammate; add a new migration.

### 7.3 Internal handoffs

| Producer -> consumer | Shared interface/data |
|---|---|
| Developer 1 auth -> Developer 2 submission | Authenticated candidate identity and ownership helpers |
| Developer 1 challenge -> Developer 2 runner | Challenge ID/version, rubric, fixed question; private runner fixture reference |
| Developer 2 runner -> assessment | Trusted per-test results, timings and execution outcome |
| Developer 2 follow-up -> assessment | Stored answer and completed submission ID |
| Developer 2 passport -> Developer 1 matching/discovery | Completed score, skills, statuses, summary and source assessment version |
| Developer 1 identity -> employer projection | Authorized candidate/opportunity reveal grant |

Suggested internal entry points are `RunSubmission(submissionID)`, `AssessSubmission(submissionID)` and `GetPassport(candidateID)`. These are implementation interfaces, not new HTTP routes. `AssessSubmission` requires completed execution and follow-up, reads stored test results, validates AI evidence, applies the fixed rubric, creates evidence and updates the passport.

## 8. AI pipeline, scoring and guardrails

### 8.1 Inputs and processing

```text
Stored candidate code + explanation
            |
Isolated execution -> objective test evidence
            |
Fixed follow-up answer completed
            |
Challenge description + expected competencies + fixed rubric
            |
AI qualitative evidence analysis
            |
Validate structured output and supporting references
            |
KAZI scoring rules + KAZI integrity rules
            |
Assessment -> evidence -> Skills Passport
```

Send the AI the challenge, competencies, rubric, submitted code, relevant automated test results, explanation, fixed question and answer. Exclude candidate name, email, country and employer identity. Treat code, comments and free text as untrusted material to analyze, never as instructions that may override the evaluator.

Ask for structured, supported observations. The following is a suggested **internal** AI response, not a v0.2 public response:

```json
{
  "observations": [
    {
      "competency": "DEBUGGING",
      "finding": "The explanation identifies the ID lookup failure addressed by the code.",
      "supporting_references": ["code:lookup-handler", "explanation:paragraph-1"]
    }
  ],
  "consistency_findings": [],
  "feedback": "The core lookup bug is addressed; edge-case handling needs more work."
}
```

Actual references must resolve to submitted material or trusted test evidence. Validate allowed competencies, field types, lengths and reference existence. The earlier conversation's suggested qualitative competency labels can be used as internal suggestions, but the backend must validate/map them before publishing evidence statuses.

### 8.2 AI must not decide

- Final numeric `score` or final `integrity_status`.
- Candidate rankings or hiring/rejection decisions.
- Cheating probability, fraud probability or AI-generated-code probability.
- New rubric weights, invented test outcomes or adaptive follow-up questions.

Never copy a model-supplied overall consistency label straight into `integrity_status`. Reject unsupported fields/claims, validate feedback and keep the provider response separate from authoritative assessment records.

### 8.3 Fixed rubric

The contract's challenge example defines these weights (Contract §4.2):

| Competency | Weight | Illustrative contribution |
|---|---:|---:|
| `CORRECT_IMPLEMENTATION` | 30% | 27/30 |
| `TESTS_EDGE_CASES` | 20% | 15/20 |
| `DEBUGGING` | 20% | 17/20 |
| `CODE_QUALITY` | 15% | 11/15 |
| `EXPLANATION` | 15% | 12/15 |
| **Total** | **100%** | **82/100** |

The contributions are the conversation's example, not values implied by 8/10 tests passing. Define each competency's evidence checks, normalized score, caps, rounding and treatment of missing evidence in `assessment-rules.md` before implementation.

Suggested calculation: `score = round(100 * sum(weight_i * competency_score_i))`, where each competency score is between 0 and 1 and is produced by documented backend rules. Objective evidence and validated qualitative observations feed those rules. The weights alone are not a complete scoring specification.

Persist rubric, test, prompt/model and assessment-rule versions so a published result can be explained and reproduced from its recorded inputs.

### 8.4 Evidence engine

Publish only the contract's evidence statuses:

```text
STRONG_EVIDENCE
DEMONSTRATED
DEVELOPING
NEEDS_IMPROVEMENT
```

Each assessment evidence entry has a competency, status and explanation (Contract §6.1). Map rubric competencies to skill labels consistently; for example, “API Development” and “REST APIs” need a shared mapping for matching.

Evidence should say what was observed, such as “2 of 10 tests failed, primarily around edge cases,” with internal references to the relevant tests. Do not publish unsupported model assertions as established facts.

### 8.5 Integrity engine

Allowed final statuses are `HIGH` and `REVIEW_RECOMMENDED`. Compare code, objective test evidence, explanation and follow-up using documented backend rules. Supported consistency and sufficient evidence may produce `HIGH`; substantiated discrepancies or insufficient evidence under the agreed rules may produce `REVIEW_RECOMMENDED`.

These statuses describe assessment evidence and review needs. They do not prove whether someone used AI or cheated. Low skill/test performance alone is not evidence of dishonesty, and provider failure must not be treated as a candidate integrity violation.

### 8.6 Failure handling

Set AI timeouts and bounded retries; validate structured responses before use. Keep valid execution results when AI processing fails. Track assessment work/failure separately from the completed submission and do not fabricate a finished passport.

During development, a clearly marked deterministic AI stub can unblock module integration. The end-to-end AI milestone requires real provider validation. Any objective-only fallback and public indication of delayed/failed assessment must be agreed and documented; v0.2 does not define those details.

## 9. Git workflow

Preserve the branch plan agreed in the conversation:

```text
main                    Stable demo/release
develop                 Shared integration
dev-core-backend        Developer 1
dev-assessment-ai       Developer 2

feature branch -> developer branch -> develop -> main
```

Example feature branches: `feature/auth`, `feature/candidate-profile`, `feature/employer`, `feature/submissions`, `feature/code-runner`, `feature/assessment`, `feature/ai-evaluation`, `feature/passport`.

1. Start small feature branches from the appropriate developer branch.
2. Keep changes focused on the owned module; coordinate shared wiring and migrations.
3. Submit reviewable pull requests with behavior, contract changes and relevant validation.
4. Merge into the developer branch and integrate into `develop` frequently to avoid drift.
5. Both review changes to authentication, identity boundaries, public contracts and shared schema.
6. Promote a tested integration to `main` for the demo.

Keep secrets out of Git; commit placeholders in `.env.example`. A change to routes, fields or enums follows: discuss -> update contract -> implement -> update API examples/tests -> update frontend when it exists. Run formatting, Go tests and applicable integration checks before merging.

## 10. Development stages and order

### Stage 1 — Shared foundation

Agree on repository layout, environment configuration, database/migrations, IDs, authentication context, error JSON, challenge fixture, test harness, rubric checks and internal interfaces. Resolve the contract gaps in Section 11. Developer 1 creates the API/database skeleton while Developer 2 prepares the runner fixture and assessment rules.

**Exit:** Both can start the API/worker locally, apply the same migrations and access identical seeded challenge data.

### Stage 2 — Implement owned modules in parallel

| Order | Developer 1 | Developer 2 |
|---|---|---|
| 1 | Authentication and role enforcement | Submission persistence and durable execution job |
| 2 | Candidate/employer profiles | Isolated runner and trusted result parsing |
| 3 | Challenge list/detail and seeds | Submission polling and fixed follow-up |
| 4 | Opportunity seeds and matching against a fixture passport | AI adapter and structured response validation |
| 5 | Anonymous employer view and invites | Assessment rules, integrity rules and evidence |
| 6 | Candidate invite responses and scoped reveal | Passport generation and assessment/passport reads |

Use agreed fixtures to unblock each other, then replace them with real module outputs during integration.

**Exit:** Each developer can demonstrate their module through API requests and relevant automated checks.

### Stage 3 — Integrate the complete backend

```text
Candidate registration -> email verification -> login -> create profile
  -> choose challenge -> submit code + explanation
  -> QUEUED -> RUNNING -> COMPLETED execution
  -> answer fixed follow-up -> AI qualitative evidence
  -> fixed rubric + integrity rules -> assessment + evidence
  -> Skills Passport -> opportunity match

Employer registration -> email verification -> login
  -> anonymous candidate evidence -> invite for opportunity
  -> candidate sees context -> accepts/declines
  -> candidate supplies name/country and explicitly reveals
  -> employer sees name/country only in authorized opportunity context
```

Use Postman, Bruno or an equivalent API collection. No manual database edits should be required after initial seeds. Invite acceptance and identity reveal remain separate requests.

**Exit:** Both roles complete the full flow, including anonymity before reveal and isolation between opportunities after reveal.

### Stage 4 — Frontend and demo readiness

Start substantial frontend work once Stage 3 passes. Give the frontend stable payload examples, the approved contract, polling behavior, loading/error states and the API collection. Poll submission status first; after execution, collect follow-up, then poll for a completed passport (the draft suggests roughly two-second intervals).

**Exit:** The demo runs through ordinary API calls with real test execution and AI evidence, plus clear loading, empty and failure states.

## 11. Draft-contract decisions to resolve together

These are identified gaps in v0.2, not silent changes to its API.

| Issue | Why it matters | Proposed resolution to record |
|---|---|---|
| Assessment ID discovery | §6.1 needs `assessment_id`, but submission/follow-up/passport responses do not supply it. | Agree an additive nullable `assessment_id` on submission polling, or another documented lookup mechanism. Passport polling works for the passport but does not expose assessment details. |
| Assessment readiness/failure | A submission can be `COMPLETED` before follow-up; submission status cannot express later AI/assessment failure. | Define separate assessment processing visibility and retry/error behavior without changing the meaning of execution status. |
| Employer opportunity context | §9.1 promises per-opportunity reveal but its URL/request carries only `candidate_id`. Multiple opportunities for the same employer are ambiguous. | Prefer an explicit opportunity context, validated against the employer, through an agreed contract addition. Keep the response anonymous when context cannot be established; never use “any reveal exists” to disclose identity. |
| Candidate profile reveal flag | §2 includes `identity_revealed` without opportunity context. | Define it as display-only metadata or revise its shape; never use it to authorize employer disclosure. |
| Matching weights | §8.2 gives a count-based formula “weighted against evidence status,” but no weights. Its example is 82 for three matched requirements out of four, so 82 is not derived by the plain count formula. | Agree evidence weights, thresholds, denominator and rounding, then correct examples. Do not hard-code 82. |
| Detailed rubric scoring | §4.2 supplies weights but not per-competency scoring rules. | Both approve explicit evidence-to-score rules and fixtures; Developer 2 implements them. |
| Opportunity ownership and reveal prerequisites | Employer/opportunity relationships are needed internally; §10.1 does not explicitly require an accepted invite. | Model opportunity ownership and document checks. Preserve the separate candidate reveal action; any accepted-invite requirement is a contract amendment, not an existing rule. |
| Cross-module naming and errors | Skills differ in display names; one common error JSON body and timing units are unspecified. | Agree canonical skill mappings, error envelope, runtime units and API examples before frontend integration. |

## 12. Priorities

### P0 — Must work for the backend MVP

- Go/API/database setup and reproducible seeds.
- Candidate and employer registration, verification flow, login, refresh/logout and role/ownership checks.
- Candidate and employer profiles; public challenge endpoints.
- Submission queue, isolated execution, trusted tests, polling and fixed follow-up.
- Real AI qualitative analysis, backend scoring/integrity rules, evidence and passport.
- Opportunity listing, rule-based matching and authenticated anonymous employer evidence.
- Invites, candidate responses and explicit opportunity-scoped identity reveal.
- Essential logs, bounded execution/jobs, retry safety, privacy checks and the end-to-end API demo.
- Resolution of contract gaps needed for that complete flow.

Email verification behavior remains P0 because the contract rejects unverified login. For local development, use a documented development delivery sink; do not quietly remove verification. Public/live delivery must work before claiming the full registration experience works.

### P1 — Improvements after the core flow

- Production email-delivery polish and monitoring beyond the development/demo setup.
- Richer operational dashboards and logging; broader rate-limit tuning beyond essential abuse/resource controls.
- More challenges using the same execution language and pipeline.
- Better assessment analytics and feedback presentation.
- Employer verification workflow exploration. Keep authorization unchanged unless the contract is revised; formal verification is roadmap work in v0.2.

### P2 — After the hackathon

- Multiple programming languages or skill domains.
- Adaptive follow-ups and AI-generated challenges.
- University integrations.
- Employer job creation, full marketplace or employer CRM.
- Payments and subscriptions.
- Advanced integrity investigations and workforce analytics, with appropriate evidence standards.

## 13. Integration verification and completion checklist

Both developers own the integrated result. Each writes tests for their business rules and repositories, then shares boundary tests with the other.

- [ ] Both registration flows and verification work; incorrect roles and suspended/unverified logins are rejected.
- [ ] Logout invalidates both token types as contracted.
- [ ] Candidate profile creation is unique; optional identity is not required for assessments.
- [ ] Submission creation returns `202` promptly; concurrent duplicate active submissions are rejected.
- [ ] Code runs only in the isolated environment; timeouts and cleanup work.
- [ ] Test failures are reported accurately and cannot be spoofed through candidate output.
- [ ] Submission polling is owner-only and distinguishes execution from later assessment readiness.
- [ ] Follow-up before completion and repeated answers return `409`.
- [ ] AI malformed output, unsupported claims and prompt-injection text do not bypass validation or scoring rules.
- [ ] Scoring fixtures reproduce expected totals; the AI cannot set final score/integrity fields.
- [ ] Worker retries do not duplicate results; provider failures do not fabricate candidate evidence.
- [ ] Passport and assessment reads enforce candidate ownership.
- [ ] Matching uses documented skill/evidence mappings and produces explainable matches/gaps.
- [ ] Employer discovery rejects unauthenticated/candidate callers and omits private identity before reveal.
- [ ] Invite creation/acceptance never implicitly reveals identity.
- [ ] Reveal requires candidate action and name/country; duplicate reveal returns `409`.
- [ ] A second employer and another opportunity remain anonymous after the first reveal, including another opportunity for the same employer.
- [ ] Both roles complete the API journey without manual database edits.

## 14. Ownership summary

| Area | Developer 1 | Developer 2 — You |
|---|---|---|
| Go setup, routing, configuration | Lead | Support |
| PostgreSQL and schema | Shared; infrastructure lead | Shared; assessment schema lead |
| Authentication and sessions | Owner | Consume shared context |
| Candidate profile | Owner | — |
| Employer profile | Owner | — |
| Challenges API and metadata | Owner | Test/rubric fixture support |
| Submissions API and async lifecycle | Infrastructure support | Owner |
| Isolated code runner | — | Owner |
| Automated tests/result processing | — | Owner |
| Fixed follow-up | — | Owner |
| AI integration | — | Owner |
| Assessment engine and API | — | Owner |
| Evidence engine | — | Owner |
| Integrity signals and rules | Joint policy review | Owner |
| Skills Passport generation and API | Integration support | Owner |
| Opportunities | Owner | — |
| Matching | Owner | Passport/evidence support |
| Employer candidate discovery | Owner | Safe evidence interface support |
| Invites and responses | Owner | — |
| Identity reveal | Owner | Boundary review support |
| Shared contract and migrations | Shared | Shared |
| Security review and integration tests | Shared | Shared |
| Backend demo readiness | Shared | Shared |

**Core delivery sequence:** Challenge -> Submission -> Objective Tests -> Fixed Follow-Up -> AI Qualitative Evidence -> Fixed Rubric -> Assessment -> Skills Passport -> Opportunity Match -> Anonymous Employer View -> Invite -> Candidate-Controlled Reveal.
