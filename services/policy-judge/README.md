# Policy Judge

The local sensitivity / complexity / routing classifier. A small (4B-param) Qwen3.5 model wrapped in a thin FastAPI server that accepts `{tool, args, identity, context}` and returns `{class, confidence, reason, suggested_route}`.

Used by:

- **ClawRails** — to decide local-vs-frontier routing per request.
- **MCP gateway** — to short-circuit Tier 1 smart-approval decisions when confidence is high.
- **EdgeClaw shell** — to surface the classifier's reasoning to the user when prompting for tier promotion.

## Why

A 4B-parameter model running on Ollama is fast (~50 ms on CPU, sub-10 ms on GPU), private (never leaves the box), and good enough at structured classification when you give it a clear rubric and few-shot examples. EdgeClaw's empirical data suggests this combination lets 60-80% of agent traffic stay on local models without quality degradation.

## Layout

```
policy-judge/
├── policy_judge/
│   ├── server.py     - FastAPI server
│   ├── cache.py      - 5-minute SHA-256 cache to avoid rerunning identical classifications
│   └── prompts/
│       ├── sensitivity.txt
│       ├── complexity.txt
│       └── routing.txt
└── Dockerfile
```

## Configuration

| Env var | Default |
|---|---|
| `JUDGE_MODEL` | `qwen3.5-4b` |
| `JUDGE_BACKEND` | `ollama` |
| `OLLAMA_HOST` | `http://ollama:11434` |
| `JUDGE_CACHE_TTL` | `300` (seconds) |

## Quality control

- Returns `{confidence}` per call. Callers reject `confidence < 0.6` and escalate to manual approval.
- Re-benchmarked every 6 months against a held-out test set — the rubrics in `prompts/` are versioned with the test set.
