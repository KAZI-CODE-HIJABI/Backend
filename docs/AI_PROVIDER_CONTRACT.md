# AI provider contract

`internal/ai.Provider` is the only boundary where an LLM provider is introduced.
Its input contains a challenge description, rubric version, candidate code and
explanation, a fixed follow-up question/answer, and trusted test results. It must
not receive candidate identity, employer identity, passwords, tokens or database URLs.

The provider returns only supported qualitative observations and feedback. Each
observation must identify an allowed rubric competency and cite submitted or trusted
evidence. The AI service rejects unsupported competencies, duplicates and unsupported
observations before they reach assessment persistence.

Do not ask the model for a score, ranking, hire/reject decision, cheating probability
or integrity status. The `assessment` module uses KAZI's versioned fixed rubric and
integrity rules to make those decisions.

Before adding a real provider, supply the provider name, model, API key through the
environment or secret manager, timeout, retry limits and a prompt version. Do not put
the key in Git or send it from the frontend.
