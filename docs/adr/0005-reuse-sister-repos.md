# ADR-0005: Reuse ClawRails, ClawBoot, and ClawSecure as canonical components

- **Status:** Accepted
- **Date:** 2026-04-30
- **Deciders:** @jimmershere

## Context

Three pre-existing MIT-licensed repos owned by the same author already cover three of the most important ClawFirm subsystems:

- [`jimmershere/clawrails`](https://github.com/jimmershere/clawrails) — hybrid local/cloud routing, cost attribution, kill switches.
- [`jimmershere/clawboot`](https://github.com/jimmershere/clawboot) — bootable USB / cloud-init installer, hardened first-boot script, SSH lockdown, UFW, optional Nvidia driver presets.
- [`jimmershere/clawsecure`](https://github.com/jimmershere/clawsecure) — local-first Express dashboard, event timeline, approval queue, priority-ordered policy rules, artifact redaction, preset bundles, OpenClaw stdin emitter.

The original ClawFirm design sketched these three subsystems as new components ("LiteLLM + JudgeRouter," "Go installer CLI," "ApprovalService gRPC"). That would be redundant.

## Decision

Treat the three sister repos as **canonical components** of ClawFirm. Reference them as separate repos (git submodules + published container images), and limit this repo to:

- The thin glue / opinionated assembly.
- Tier-aware deployment manifests (`deploy/`).
- The four new services that don't exist elsewhere (MCP gateway, memory service, approval shim, policy judge).
- Integration adapters that wire the three sister repos into the broader stack (`integrations/`).
- Documentation, ADRs, install scripts, signing pipeline.

Naming: `ClawFirm = ClawBoot (install) + ClawRails (route) + ClawSecure (govern) + curated stack`.

## Consequences

**Easier:**
- Three subsystems already exist, are MIT-licensed, and have working code.
- Each of the three repos retains its independent identity, contributors, and release cadence.
- ClawFirm becomes much smaller in scope — primarily an integration project.
- License story stays clean (all four repos MIT).

**Harder:**
- Coordinated releases across four repos. Mitigation: a release-train cadence (monthly minor, weekly patch) and a `go.mod`-style version constraints file in this repo.
- Contributors have to know which repo to send a PR to. Mitigation: explicit guidance in `CONTRIBUTING.md`.
- API stability between repos becomes a contract. Mitigation: each sister repo publishes a versioned API surface, and ClawFirm pins to a specific version.

## Consequences for the new services

The four ClawFirm-original services (`services/mcp-gateway`, `services/memory-service`, `services/approval-shim`, `services/policy-judge`) exist precisely because no sister repo or upstream component covers them. They should be treated as candidates to spin out into their own repos if and when they outgrow the integration scope (a future ADR may do this for the MCP gateway in particular).

## Alternatives considered

- **Vendor (copy) the three repos into a monorepo.** Rejected — destroys their independent value and changes the contributor model for each.
- **Re-implement them inside ClawFirm.** Rejected — pure duplication of work.
- **Make ClawFirm itself a thin distribution that doesn't add new code.** Rejected — there are real gaps (MCP gateway, memory service, cross-framework approval shim, local policy judge) that no sister repo or upstream covers.
