# Licensing

ClawFirm is MIT-licensed. Bundled upstream components are not all MIT, and a few have meaningful restrictions on commercial redistribution that integrators must understand before shipping ClawFirm-as-a-service to third parties.

This document is informational, not legal advice. Get counsel before commercial redistribution.

---

## Component matrix

| Component | License | Redistribute commercially? | Implication |
|---|---|---|---|
| **ClawFirm CLI + glue + manifests** (this repo) | MIT | Yes | None. |
| **ClawRails** | MIT | Yes | None. |
| **ClawBoot** | MIT | Yes | None. |
| **ClawSecure** | MIT | Yes | None. |
| LangChain / LangGraph | MIT | Yes | None. |
| Microsoft Agent Framework | MIT | Yes | None. |
| OpenHands core SDK + agent-server | MIT | Yes | The `enterprise/` directory is source-available with paid license required after 1 month — **ClawFirm bundles only the MIT core**. |
| CrewAI | MIT | Yes | AMP enterprise tier is separate. |
| Goose (AAIF / Linux Foundation) | Apache-2.0 | Yes | None. |
| Mastra **core** | Apache-2.0 | Yes | The `ee/` directories (RBAC/SSO/ACL) are under a separate Mastra Enterprise License — **not bundled in default builds**. |
| Authentik | MIT | Yes | None. |
| Keycloak | Apache-2.0 | Yes | None. |
| Headscale | BSD-3 | Yes | None. |
| Tailscale **client** | BSD-3 | Yes | The Tailscale **control plane** is proprietary and is **not** bundled — we ship Headscale as the open control plane. |
| Caddy | Apache-2.0 | Yes | None. |
| Traefik | MIT | Yes | None. |
| Flux | Apache-2.0 | Yes | None. |
| ArgoCD | Apache-2.0 | Yes | None. |
| k3s / k0s / Talos Linux | Apache-2.0 | Yes | None. |
| Coolify | Apache-2.0 | Yes | None. |
| OpenBao | MPL-2.0 | Yes | Modifications to OpenBao itself must be disclosed (file-level copyleft); using it as a service does not trigger this. |
| pgvector | PostgreSQL License (BSD-style) | Yes | None. |
| Qdrant | Apache-2.0 | Yes | None. |
| LanceDB | Apache-2.0 | Yes | None. |
| Langfuse | Apache-2.0 | Yes | Some enterprise features (SSO/RBAC/audit) are commercial — **not bundled by default**. |
| Ollama | MIT | Yes | None. |
| vLLM | Apache-2.0 | Yes | None. |
| llama.cpp | MIT | Yes | None. |
| **n8n** (`n8n-io/n8n`) | **Sustainable Use License (SUL-1.0)** | **Conditional.** | Internal business use is allowed. **Hosting n8n as a paid service for end users, white-labeling it, or letting customers connect their own n8n credentials is not.** `*.ee.*` files require a separate n8n Enterprise License. ClawFirm self-hosted is fine; "ClawFirm Cloud" with n8n exposed to tenants is not, without a commercial agreement with n8n GmbH. |
| **Dify** (`langgenius/dify`) | **Dify Open Source License** (Apache-2.0 + extra clauses; not OSI-approved) | **Conditional.** | Self-hosted bundling is permitted. Repackaging Dify as a competing multi-tenant SaaS, or removing Dify branding from the `web/` frontend, requires permission from LangGenius. |
| **Arize Phoenix** | Elastic License v2 (ELv2) | Conditional. | "Providing Phoenix as a managed service" is restricted. ClawFirm uses Phoenix only as the **optional fallback** observability; **Langfuse is the primary**. |
| **HashiCorp Vault** | BSL 1.1 | **No** vs HashiCorp commercial offerings. | **Excluded** from ClawFirm. We use OpenBao instead. |
| **Zitadel v3+** | AGPL-3.0 | Yes, but copyleft on modifications. | **Not the default** — Authentik (MIT) is our default identity provider. Available as an optional swap for B2B multi-tenant. |
| OpenClaw / EdgeClaw upstream | Source-available, custom (verify per release) | Verify per release | Mark as "subject to upstream terms" in any redistribution; clearly attribute. |

---

## Build profiles

ClawFirm ships multiple build profiles to address the redistribution-restriction problem:

### `default` — full stack
Everything. Best UX, broadest features. Subject to the n8n SUL and Dify Open Source License conditions.

### `permissive-only` — for redistribution
Drops Dify and n8n. Substitutes:
- **Dify** → [Flowise](https://github.com/FlowiseAI/Flowise) (Apache-2.0) for the visual builder.
- **n8n** → [Activepieces](https://github.com/activepieces/activepieces) (MIT) for the workflow engine.

This profile is for ClawFirm operators who plan to ship a managed ClawFirm-as-a-service to third parties.

### `air-gap-fips` — for regulated environments
Adds FIPS-mode builds, removes any component without a clean upstream FIPS story. Tracks the Q3-26 roadmap milestone.

Select via:

```bash
clawfirm init --tier enterprise --profile permissive-only
```

---

## Per-image SBOMs

Every ClawFirm-published OCI image ships with an SPDX SBOM as an OCI artifact. Inspect with:

```bash
cosign download sbom ghcr.io/clawfirm/mcp-gateway:VERSION
```

The SBOM lists every transitive dependency and its license, so you can audit before deploying.

---

## Trademark notes

"OpenClaw," "EdgeClaw," "Dify," "n8n," "Tailscale," "Authentik," and other named upstream components are trademarks of their respective owners. ClawFirm's references to them are nominative use; ClawFirm is not affiliated with or endorsed by any of them.
