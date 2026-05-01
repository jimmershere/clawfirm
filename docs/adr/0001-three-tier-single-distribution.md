# ADR-0001: Single distribution serves all three tiers (solo / SMB / enterprise)

- **Status:** Accepted
- **Date:** 2026-04-30
- **Deciders:** @jimmershere

## Context

ClawFirm needs to serve solo prosumers, SMB teams, and enterprise customers with HA + air-gap. The classic options are (a) one product with config-driven tiers, (b) three separate products, or (c) a freemium SaaS with a self-hosted "open core."

## Decision

One product, one repo, one CLI. Tier is selected at install time (`clawfirm init --tier {solo,smb,enterprise}`) and produces different deployment artifacts (Compose vs k3s vs k0s+Flux), different default sandboxes (DifySandbox+rootless vs +gVisor vs +Firecracker), and different default governance (Tier 0 vs Tier 1 vs Tier 1+ Tier 2 enforcement on write/egress).

The intelligence — ClawRails, ClawSecure, ClawBoot, the MCP gateway, the Memory Service, the policy-judge — is identical across tiers.

## Consequences

**Easier:**
- One codebase to maintain.
- Operators can grow with the product (Solo → SMB → Enterprise) without re-platforming.
- Upstream upgrades land for everyone simultaneously.

**Harder:**
- The CLI carries Compose, k3s, AND k0s drivers.
- Documentation has to be tier-aware.
- Test matrix is 3× larger (we plan to use a tier-tagged e2e suite).

**Accepted tradeoffs:**
- Solo tier carries a slight overhead (~50 MB for the unused k3s/k0s drivers in the binary) — acceptable.
- Enterprise customers who want a "stripped-down enterprise-only" build can produce one with `clawfirm build --profile enterprise-only`.

## Alternatives considered

- **Three separate products.** Rejected — three brand surfaces, three install paths, three docs sites. Operators who outgrow Solo would have to migrate to a different product.
- **Open-core SaaS.** Rejected — explicit user constraint (cost-first, secure-by-default, self-hosted) makes this the wrong shape.
- **One product, one tier (auto-detect everything).** Rejected — too much magic; operators legitimately need to declare intent (a 16 vCPU box could be solo-on-overkill or an SMB master node).
