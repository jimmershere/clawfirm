# Memory Service

The Memory Service unifies two memory representations behind one API:

1. **EdgeClaw-style hierarchical memory** — project / timeline / profile → fragment → conversation. Useful for "remember this about me," "what did we decide last sprint," and personalization across sessions.
2. **Dify-style flat RAG** — chunked-and-embedded retrieval over user-supplied knowledge bases. Useful for grounded Q&A over docs.

Both are projections of an **append-only event log**. We never have to reconcile two writable stores.

## Why event sourcing

The naïve approach is to write to both stores from every component. That creates consistency hell — Dify writes a chunk, the OpenClaw shell writes a memory fragment, the LangGraph agent writes a conversation summary, and now they all disagree about the canonical state. Event sourcing solves this: every change is one append to an event log, and the two stores are projections that get rebuilt deterministically.

This mirrors the OpenHands V1 SDK design (`Conversation` + `EventLog`), so the pattern is well-trodden in the agent space.

## API

The service exposes a gRPC API (with a thin REST gateway). The four primary methods:

```proto
service MemoryService {
  rpc Append(AppendRequest) returns (AppendReply);
  rpc Recall(RecallRequest) returns (RecallReply);
  rpc Forget(ForgetRequest) returns (ForgetReply);
  rpc Snapshot(SnapshotRequest) returns (SnapshotReply);
}
```

- `Append` — write a new event. Returns the event ID and the resulting hash-chain head.
- `Recall` — query. Specify hierarchical scope (project / timeline / profile) OR flat RAG retrieval (top-k by embedding similarity). Both can be combined.
- `Forget` — soft-delete via tombstone event. Hard-delete is a separate admin operation that compacts the event log and rebuilds projections (used for GDPR right-to-erasure).
- `Snapshot` — produce a point-in-time export.

## Stores

Three pluggable storage backends, selected by tier:

| Store | When |
|---|---|
| **pgvector** | Default. Uses the same Postgres instance Dify and n8n already need — zero new ops. |
| **Qdrant** | When workload >5–10M vectors or hybrid search / late-interaction is required. |
| **LanceDB** | Embedded; solo tier on a Raspberry Pi or laptop with no Postgres. |

The store choice is a runtime config; the API surface is identical.

## Hash chain

Every event includes the SHA-256 of the previous event ID + canonical JSON of the payload. The chain head is published periodically to a tamper-evident location (Sigstore Rekor at Enterprise tier; a local file at Solo). This makes silent post-hoc memory editing detectable.

## Memory writes are themselves a tool call

A common attack on agentic systems is for a malicious tool to write into the agent's memory to influence future sessions. ClawFirm treats `MemoryService.Append` as **just another tool call**, so it goes through the same MCP gateway, the same approval ladder, and the same audit log. At Tier 1 the policy-judge sees the write, classifies it (e.g., "agent wants to record a user preference") and auto-approves the safe ones. At Tier 2+ the operator sees them.

## Privacy scopes

Every event is tagged with a privacy scope:

- `public` — visible to all users on this ClawFirm instance.
- `tenant` — visible to all users within a tenant (Enterprise tier).
- `user` — visible only to the user who owns the event.
- `session` — discarded at session end.
- `secret` — encrypted with a per-user key from OpenBao; never sent to a frontier API.

The cost router checks the scope before sending content to a frontier API; `secret`-scoped content is silently filtered out (or the call is downgraded to a local model).
