# Memory Service

Event-sourced memory for ClawFirm. Exposes both EdgeClaw-style hierarchical memory (project/timeline/profile -> fragment -> conversation) and Dify-style flat RAG behind a single gRPC API. Both representations are projections of an append-only event log.

See [`docs/architecture/memory-service.md`](../../docs/architecture/memory-service.md) for the full design.

## Layout

```
memory-service/
├── memory_service/
│   ├── api.py        - FastAPI + gRPC server entrypoint
│   ├── events.py     - append-only event log + hash chain
│   ├── tree.py       - hierarchical projection (EdgeClaw-style)
│   ├── rag.py        - flat retrieval projection (Dify-compatible)
│   └── stores/
│       ├── pgvector.py  - default
│       ├── qdrant.py    - SMB+ scale
│       └── lancedb.py   - embedded for solo
├── proto/
│   └── memory.proto
├── pyproject.toml
└── Dockerfile
```

## Configuration

| Env var | Default | Notes |
|---|---|---|
| `MEMORY_STORE` | `pgvector` | `pgvector` / `qdrant` / `lancedb` |
| `MEMORY_PG_DSN` | (required for pgvector) | Postgres DSN |
| `MEMORY_EVENT_LOG` | `/var/lib/memory/events` | Path to event log directory |
| `MEMORY_HASH_CHAIN` | `true` | Enable tamper-evident hash chain |
