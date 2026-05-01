# Approval Shim

A thin gRPC service that adapts each agent framework's native approval primitive to a single backend (ClawSecure REST). One uniform interface — `ApprovalService.RequestApproval(...)` — regardless of whether the caller is Dify, n8n, LangGraph, OpenHands, Goose, or the OpenClaw shell.

## Why

Every framework expresses approvals differently:

- **OpenHands** has a `Conversation` confirmation policy.
- **Goose** has approval modes.
- **LangGraph** uses `interrupt_before` and human-in-the-loop checkpointing.
- **n8n** has dedicated Wait / Approve nodes.
- **Dify** has human-input nodes inside workflows.
- **OpenClaw** uses a per-skill `autoAllowSkills` model with a stdin emitter.

Without a shim, every framework would have its own approval queue, its own audit log, and its own policy engine — a recipe for inconsistency and gaps. The shim collapses all of these into one `ApprovalService` interface that posts to ClawSecure.

## Layout

```
approval-shim/
├── cmd/approval-shim/main.go
├── internal/
│   └── adapters/
│       ├── clawsecure.go    - the backend - posts to /api/approvals
│       ├── langgraph.go     - HTTP shim that LangGraph's interrupt_before hits
│       ├── n8n.go           - webhook adapter for n8n approval nodes
│       ├── dify.go          - webhook adapter for Dify human-input nodes
│       └── openhands.go     - gRPC adapter for OpenHands Conversation policy
├── go.mod
└── Dockerfile
```

## Configuration

| Env var | Default |
|---|---|
| `SHIM_LISTEN` | `:50051` |
| `CLAWSECURE_URL` | `http://clawsecure:3188` |
