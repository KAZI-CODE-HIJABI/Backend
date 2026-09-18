# Assessment rules for the hackathon MVP

The final score comes from KAZI's fixed rubric, not from the AI provider:

| Component | Weight |
|---|---:|
| Correct implementation | 30% |
| Tests and edge cases | 20% |
| Debugging | 20% |
| Code quality | 15% |
| Explanation | 15% |

Each component is a value from 0 to 1. The backend computes and rounds:

```text
100 × (0.30×correct + 0.20×tests + 0.20×debugging + 0.15×quality + 0.15×explanation)
```

Automated test results supply objective evidence. AI output is restricted to
validated observations that cite candidate work or trusted test results. It cannot
return a final score, candidate ranking, hiring recommendation, cheating probability
or integrity status.

Evidence statuses use these thresholds: `STRONG_EVIDENCE` at 0.85 or above,
`DEMONSTRATED` at 0.60–0.849, `DEVELOPING` at 0.35–0.599 and
`NEEDS_IMPROVEMENT` below 0.35.

Integrity is `HIGH` only when required evidence is present and no material
contradiction is recorded. Otherwise it is `REVIEW_RECOMMENDED`. It is a prompt for
human review, not a claim that a candidate used AI or cheated.
