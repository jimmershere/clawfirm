# Roadmap

## 90-day MVP (Solo + SMB GA)

| Sprint | Deliverable |
|---|---|
| **Weeks 1-2** | Repo scaffold; Go `clawfirm` CLI; signed-release pipeline (Sigstore + cosign); Compose template for Solo. |
| **Weeks 3-4** | Bundle Ollama + Dify + n8n + Langfuse + pgvector + Authentik in Compose; first `clawfirm up` works on a fresh Ubuntu/macOS box. Tier-0 default enforced. |
| **Weeks 5-6** | OpenClaw/EdgeClaw shell on `127.0.0.1`, web UI only. **ClawSecure** integrated as the approval/audit backend. **ClawRails** integrated as the inference router. `policy-judge` Qwen3.5-4B classifier wired in. |
| **Weeks 7-8** | LangGraph runtime + LangChain glue. OpenHands V1 DockerWorkspace integration. **MCP gateway v1** — allowlist + audit + identity scoping. |
| **Weeks 9-10** | Memory Service unifying pgvector + EdgeClaw memory tree (event-sourced). Dify retriever plugin. n8n MCP server preconfigured. |
| **Weeks 11-12** | k3s install path for SMB tier. Coolify integration. Headscale + Tailscale automatic enrollment. First signed model + skill registry mirror. **ClawBoot** seeds wired into the `clawfirm init` flow. |
| **End of 90 days** | One-line installer. Solo & SMB tiers GA. 4-tier governance ladder. 5 channel adapters (web, CLI, Slack, Telegram, WhatsApp). Local + frontier routing via ClawRails. Langfuse traces. Air-gap bundle builder (alpha). |

---

## 12 months

### Q3-26 — Enterprise GA
- k0s HA installation path
- OpenBao HA cluster (Raft)
- Firecracker / Kata Containers as the **default** sandbox
- Flux GitOps for declarative configuration
- Air-gap installer GA (signed bundle workflow)
- FIPS-mode build profile
- SOC 2 Type II preparation
- Microsoft Agent Framework as alternative orchestrator backend
- Mastra TypeScript path for JS-shop customers

### Q4-26 — Confidential & marketplace
- Confidential-compute support (AMD SEV-SNP / Intel TDX) for Tier 3 autonomous workloads
- LiteBox / library-OS sandbox option (experimental)
- **Marketplace v2**: quality-scored, signed skills + Dify plugins + n8n nodes + MCP servers; dependency SBOMs; vulnerability scanner
- ClawSecure enforcement engine GA (today's status: visibility + triage; next: blocking enforcement)

### Q1-27 — Federation & policy UX
- Multi-region federation (geo-local Postgres + Qdrant; cross-region MCP-gateway routing)
- Zero-knowledge memory (encrypt embeddings on the user-controlled side)
- **Visual policy designer** for Tier 1 smart-approval rules (extends ClawSecure)
- Integration with ServiceNow / Jira / PagerDuty for Tier 3 ticket-gating

### Q2-27 — Mobile & on-device
- Native mobile shell (iOS + Android)
- On-device tiny model for offline chat (Qwen3.5-2B)
- Voice via Whisper (local) + Piper TTS

### Throughout
- Track upstream model releases (Qwen 3.7?, Llama 5?, Gemma 5?)
- Quarterly re-benchmarking of routing thresholds
- Monthly upstream upgrades (LangGraph, Dify, n8n, OpenHands, MAF)

---

## Open questions tracked as design risks

See [`docs/adr/`](./docs/adr/) for the architecture decision records that capture the design rationale, and the "Key open risks" section of `docs/architecture/overview.md` for the live risk log.
