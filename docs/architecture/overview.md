# Architecture Overview

ClawFirm is a layered system designed for three priorities, in order:

1. **Lowest cost to run** — local models first; frontier API only when a local judge says it's worth it.
2. **Secure by default** — sandbox-first, default-deny egress, signed skills, mandatory audit log.
3. **Easiest install / UX** — one-line installer; opinionated defaults that flip OpenClaw's permissive ones.
4. **Maximum flexibility** — BYO models, BYO tools via MCP — once the safe defaults are in place.

## Read in this order

1. [`reference-architecture.md`](./reference-architecture.md) — the diagram + the layer-by-layer explanation.
2. [`routing-logic.md`](./routing-logic.md) — what happens when a user message arrives.
3. [`governance-ladder.md`](./governance-ladder.md) — the four-tier approval model.
4. [`memory-service.md`](./memory-service.md) — the event-sourced memory design.
5. [`../adr/`](../adr/) — the architecture decision records (why each major choice was made).

## Component map (one-screen reference)

| Layer | Primary | Where it lives |
|---|---|---|
| LLM serving (single-user) | Ollama | `deploy/{solo,smb}/` |
| LLM serving (multi-user GPU) | vLLM | `deploy/{smb,enterprise}/` |
| Inference router + cost layer | **ClawRails** (sister repo) | `integrations/clawrails/` |
| Cost-routing judge | Qwen3.5-4B local | `services/policy-judge/` |
| Personal-assistant shell | EdgeClaw (with ClawFirm hardening) | container in `deploy/` |
| Visual low-code builder | Dify | container in `deploy/` |
| Workflow engine | n8n | container in `deploy/` |
| Programmable orchestration | LangGraph 1.x | container in `deploy/` |
| Coding-agent runtime | OpenHands V1 SDK | container in `deploy/` |
| Memory / RAG store (default) | pgvector on shared Postgres | `services/memory-service/` |
| Code/tool sandbox | DifySandbox + gVisor + Firecracker (tiered) | image refs in `deploy/` |
| Secrets | OpenBao | container in `deploy/` |
| Identity / SSO | Authentik | container in `deploy/` |
| Zero-trust net | Headscale + Tailscale clients | container in `deploy/`, host install via ClawBoot |
| Observability | Langfuse | container in `deploy/` |
| Ingress | Caddy 2 | container in `deploy/` |
| Orchestration (SMB) | k3s | `deploy/smb/k3s/` |
| Orchestration (Enterprise) | k0s + Flux | `deploy/enterprise/` |
| One-command PaaS UX | Coolify | optional, `deploy/smb/` |
| MCP gateway | **In-house** | `services/mcp-gateway/` |
| Approvals + policy + audit | **ClawSecure** (sister repo) | `integrations/clawsecure/` |
| Bare-metal install / cloud-init | **ClawBoot** (sister repo) | `integrations/clawboot/` |

## Key open risks

These are tracked as live risks; ADRs may supersede the current choices as we learn.

1. **Upstream OpenClaw / EdgeClaw stability.** Treat the shell as swappable; alternatives (Goose, custom React) can replace it via the Channel API.
2. **n8n SUL ambiguity for managed-service operators.** See [`LICENSING.md`](../../LICENSING.md). The `permissive-only` build profile drops n8n in favor of Activepieces.
3. **Dify Open Source License + appearance-patent clause.** Reskinning the Dify UI may trigger restrictions; `permissive-only` profile drops Dify in favor of Flowise.
4. **MCP protocol churn.** 43% of public MCP servers had command-injection vulns in 2026 audits — we treat third-party MCP servers as untrusted and sandbox them.
5. **Apple Silicon sandboxing.** gVisor and Firecracker don't run natively on macOS; Solo tier on macOS uses Docker Desktop's VZ + DifySandbox.
6. **Frontier-model dependency drift.** Pin model digests and re-benchmark every 6 months.
7. **Tier 3 prompt-injection exfiltration.** Even sandboxed, ingested-document prompt injection can cause logical exfiltration. Tier 3 requires explicit human approval for outbound tool calls and a two-person rule at Enterprise.
8. **Memory consistency.** Event sourcing solves this; both projections rebuild deterministically from the log.
9. **License compatibility of the assembled image.** Per-image SBOM with SPDX tags; permissive-only profile available.
10. **Air-gapped model updates.** Detached Sigstore signatures; customer brings updates via approved cross-domain solution; ClawFirm verifies before applying.
