# Air-gapped install

Air-gapped ClawFirm runs with zero external network access. Useful for regulated industries, classified environments, and customers with strict data-sovereignty requirements.

## Overview

Air-gapped install is a two-step process:

1. **On a staging machine with internet access**, build a self-contained bundle.
2. **Transfer the bundle** to the target environment via approved means (data diode, sneakernet, cross-domain solution).
3. **On the target**, verify and install.

## Build the bundle

```bash
clawfirm bundle build \
  --tier enterprise \
  --profile permissive-only \
  --models qwen3.5-4b,qwen3.6-27b,llama-3.3-70b-q4 \
  --skills core,fs,git,browser-readonly,db,search \
  --output ./clawfirm.airgap.tar.zst
```

Flags:

- `--tier {solo,smb,enterprise}` — selects the orchestration backend.
- `--profile {default,permissive-only,air-gap-fips}` — selects which upstream components are bundled. `air-gap-fips` (Q3-26) drops anything without a clean FIPS story.
- `--models <list>` — model digests to embed. ClawFirm pins to specific digests, not floating tags.
- `--skills <list>` — skills + plugins + MCP servers to embed. Only signed entries are included.
- `--output <path>` — the bundle file. Typically 30–80 GB depending on model selection.

The bundle is a Zstandard-compressed tar containing:

```
clawfirm.airgap/
├── manifest.json              # SHA-256 of every artifact + signatures
├── manifest.json.sig          # Sigstore signature of the manifest
├── images/                    # OCI image layers (deduplicated)
├── models/                    # Model weights, digest-named
├── skills/                    # Signed skill packages
├── helm/                      # Helm charts
├── k0sctl.yaml.template       # Cluster topology template
├── flux/                      # Flux GitOps repo seed
└── installer/                 # The clawfirm CLI binary itself, signed
```

## Verify the bundle

Before transferring, verify on the staging machine:

```bash
clawfirm bundle verify ./clawfirm.airgap.tar.zst
```

This runs `cosign verify-blob` against the manifest and every signed artifact within. The verification is reproducible; you can also run it after transfer to confirm bit-for-bit integrity.

## Transfer

How you transfer is up to your environment:

- USB drive, sneakernet
- One-way data diode (the bundle is small enough for most CDS throughput)
- Customer-approved CDS

## Install on the target

```bash
sha256sum -c clawfirm.airgap.tar.zst.sha256       # local integrity
clawfirm bundle verify clawfirm.airgap.tar.zst    # signature integrity
clawfirm install --from-bundle clawfirm.airgap.tar.zst
```

The installer:

1. Imports OCI images into the local registry (or k0s' embedded one).
2. Imports model weights to a digest-addressable model store.
3. Bootstraps the k0s control plane (or k3s for SMB tier).
4. Applies Helm charts.
5. Seeds Authentik with a single bootstrap admin user; rotate immediately.

## Updates

Air-gapped updates use the same bundle mechanism:

```bash
# Build on the staging machine:
clawfirm bundle build --update --from-version 1.4.2 --to-version 1.5.0 --output update-1.5.0.tar.zst

# Transfer, then on the target:
clawfirm update apply --bundle update-1.5.0.tar.zst --verify
```

Update bundles are delta-encoded against the previous version where possible, so they're typically much smaller than full bundles.

## What does NOT work air-gapped

- Frontier API egress (by definition — no internet).
- Pulling new models / skills directly from upstream (must come via bundle).
- Telemetry to vendor (already off by default; impossible air-gapped).
- Sigstore Rekor public log (Enterprise tier publishes hash-chain heads to a local Rekor instance instead).

## What DOES work air-gapped

- Everything else. Local LLMs, all engines (Dify, n8n, LangGraph, OpenHands), full RAG, MCP gateway, ClawSecure approvals, ClawRails routing (between local backends), full audit, full backup/restore, full HA.
