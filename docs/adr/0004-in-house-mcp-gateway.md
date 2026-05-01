# ADR-0004: In-house MCP gateway, not a third-party SaaS

- **Status:** Accepted
- **Date:** 2026-04-30
- **Deciders:** @jimmershere

## Context

ClawFirm needs an MCP gateway: a single point that mediates every tool call, enforcing identity, allowlists, audit, and the MCP 2025-03-26 spec rules (notably the prohibition on token passthrough).

Several commercial MCP gateways exist (TrueFoundry, MintMCP, Cloudflare's reference architecture). Several open implementations exist as libraries.

Industry context: 43% of public MCP servers were reported vulnerable to command injection in 2026 audits. The MCP spec has changed materially between 2024-11 and 2025-03-26, with OAuth 2.1 / Dynamic Client Registration adoption still uneven.

## Decision

Ship our own MCP gateway as a small Go binary in `services/mcp-gateway/`. Treat it as critical infrastructure that ClawFirm owns end-to-end.

The gateway enforces:

- **mTLS** on every connection (gateway ↔ engine, gateway ↔ MCP server).
- **OIDC identity** on every request, propagated as a verified claim.
- **Per-tool allowlists** in YAML, version-controlled.
- **Per-agent egress policy** — what hosts/ports a tool may reach.
- **Token passthrough prohibition** — the gateway never forwards a user's bearer token to an MCP server. It mints its own short-lived token bound to the request.
- **Hash-chained audit log**, written to ClawSecure.
- **Confused-deputy protection** — explicit consent required when a tool calls another tool on behalf of an end user (per MCP spec).

All bundled MCP servers run inside the sandbox layer (ADR-0003), so even a compromised MCP server cannot directly escape.

## Consequences

**Easier:**
- No external SaaS dependency.
- Air-gapped installs work without internet.
- We can move at the speed of the spec, not at the speed of a commercial vendor's roadmap.
- One place to fix any new MCP-related CVE for all customers.

**Harder:**
- We own the maintenance burden, including tracking spec churn.
- We have to write our own observability and rate-limiting (mitigated by emitting OTel and reusing existing rate-limit libraries).

## Alternatives considered

- **TrueFoundry MCP Gateway / MintMCP.** Rejected as default — commercial SaaS / hosted dependency conflicts with the air-gap and cost-first goals. Available as opt-in.
- **Build on the official MCP SDK only, no gateway.** Rejected — leaves identity, audit, and policy enforcement to each engine, which is exactly the OpenClaw failure mode we're trying to fix.
- **Cloudflare's MCP "server portal" pattern.** Adopted in part — we use the "Code Mode" token-compression idea but run it ourselves rather than depending on Cloudflare Workers.
