# Architecture Decision Records

This directory contains the architecture decision records (ADRs) for ClawFirm. Each ADR captures one significant design choice, its context, the alternatives considered, and the consequences accepted.

When you make a significant architectural change, write a new ADR. When you supersede an existing ADR, write a new one that references it; do not edit the superseded one.

## Index

| # | Title | Status |
|---|---|---|
| [0000](./0000-template.md) | Template | n/a |
| [0001](./0001-three-tier-single-distribution.md) | Single distribution serves all three tiers | Accepted |
| [0002](./0002-k3s-for-smb-k0s-for-enterprise.md) | k3s for SMB, k0s for Enterprise | Accepted |
| [0003](./0003-tiered-sandbox.md) | Tiered sandbox (DifySandbox + gVisor + Firecracker) | Accepted |
| [0004](./0004-in-house-mcp-gateway.md) | In-house MCP gateway, not third-party SaaS | Accepted |
| [0005](./0005-reuse-sister-repos.md) | Reuse ClawRails, ClawBoot, ClawSecure as canonical components | Accepted |
| [0006](./0006-ollama-and-vllm.md) | Bundle both Ollama and vLLM behind a unified router | Accepted |
