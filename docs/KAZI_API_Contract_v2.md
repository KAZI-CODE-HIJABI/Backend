# KAZI API Contract Specification

**Version:** `0.2` (Hackathon MVP)
**Base URL:** `/api/v1`
**Content-Type:** `application/json`
**Designed By:** Backend Team
**Status:** Draft

---

## Roles

| Role | Who | Registers via |
|---|---|---|
| `CANDIDATE` | The person proving skills | `1.1 Register Candidate` |
| `EMPLOYER` | The company hiring | `1.2 Register Employer` |

Registration, profile shape, and profile fields are **separate per role** — a candidate and an employer are different entities with different data, not the same `users` record with optional fields.

---

## 1. Authentication Domain (`/auth`)

### 1.1 Register Candidate

- **Method:** `POST`
- **URL:** `/auth/register/candidate`
- **Authentication:** None
- **Authorization:** Public
- **Request:**
  ```json
  {
    "email": "aisha@example.com",
    "password": "SecurePassword123!"
  }
  ```
- **Response (`201 Created`):**
  ```json
  {
    "message": "Verification code sent.",
    "expires_in": 600
  }
  ```
- **Note:** No name is collected at registration. Identity fields are added later, only if/when the candidate chooses (see `2.3 Update My Candidate Profile`).
- **Error Responses:**
  - `400 Bad Request` — Missing mandatory fields (`email`, `password`) or malformed input
  - `409 Conflict` — Email address already registered

---

### 1.2 Register Employer

- **Method:** `POST`
- **URL:** `/auth/register/employer`
- **Authentication:** None
- **Authorization:** Public
- **Request:**
  ```json
  {
    "email": "recruiter@acme.com",
    "password": "SecurePassword123!",
    "company_name": "Acme Ltd",
    "contact_first_name": "Almaz",
    "contact_last_name": "Kebede",
    "phone": "+251911556677",
    "location": "Addis Ababa, Ethiopia",
    "industry": "Software / Technology",
    "company_size": "11-50",
    "website_url": "https://acme.com"
  }
  ```
- **Response (`201 Created`):**
  ```json
  {
    "message": "Verification code sent.",
    "expires_in": 600
  }
  ```
- **Note:** Unlike candidate registration, employer registration collects company identity up front — there's no anonymity concept on the employer side. `location`, `industry`, and `company_size` are optional but recommended; they're not required to register, only to complete the profile (see `3.2`).
- **Error Responses:**
  - `400 Bad Request` — Missing mandatory fields (`email`, `password`, `company_name`) or malformed input
  - `409 Conflict` — Email address already registered

---

### 1.3 Verify Email

- **Method:** `POST`
- **URL:** `/auth/verify-email`
- **Authentication:** None
- **Authorization:** Public
- **Request:**
  ```json
  {
    "email": "aisha@example.com",
    "verification_code": "482913"
  }
  ```
- **Response (`200 OK`):**
  ```json
  { "message": "Email verified successfully" }
  ```
- **Error Responses:**
  - `400 Bad Request` — Invalid or expired verification code

---

### 1.4 Login

- **Method:** `POST`
- **URL:** `/auth/login`
- **Authentication:** None
- **Authorization:** Public
- **Request:**
  ```json
  {
    "email": "aisha@example.com",
    "password": "SecurePassword123!"
  }
  ```
- **Response (`200 OK`):**
  ```json
  {
    "access_token": "eyJhbGciOiJIUzI1NiIsIn...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsIn...",
    "token_type": "Bearer",
    "expires_in": 900,
    "user": {
      "id": "019171b3-4f9e-711b-8e7c-a49688bc8a91",
      "email": "aisha@example.com",
      "role": "CANDIDATE",
      "account_status": "ACTIVE"
    }
  }
  ```
- **Note:** One shared login endpoint for both roles — `role` in the response (`CANDIDATE` or `EMPLOYER`) tells the frontend which profile endpoint to call next (`2.2` vs `3.1`) and which UI to route into.
- **Error Responses:**
  - `400 Bad Request` — Missing email or password
  - `401 Unauthorized` — Invalid credentials, unverified email, or account `SUSPENDED`

---

### 1.5 Refresh Token

- **Method:** `POST`
- **URL:** `/auth/refresh`
- **Authentication:** None
- **Authorization:** Public
- **Request:**
  ```json
  { "refresh_token": "eyJhbGciOiJIUzI1NiIsIn..." }
  ```
- **Response (`200 OK`):**
  ```json
  {
    "access_token": "eyJhbGciOiJIUzI1NiIsIn...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsIn...",
    "token_type": "Bearer",
    "expires_in": 900
  }
  ```
- **Error Responses:**
  - `400 Bad Request` — Missing refresh token
  - `401 Unauthorized` — Token expired, invalidated, or revoked

---

### 1.6 Logout

- **Method:** `POST`
- **URL:** `/auth/logout`
- **Authentication:** `Bearer <access_token>`
- **Authorization:** Authenticated User (`CANDIDATE`, `EMPLOYER`)
- **Request:**
  ```json
  { "refresh_token": "eyJhbGciOiJIUzI1NiIsIn..." }
  ```
- **Response (`204 No Content`):** Empty body
- **Note:** Logging out invalidates both access and refresh tokens.
- **Error Responses:**
  - `401 Unauthorized` — Missing or expired session

---

## 2. Candidate Domain (`/candidates`)

### 2.1 Create Candidate Profile

- **Method:** `POST`
- **URL:** `/candidates`
- **Authentication:** `Bearer <access_token>`
- **Authorization:** Authenticated `CANDIDATE` (post-registration, pre-onboarding)
- **Request:** None
- **Response (`201 Created`):**
  ```json
  {
    "candidate_id": "019171b3-4fa0-711b-8e7c-a49688bc8a92",
    "anonymous_id": "KZ-1042",
    "identity_revealed": false,
    "created_at": "2026-08-21T10:00:00Z"
  }
  ```
- **Note:** This assigns the anonymous identifier shown throughout the candidate journey. It is called once, immediately after email verification, before any challenge is started.
- **Error Responses:**
  - `401 Unauthorized` — Token missing or expired
  - `403 Forbidden` — Caller is not a `CANDIDATE`
  - `409 Conflict` — Candidate profile already exists for this user

---

### 2.2 Get My Candidate Profile

- **Method:** `GET`
- **URL:** `/candidates/me`
- **Authentication:** `Bearer <access_token>`
- **Authorization:** Candidate (Owner)
- **Request:** None
- **Response (`200 OK`):**
  ```json
  {
    "candidate_id": "019171b3-4fa0-711b-8e7c-a49688bc8a92",
    "anonymous_id": "KZ-1042",
    "name": null,
    "email": "aisha@example.com",
    "country": null,
    "identity_revealed": false,
    "created_at": "2026-08-21T10:00:00Z",
    "updated_at": "2026-08-21T10:00:00Z"
  }
  ```
- **Note:** This is the candidate's own view and always includes real identity fields, regardless of `identity_revealed`. That flag only governs what the **employer-facing** endpoint (`9.1`) returns.
- **Error Responses:**
  - `401 Unauthorized` — Token missing or expired
  - `404 Not Found` — Candidate profile not found

---

### 2.3 Update My Candidate Profile

- **Method:** `PATCH`
- **URL:** `/candidates/me`
- **Authentication:** `Bearer <access_token>`
- **Authorization:** Candidate (Owner)
- **Request:**
  ```json
  {
    "name": "Aisha Mohammed",
    "country": "Nigeria"
  }
  ```
- **Response (`200 OK`):** Full updated candidate profile object (same shape as `2.2`)
- **Note:** Identity fields are collected lazily and are never required to complete a challenge or receive a passport — only to reveal identity to an employer.
- **Error Responses:**
  - `400 Bad Request` — Invalid field values
  - `401 Unauthorized` — Missing or expired token

---

### 2.4 List My Invites

- **Method:** `GET`
- **URL:** `/candidates/me/invites`
- **Authentication:** `Bearer <access_token>`
- **Authorization:** Candidate (Owner)
- **Query Parameters:**
  - `status` (string, optional: `PENDING`, `ACCEPTED`, `DECLINED`)
- **Response (`200 OK`):**
  ```json
  {
    "data": [
      {
        "invite_id": "019171b3-4fa8-711b-8e7c-a49688bc8a9a",
        "status": "PENDING",
        "opportunity": {
          "id": "019171b3-4fa7-711b-8e7c-a49688bc8a99",
          "title": "Junior Software Engineer",
          "company_name": "Acme Ltd",
          "location": "Addis Ababa / Remote"
        },
        "created_at": "2026-08-21T12:10:00Z"
      }
    ]
  }
  ```
- **Note:** This is what tells the candidate what they were actually invited to — job title, company, location — not just that an invite exists. Without this endpoint, the candidate has no way to make an informed decision.
- **Error Responses:**
  - `401 Unauthorized` — Token missing or expired

---

### 2.5 Respond to Invite

- **Method:** `PATCH`
- **URL:** `/candidates/me/invites/{invite_id}`
- **Authentication:** `Bearer <access_token>`
- **Authorization:** Candidate (Owner)
- **Request:**
  ```json
  { "response": "ACCEPTED" }
  ```
- **Response (`200 OK`):**
  ```json
  {
    "invite_id": "019171b3-4fa8-711b-8e7c-a49688bc8a9a",
    "status": "ACCEPTED",
    "updated_at": "2026-08-21T12:20:00Z"
  }
  ```
- **Note:** Accepting an invite does **not** automatically reveal identity — that remains a separate, explicit action (`10.1`). In the product flow, accepting is the natural moment for the frontend to prompt: *"Reveal your identity to Acme Ltd?"*
- **Error Responses:**
  - `400 Bad Request` — `response` missing or not one of `ACCEPTED`/`DECLINED`
  - `401 Unauthorized` — Token missing or expired
  - `403 Forbidden` — Caller is not the invite's recipient
  - `404 Not Found` — Invite not found
  - `409 Conflict` — Invite already responded to

---

## 3. Employer Domain (`/employers`)

**Employer profile status enum:** `{ PENDING, VERIFIED, UNVERIFIED }` — `verification_proof_url` and formal verification are reserved for a future roadmap step; for the MVP the field is stored but not enforced by any authorization check.

### 3.1 Get My Employer Profile

- **Method:** `GET`
- **URL:** `/employers/me`
- **Authentication:** `Bearer <access_token>`
- **Authorization:** Employer (Owner)
- **Request:** None
- **Response (`200 OK`):**
  ```json
  {
    "employer_id": "019171b3-4fa2-711b-8e7c-a49688bc8a94",
    "email": "recruiter@acme.com",
    "company_name": "Acme Ltd",
    "contact_first_name": "Almaz",
    "contact_last_name": "Kebede",
    "phone": "+251911556677",
    "location": "Addis Ababa, Ethiopia",
    "industry": "Software / Technology",
    "company_size": "11-50",
    "website_url": "https://acme.com",
    "verification_status": "PENDING",
    "created_at": "2026-08-21T10:00:00Z",
    "updated_at": "2026-08-21T10:00:00Z"
  }
  ```
- **Error Responses:**
  - `401 Unauthorized` — Token missing or expired
  - `404 Not Found` — Employer profile not found

---

### 3.2 Update My Employer Profile

- **Method:** `PATCH`
- **URL:** `/employers/me`
- **Authentication:** `Bearer <access_token>`
- **Authorization:** Employer (Owner)
- **Request:**
  ```json
  {
    "location": "Nairobi, Kenya",
    "industry": "Fintech",
    "company_size": "51-200",
    "website_url": "https://acme.com"
  }
  ```
- **Response (`200 OK`):** Full updated employer profile object (same shape as `3.1`)
- **Error Responses:**
  - `400 Bad Request` — Invalid field values
  - `401 Unauthorized` — Missing or expired token

---

## 4. Challenge Domain (`/challenges`)

**Challenge category enum:** `{ SOFTWARE_ENGINEERING }` — additional categories (`DATA_ANALYSIS`, `UI_UX`, etc.) are reserved for future expansion and are not active in the MVP.

### 4.1 List Challenges

- **Method:** `GET`
- **URL:** `/challenges`
- **Authentication:** None
- **Authorization:** Public
- **Query Parameters:**
  - `category` (string, optional, default: `SOFTWARE_ENGINEERING`) — filters the list to one category. Omit it and you get the default; pass `?category=SOFTWARE_ENGINEERING` explicitly for the same result. This is just how filters are passed on a `GET` request, since `GET` has no request body.
- **Response (`200 OK`):**
  ```json
  {
    "data": [
      {
        "id": "019171b3-4fa4-711b-8e7c-a49688bc8a96",
        "title": "Fix a broken task management API",
        "category": "SOFTWARE_ENGINEERING",
        "difficulty": "JUNIOR"
      }
    ]
  }
  ```
- **Error Responses:**
  - `400 Bad Request` — Invalid `category` value

---

### 4.2 Get Challenge

- **Method:** `GET`
- **URL:** `/challenges/{challenge_id}`
- **Authentication:** None
- **Authorization:** Public
- **Request:** Path parameter `challenge_id` (UUIDv7)
- **Response (`200 OK`):**
  ```json
  {
    "id": "019171b3-4fa4-711b-8e7c-a49688bc8a96",
    "title": "Fix a broken task management API",
    "category": "SOFTWARE_ENGINEERING",
    "difficulty": "JUNIOR",
    "description": "The repository contains a task management API with an intentional bug in ID lookup handling...",
    "repo_url": "https://github.com/kazi-challenges/task-api-001",
    "rubric": [
      { "competency": "CORRECT_IMPLEMENTATION", "weight": 0.30 },
      { "competency": "TESTS_EDGE_CASES", "weight": 0.20 },
      { "competency": "DEBUGGING", "weight": 0.20 },
      { "competency": "CODE_QUALITY", "weight": 0.15 },
      { "competency": "EXPLANATION", "weight": 0.15 }
    ],
    "followup_question": "Your implementation currently handles valid user IDs. What would you change if the API receives an ID that doesn't exist?"
  }
  ```
- **Error Responses:**
  - `404 Not Found` — Challenge does not exist

---

## 5. Submission Domain (`/submissions`)

**Submission status enum:** `{ QUEUED, RUNNING, COMPLETED, FAILED }`

### 5.1 Create Submission

- **Method:** `POST`
- **URL:** `/submissions`
- **Authentication:** `Bearer <access_token>`
- **Authorization:** Candidate
- **Request:**
  ```json
  {
    "challenge_id": "019171b3-4fa4-711b-8e7c-a49688bc8a96",
    "code": "string",
    "explanation": "string"
  }
  ```
- **Response (`202 Accepted`):**
  ```json
  {
    "submission_id": "019171b3-4fa5-711b-8e7c-a49688bc8a97",
    "status": "QUEUED"
  }
  ```
- **Note:** Code is dispatched to an isolated sandbox runner and executed asynchronously — it never runs on the main application server. Frontend should poll `5.2` for status.
- **Error Responses:**
  - `400 Bad Request` — Missing `challenge_id`, `code`, or `explanation`
  - `401 Unauthorized` — Token missing or expired
  - `404 Not Found` — Challenge does not exist
  - `409 Conflict` — Candidate already has an active submission for this challenge

---

### 5.2 Get Submission

- **Method:** `GET`
- **URL:** `/submissions/{submission_id}`
- **Authentication:** `Bearer <access_token>`
- **Authorization:** Submission Owner (Candidate)
- **Request:** Path parameter `submission_id` (UUIDv7)
- **Response (`200 OK`):**
  ```json
  {
    "submission_id": "019171b3-4fa5-711b-8e7c-a49688bc8a97",
    "challenge_id": "019171b3-4fa4-711b-8e7c-a49688bc8a96",
    "status": "COMPLETED",
    "test_results": {
      "passed": 8,
      "failed": 2,
      "total": 10,
      "execution_time": 1.82
    },
    "followup_status": "PENDING",
    "created_at": "2026-08-21T11:45:00Z"
  }
  ```
- **Error Responses:**
  - `401 Unauthorized` — Token missing or expired
  - `403 Forbidden` — Caller is not the submission owner
  - `404 Not Found` — Submission not found

---

### 5.3 Answer Follow-Up

- **Method:** `POST`
- **URL:** `/submissions/{submission_id}/followup`
- **Authentication:** `Bearer <access_token>`
- **Authorization:** Submission Owner (Candidate)
- **Request:**
  ```json
  { "answer": "I'd add a check that returns a 404 with a clear error message before attempting to look up the task, instead of letting it fail deeper in the service layer." }
  ```
- **Response (`200 OK`):**
  ```json
  {
    "submission_id": "019171b3-4fa5-711b-8e7c-a49688bc8a97",
    "followup_status": "COMPLETED",
    "updated_at": "2026-08-21T11:50:00Z"
  }
  ```
- **Error Responses:**
  - `400 Bad Request` — Missing `answer`
  - `401 Unauthorized` — Token missing or expired
  - `403 Forbidden` — Caller is not the submission owner
  - `404 Not Found` — Submission not found
  - `409 Conflict` — Submission is not yet `COMPLETED`, or follow-up already answered

---

## 6. Assessment Domain (`/assessments`)

**Integrity status enum:** `{ HIGH, REVIEW_RECOMMENDED }`
**Evidence status enum:** `{ STRONG_EVIDENCE, DEMONSTRATED, DEVELOPING, NEEDS_IMPROVEMENT }`

### 6.1 Get Assessment

- **Method:** `GET`
- **URL:** `/assessments/{assessment_id}`
- **Authentication:** `Bearer <access_token>`
- **Authorization:** Related Candidate (Owner)
- **Request:** Path parameter `assessment_id` (UUIDv7)
- **Response (`200 OK`):**
  ```json
  {
    "assessment_id": "019171b3-4fa6-711b-8e7c-a49688bc8a98",
    "submission_id": "019171b3-4fa5-711b-8e7c-a49688bc8a97",
    "score": 82,
    "integrity_status": "HIGH",
    "feedback": "Correctly identified and fixed the core bug; explanation was consistent with the implementation.",
    "evidence": [
      {
        "competency": "API Development",
        "status": "STRONG_EVIDENCE",
        "explanation": "Correctly implemented response handling for the failing endpoint."
      },
      {
        "competency": "Testing",
        "status": "DEVELOPING",
        "explanation": "2 of 10 tests failed, primarily around edge cases."
      }
    ],
    "created_at": "2026-08-21T11:55:00Z"
  }
  ```
- **Note:** `score` is computed against the challenge's fixed rubric (`4.2`). The AI service contributes qualitative signal only — it never sets `score` or `integrity_status` directly.
- **Error Responses:**
  - `401 Unauthorized` — Token missing or expired
  - `403 Forbidden` — Caller is not the related candidate
  - `404 Not Found` — Assessment not found (may still be processing — check submission status via `5.2` first)

---

## 7. Skills Passport Domain (`/passport`)

### 7.1 Get My Skills Passport

- **Method:** `GET`
- **URL:** `/passport/me`
- **Authentication:** `Bearer <access_token>`
- **Authorization:** Candidate (Owner)
- **Request:** None
- **Response (`200 OK`):**
  ```json
  {
    "candidate_id": "019171b3-4fa0-711b-8e7c-a49688bc8a92",
    "anonymous_id": "KZ-1042",
    "category": "SOFTWARE_ENGINEERING",
    "overall_score": 82,
    "skills": [
      { "name": "JavaScript", "status": "STRONG_EVIDENCE" },
      { "name": "API Development", "status": "STRONG_EVIDENCE" },
      { "name": "Debugging", "status": "DEMONSTRATED" },
      { "name": "Testing", "status": "DEVELOPING" }
    ],
    "evidence_summary": {
      "challenges_completed": 1,
      "tests_passed": "8/10",
      "explanation_submitted": true,
      "followup_completed": true
    }
  }
  ```
- **Error Responses:**
  - `401 Unauthorized` — Token missing or expired
  - `404 Not Found` — No completed assessment yet for this candidate

---

## 8. Opportunity Domain (`/opportunities`)

**Note:** Per the PRD, KAZI is explicitly *not* building a full employer job-posting platform for the MVP ("No employer CRM", "No full job marketplace"). Opportunities are staged/seeded records for the demo, not created by employers through the API.

### 8.1 List Opportunities

- **Method:** `GET`
- **URL:** `/opportunities`
- **Authentication:** None
- **Authorization:** Public
- **Response (`200 OK`):**
  ```json
  {
    "data": [
      {
        "id": "019171b3-4fa7-711b-8e7c-a49688bc8a99",
        "title": "Junior Software Engineer",
        "company": "Acme Ltd",
        "requirements": ["JavaScript", "REST APIs", "Debugging", "Testing"]
      }
    ]
  }
  ```
- **Error Responses:**
  - None (public, static-ish list in MVP)

---

### 8.2 Get Match for Candidate

- **Method:** `GET`
- **URL:** `/opportunities/{opportunity_id}/match`
- **Authentication:** `Bearer <access_token>`
- **Authorization:** Candidate (self-match only in MVP)
- **Request:** Path parameter `opportunity_id` (UUIDv7)
- **Response (`200 OK`):**
  ```json
  {
    "opportunity_id": "019171b3-4fa7-711b-8e7c-a49688bc8a99",
    "candidate_id": "019171b3-4fa0-711b-8e7c-a49688bc8a92",
    "match_score": 82,
    "matched_requirements": ["JavaScript", "REST APIs", "Debugging"],
    "gaps": ["Testing"]
  }
  ```
- **Note:** Match score is `(matched requirement count / total requirement count) * 100`, weighted against evidence status — no ML model required for MVP.
- **Error Responses:**
  - `401 Unauthorized` — Token missing or expired
  - `404 Not Found` — Opportunity not found, or candidate has no passport yet

---

## 9. Candidate Discovery Domain (`/employer`)

**This domain is the anonymity boundary — and requires a signed-in, authenticated `EMPLOYER`.** There is no public or unauthenticated way to browse candidates; every endpoint below requires `Bearer <access_token>` from a registered, logged-in employer. Every response here must return only anonymized candidate data unless `identity_revealed = true` for that specific candidate/opportunity pair. This is enforced at the query layer, not filtered client-side.

### 9.1 Get Candidate Evidence (Anonymous View)

- **Method:** `GET`
- **URL:** `/employer/candidates/{candidate_id}`
- **Authentication:** `Bearer <access_token>`
- **Authorization:** `EMPLOYER`
- **Request:** Path parameter `candidate_id` (UUIDv7)
- **Response (`200 OK`, not revealed):**
  ```json
  {
    "anonymous_id": "KZ-1042",
    "category": "SOFTWARE_ENGINEERING",
    "overall_score": 82,
    "skills": [
      { "name": "JavaScript", "status": "STRONG_EVIDENCE" }
    ],
    "evidence_summary": {
      "challenge_completed": true,
      "tests_passed": "8/10",
      "followup_completed": true
    },
    "identity_revealed": false
  }
  ```
- **Response (`200 OK`, revealed):** Same shape, plus:
  ```json
  {
    "name": "Aisha Mohammed",
    "country": "Nigeria",
    "identity_revealed": true
  }
  ```
- **Note:** `identity_revealed` here is scoped to *this employer's* relationship with the candidate (via the opportunity they invited on) — never a global, publicly visible flag. A different employer querying the same `candidate_id` for a different opportunity the candidate hasn't revealed to would still get `identity_revealed: false`.
- **Error Responses:**
  - `401 Unauthorized` — Token missing or expired
  - `403 Forbidden` — Caller is not `EMPLOYER`
  - `404 Not Found` — Candidate not found or has no completed passport

---

### 9.2 Invite Candidate

- **Method:** `POST`
- **URL:** `/employer/candidates/{candidate_id}/invite`
- **Authentication:** `Bearer <access_token>`
- **Authorization:** `EMPLOYER`
- **Request:**
  ```json
  { "opportunity_id": "019171b3-4fa7-711b-8e7c-a49688bc8a99" }
  ```
- **Response (`201 Created`):**
  ```json
  {
    "invite_id": "019171b3-4fa8-711b-8e7c-a49688bc8a9a",
    "candidate_id": "019171b3-4fa0-711b-8e7c-a49688bc8a92",
    "opportunity_id": "019171b3-4fa7-711b-8e7c-a49688bc8a99",
    "status": "PENDING",
    "created_at": "2026-08-21T12:10:00Z"
  }
  ```
- **Note:** This only creates the invite record — it does **not** reveal identity or notify the candidate of who's asking beyond the opportunity context. The candidate sees and responds to it via `2.4`/`2.5`.
- **Error Responses:**
  - `401 Unauthorized` — Token missing or expired
  - `403 Forbidden` — Caller is not `EMPLOYER`
  - `404 Not Found` — Candidate or opportunity not found
  - `409 Conflict` — Invite already sent for this candidate/opportunity pair

---

## 10. Identity Domain (`/identity`)

### 10.1 Reveal Identity

- **Method:** `POST`
- **URL:** `/identity/reveal`
- **Authentication:** `Bearer <access_token>`
- **Authorization:** Candidate (Owner) — **this action can only be initiated by the candidate, never the employer**
- **Request:**
  ```json
  { "opportunity_id": "019171b3-4fa7-711b-8e7c-a49688bc8a99" }
  ```
- **Response (`200 OK`):**
  ```json
  {
    "candidate_id": "019171b3-4fa0-711b-8e7c-a49688bc8a92",
    "opportunity_id": "019171b3-4fa7-711b-8e7c-a49688bc8a99",
    "identity_revealed": true,
    "revealed_at": "2026-08-21T12:15:00Z"
  }
  ```
- **Note:** Reveal is scoped per opportunity, not global — a candidate can be anonymous to one employer and revealed to another for a different opportunity.
- **Error Responses:**
  - `400 Bad Request` — Candidate has no `name`/`country` set yet (see `2.3`)
  - `401 Unauthorized` — Token missing or expired
  - `403 Forbidden` — Caller is not the candidate
  - `404 Not Found` — Opportunity not found
  - `409 Conflict` — Identity already revealed for this opportunity

---

## Status Code Conventions

| Code | Meaning |
|---|---|
| `200` | Success (sync) |
| `201` | Resource created |
| `202` | Accepted, processing async (submissions, assessment) |
| `204` | Success, no content |
| `400` | Validation error |
| `401` | Missing/invalid/expired auth |
| `403` | Authenticated but not permitted for this resource |
| `404` | Not found |
| `409` | Conflict (duplicate action, invalid state transition) |
| `500` | Server error |

## Full Invite-to-Reveal Flow (for reference)

```
Employer views anonymous candidate evidence  (9.1)
        ↓
Employer sends invite tied to an opportunity  (9.2)
        ↓
Candidate sees job title + company via their invite list  (2.4)
        ↓
Candidate accepts or declines  (2.5)
        ↓
If accepted, candidate is prompted to reveal identity  (10.1)
        ↓
Employer now sees name + country for that candidate/opportunity only  (9.1, revealed=true)
```

## Cross-Cutting Notes for Frontend

- **Two separate registration flows:** candidates and employers hit different endpoints (`1.1` vs `1.2`) with different required fields — don't build one shared signup form.
- **No anonymous employer access:** every `/employer/*` and `/employers/*` call requires a logged-in `EMPLOYER` token. There is no public candidate browsing.
- **Async flows:** `5.1` (Create Submission) and the assessment step behind it are asynchronous. Poll `5.2 GET /submissions/{id}` until `status: COMPLETED`, then `7.1` for the finished passport. Consider a light polling interval (e.g. every 2s) rather than long-polling for MVP.
- **Anonymity boundary is structural:** the `/employer/*` domain is deliberately separate from `/candidates/*` and `/passport/*` — never call candidate-facing endpoints from employer-facing screens, even if the data looks similar.
- **Reveal is per-opportunity, not global:** don't cache `identity_revealed` as a single boolean on the candidate — it varies by `opportunity_id`.
- **Enums are fixed sets** (`status`, `category`, `integrity_status`, evidence `status`, invite `status`) — confirm with backend before relying on any value not listed above.
