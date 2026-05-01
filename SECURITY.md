# Security

ClawFirm exists in part because the OpenClaw family ships insecure-by-default. This document is the threat model and the secure-default matrix.

---

## Threat model

ClawFirm is designed against the following threats, in roughly descending order of likelihood:

1. **Internet exposure of the agent shell.** Bitsight observed >30,000 OpenClaw instances exposed on the public internet in early 2026. ClawFirm binds to `127.0.0.1` by default and routes all remote access through Headscale + Tailscale.
2. **Prompt injection via ingested documents, web fetches, MCP server output, or chat-app channels** that causes the agent to invoke tools the user did not intend. Mitigation: every tool call goes through the MCP gateway, which applies allowlists, identity scoping, and (at Tier 2+) human approval via ClawSecure.
3. **Compromised or malicious skills / plugins / MCP servers.** McCarty's research identified 386 malicious skill packages in 2026. Mitigation: signed skill registry (Sigstore); blocking dangerous-code scan; default-deny egress from the sandbox; `--dangerously-force-unsafe-install` removed entirely on Enterprise tier.
4. **Credential exfiltration.** Plaintext credential storage was historically common in this family. Mitigation: all credentials → OpenBao with envelope encryption. n8n's `KeyManagerService` rotation enabled. First-launch wizard auto-generates a sealed root key.
5. **Lateral movement from a compromised tool/agent into the host or other agents.** Mitigation: tiered sandbox — DifySandbox (seccomp whitelist) for Python/Node; gVisor (user-space kernel) for arbitrary tools; Firecracker microVM for autonomous Tier 3.
6. **Memory tampering.** A malicious tool writing into the agent's long-term memory to influence future sessions. Mitigation: Memory Service writes are append-only with a hash-chained audit log; LLM-driven writes are themselves a tool call subject to approval.
7. **Token passthrough / confused-deputy attacks against MCP.** The MCP 2025-03-26 spec explicitly prohibits token passthrough; ClawFirm enforces this in the gateway.
8. **Supply-chain attacks on container images.** Mitigation: Sigstore signing on all ClawFirm-published images; Cosign verification in `clawfirm upgrade`; per-image SBOMs.
9. **Insider abuse / accidental misuse.** Mitigation: ClawSecure timeline + approval queue; per-action audit trail; role-based access via Authentik.

Out of scope for v1: hardware-level attacks (Spectre/Meltdown class), nation-state-grade supply-chain compromise of upstream components, and zero-days in the host kernel. Confidential-compute support (SEV-SNP / TDX) is on the Q4-26 roadmap to address the first two for Tier 3 workloads.

---

## Secure-default matrix

These are the defaults ClawFirm ships out of the box, side-by-side with what upstream components default to.

| Setting | Upstream default | **ClawFirm default** |
|---|---|---|
| Bind address | OpenClaw historically bound `0.0.0.0:18789` | `127.0.0.1` + Tailscale/Headscale overlay |
| Approval mode | OpenClaw: trust-leaning; Goose: autonomous; AnyClaw: `danger-full-access` | **Tier 0 (chat-only)** at first launch |
| Sandbox scope | `agent` or `shared` | `session` — one sandbox per conversation |
| Workspace access from sandbox | varies | `none` — sandbox cannot read host workspace by default |
| Skill install | optional dangerous-scan | **mandatory blocking** dangerous-scan + Sigstore signature required |
| Egress from tools | open | **default-deny**; per-agent allowlist; enforced at MCP gateway |
| Channel adapters (WhatsApp/Slack/Telegram/etc.) | enabled if creds present | disabled until wizard explicitly enables them |
| Headless mode sandbox | optional | **required** — must be DockerWorkspace or stronger |
| `--dangerously-*` flags | available | removed entirely on Enterprise tier; gated behind `--i-understand-the-risk` elsewhere |
| Vendor telemetry | varies | off by default; explicit opt-in |
| Credential storage | sometimes plaintext | impossible — OpenBao envelope encryption is mandatory |
| Tool-call audit log | optional | always-on, hash-chained, exportable |
| Frontier API egress | open if API key present | gated by ClawRails policy + ClawSecure rule + per-user policy |
| MCP token passthrough | varies | **prohibited** at the gateway (per MCP 2025-03-26 spec) |
| Default model for new sessions | varies | local small model first; frontier API requires explicit policy enable |

---

## Governance ladder

ClawFirm's four-tier approval ladder is enforced uniformly across every framework (Dify, n8n, LangGraph, OpenHands, Goose, OpenClaw shell) via the `approval-shim` gRPC service, which adapts each framework's native approval primitive to a single `ClawSecure` REST backend.

| Tier | Behavior | Recommended for |
|---|---|---|
| **0 — Chat-only** | No tools accessible. Pure local LLM Q&A. All channels disabled. | First-run default. Untrusted environments. |
| **1 — Smart approval** | Local policy LLM (Qwen3.5-4B) classifies each tool call as `safe` / `risky` / `dangerous`. `safe` auto-approved; others prompt. | Solo and SMB normal use. |
| **2 — Manual approval** | Every tool call prompts via web UI / Slack / Telegram. | High-stakes work; Enterprise default for write/egress tools. |
| **3 — Autonomous** | No prompts. **Only inside a fresh Firecracker microVM** with default-deny egress and a wall-clock budget. Explicit `--i-understand-the-risk` opt-in. | Background batch jobs in tightly scoped sandboxes. |

---

## Reporting a vulnerability

Email `security@clawfirm.io` with PGP key fingerprint published in `.well-known/security.txt`. Please do not file public issues for security reports. We aim to acknowledge within 48 hours and to ship a fix or mitigation within 14 days for high-severity issues.

---

## Verifying releases

Every ClawFirm release is signed with Sigstore. Verify before installing:

```bash
cosign verify-blob \
  --certificate clawfirm-${VERSION}.tar.gz.crt \
  --signature clawfirm-${VERSION}.tar.gz.sig \
  --certificate-identity-regexp 'https://github.com/jimmershere/clawfirm/.github/workflows/release\.yml@.*' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com \
  clawfirm-${VERSION}.tar.gz
```

`scripts/verify-signatures.sh` does this for you.
