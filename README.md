# ClawFirm

[![CI](https://img.shields.io/github/actions/workflow/status/jimmershere/clawfirm/ci.yml?branch=main&label=build)](https://github.com/jimmershere/clawfirm/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/license-MIT-green.svg)](./LICENSE)
[![Status: Alpha](https://img.shields.io/badge/status-alpha-orange.svg)](./ROADMAP.md)

**The AI platform that doesn't get you hacked. One install. Three tiers. Zero trust from day one.**

ClawFirm is the secure-by-default AI factory in a box for builders who want self-hosted agents, local inference, and governance that starts locked down instead of wide open.

ClawFirm is a meta-distribution that composes:

- [`clawrails`](https://github.com/jimmershere/clawrails) — hybrid local/cloud routing, cost attribution, kill switches.
- [`clawboot`](https://github.com/jimmershere/clawboot) — bootable USB / cloud-init installer that provisions a hardened host.
- [`clawsecure`](https://github.com/jimmershere/clawsecure) — local-first dashboard for events, approvals, policy rules, and artifact controls.

…with a curated, opinionated upstream stack: **Ollama** + **vLLM** for inference, **Dify** for low-code agents and RAG, **n8n** for workflows, **LangGraph 1.x** for durable orchestration, **OpenHands V1** for coding agents, **OpenBao** for secrets, **Authentik** for identity, **Headscale** + Tailscale clients for zero-trust networking, **Langfuse** for observability, and a thin in-house **MCP gateway** to mediate every tool call.

One bundle. Three tiers. Same stack.

```
+---------------------------------------------------+
|  Solo  (laptop / single VPS / Pi)                 |  Docker Compose
+---------------------------------------------------+
|  SMB   (1-3 node mini-cluster)                    |  k3s + Coolify
+---------------------------------------------------+
|  Enterprise  (HA, RBAC, air-gapped capable)       |  k0s + Flux GitOps
+---------------------------------------------------+
```

---

## Why ClawFirm exists

> [!IMPORTANT]
> Bitsight observed **30,000+ internet-exposed OpenClaw instances** in early 2026. Microsoft Defender, McAfee, and other researchers all warned that the default posture was too permissive for real-world deployment.
>
> ClawFirm flips that model: **chat-only at first run**, sandbox-first execution, default-deny egress, signed skills, and a mandatory audit trail.

Upstream OpenClaw, AnyClaw, and Goose all make it too easy to hand your host over to an agent. ClawFirm exists because we were tired of security being left as homework.

See [`SECURITY.md`](./SECURITY.md) for the full threat model and the secure-default matrix.

---

## Quick start (Solo tier)

```bash
curl -fsSL https://get.clawfirm.io | sh
clawfirm init --tier solo
clawfirm up
```

Open <http://localhost:7878> for the assistant shell, <http://localhost:3188> for the ClawSecure dashboard, <http://localhost:3000> for Langfuse.

That's it. Local Qwen3.5-4B model is pulled and ready. No frontier API key required. No tool can run yet — you're in **Tier 0 (chat-only)** by design. Promote yourself up the [governance ladder](./docs/architecture/governance-ladder.md) when you're ready.

---

## Repo layout

```
clawfirm/
├── cmd/clawfirm/              # Go CLI entry point
├── internal/                  # CLI internals (tier detect, installer, seed)
├── deploy/
│   ├── solo/                  # Docker Compose for single-node
│   ├── smb/k3s/               # Kustomize manifests for k3s
│   ├── enterprise/            # k0sctl + Flux GitOps for HA clusters
│   └── airgap/                # Bundle build + verify scripts
├── services/
│   ├── mcp-gateway/           # In-house MCP gateway (Go) — mTLS + OIDC + allowlist
│   ├── memory-service/        # Hierarchical memory + RAG (Python/FastAPI)
│   ├── approval-shim/         # gRPC adapter so every framework calls ClawSecure
│   └── policy-judge/          # Local Qwen3.5-4B classifier for routing/sensitivity
├── integrations/
│   ├── clawrails/             # Glue + config for the routing layer
│   ├── clawboot/              # Tier-aware seeds the installer consumes
│   └── clawsecure/            # ClawFirm-branded preset bundles
├── skills/                    # Sigstore-signed skill registry
├── docs/                      # Architecture, install, security, ADRs
├── scripts/                   # install.sh, airgap-build.sh, doctor.sh
└── benchmarks/                # Routing-decision quality + inference perf
```

---

## Tiers at a glance

| | **Solo** | **SMB** | **Enterprise** |
|---|---|---|---|
| Hardware floor | 8-core CPU, 16 GB RAM | 16 vCPU, 64 GB RAM, 1× 24 GB GPU | 3+ control planes, ≥2 GPU nodes |
| Orchestrator | Docker Compose | k3s (single binary) | k0s HA + Flux |
| Default sandbox | DifySandbox + Docker rootless | DifySandbox + gVisor | DifySandbox + gVisor + **Firecracker** |
| Secrets | OpenBao single-shard | OpenBao | OpenBao HA (Raft) |
| Identity | Authentik (single-user) | Authentik + SSO | Authentik + Keycloak federation |
| Default governance tier | Tier 0 (chat-only) | Tier 1 (smart-approval) | Tier 1 with Tier 2 enforced for write/egress |
| Air-gap support | n/a | manual bundle | first-class `clawfirm bundle` workflow |

See [`docs/install/`](./docs/install/) for the full per-tier install guides.

---

## Composition: how the three repos plug in

| Repo | ClawFirm role | Wire-up location |
|---|---|---|
| **ClawRails** | Inference router + cost layer. Sits in front of every model call. | `integrations/clawrails/` and `services/policy-judge/` adds the local sensitivity/complexity classifier. |
| **ClawBoot** | Bare-metal installer. Provides the cloud-init seed format we extend per-tier. | `integrations/clawboot/seeds/{solo,smb,enterprise}.seed.yaml` — the `clawfirm` CLI emits these. |
| **ClawSecure** | Approval queue + policy engine + audit timeline. | `integrations/clawsecure/presets/` ships ClawFirm-branded `balanced-default` and `strict-lockdown` bundles. `services/approval-shim/` adapts ClawSecure's REST API to gRPC for LangGraph / n8n / Dify / OpenHands. |

ClawFirm itself is the thin glue + opinionated assembly + tier-aware deployer. The intelligence lives in those three repos and the curated upstream components.

---

## Documentation

- [`docs/architecture/`](./docs/architecture/) — reference architecture, request routing, governance ladder, memory service.
- [`docs/install/`](./docs/install/) — tier-by-tier install, including air-gapped.
- [`docs/security/`](./docs/security/) — threat model, secure-default matrix, MCP gateway design.
- [`docs/operations/`](./docs/operations/) — upgrade, backup/restore, observability.
- [`docs/adr/`](./docs/adr/) — architecture decision records.
- [`LICENSING.md`](./LICENSING.md) — full per-component license matrix and redistribution implications.
- [`ROADMAP.md`](./ROADMAP.md) — 90-day MVP and 12-month plan.

---

## Design partners wanted

> [!NOTE]
> We're looking for a small group of design partners: AI consultancies, indie builders, SMB SaaS teams, and security-conscious operators who want to shape the Solo, SMB, and Enterprise tiers in public.
>
> If that's you, [open an issue](https://github.com/jimmershere/clawfirm/issues) with your use case or reach out through the repo discussions once enabled.

## License

ClawFirm itself is MIT-licensed (matching ClawRails, ClawBoot, ClawSecure). Bundled upstream components retain their own licenses — see [`LICENSING.md`](./LICENSING.md) for the full matrix and the implications for commercial redistribution. In particular, the n8n Sustainable Use License and the Dify Open Source License place restrictions on offering ClawFirm as a managed multi-tenant service to third parties; an "permissive-only" build profile that swaps both out is available.

---

## Status

**Alpha.** Solo tier installs and boots end-to-end. SMB and Enterprise tiers track the [90-day roadmap](./ROADMAP.md). Looking for design partners — open an issue on any of the four repos.
