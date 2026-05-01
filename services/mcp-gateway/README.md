# MCP Gateway

The in-house ClawFirm MCP gateway. Every tool call from every engine (Dify, n8n, LangGraph, OpenHands, the OpenClaw shell) flows through this binary.

## What it enforces

- **mTLS** on every connection (gateway ↔ engine, gateway ↔ MCP server).
- **OIDC identity** on every request, propagated as a verified claim into the upstream MCP server.
- **Per-tool allowlists** via `config/allowlist.yaml` (version-controlled).
- **Per-agent egress policy** — what hosts/ports a tool is permitted to reach.
- **MCP 2025-03-26 token-passthrough prohibition** — the gateway never forwards a user's bearer token to an MCP server. It mints its own short-lived token bound to the request.
- **Hash-chained audit log**, posted to ClawSecure's `/api/events` endpoint.
- **Confused-deputy protection** — explicit consent required when a tool calls another tool on behalf of an end user.

## Why in-house

ADR [0004](../../docs/adr/0004-in-house-mcp-gateway.md) explains the rationale. tl;dr: critical infrastructure, has to work air-gapped, has to move at the speed of the spec rather than at the speed of a SaaS vendor.

## Layout

```
mcp-gateway/
├── cmd/mcp-gateway/main.go
├── internal/
│   ├── allowlist/   - YAML allowlist parser, per-tool / per-agent rules
│   ├── audit/       - hash-chained audit log, posts to ClawSecure
│   ├── identity/    - OIDC verification + token minting
│   └── policy/      - request decision evaluation
├── config/
│   └── allowlist.yaml.example
├── go.mod
└── Dockerfile
```

## Configuration

`MCP_GW_LISTEN`, `MCP_GW_CONFIG`, `MCP_GW_AUDIT_SINK`, `MCP_GW_OIDC_ISSUER` — see `cmd/mcp-gateway/main.go`.
