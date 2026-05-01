# Enterprise tier install

For HA, RBAC, and air-gap-capable deployments. ClawFirm Enterprise is designed to run on customer infrastructure (your own hardware, your own k0s cluster, your own VPC) with no external SaaS dependencies.

## Hardware floor

- **3+ control-plane nodes** (k0s HA) — small/medium VMs, 4 vCPU / 16 GB each is sufficient
- **≥2 GPU nodes** (A100, H100, MI300X, or equivalent)
- **PB-class S3-compatible object storage** (MinIO, Ceph, AWS S3, etc.) for backups and air-gap bundles
- **Dedicated Postgres HA cluster** (Patroni or CloudNative-PG)
- Optional: **Confidential-compute hosts** (AMD SEV-SNP / Intel TDX) for Tier 3 autonomous workloads (Q4-26)

## Online install

```bash
clawfirm init --tier enterprise --topology ha --gitops flux
```

This generates:

- A `k0sctl.yaml` for the control plane
- A Flux GitOps repository structure for declarative configuration
- Helm value files for every bundled service
- An Authentik config federated with your existing IdP
- An OpenBao HA cluster spec (Raft storage, 3 nodes)

You commit the generated repo, point Flux at it, and wait. Flux applies, Authentik comes up, OpenBao auto-unseals via your KMS, and the rest of the stack provisions in dependency order.

## Air-gapped install

```bash
# On an internet-connected staging machine:
clawfirm bundle build \
  --tier enterprise \
  --models qwen3.5-4b,qwen3.6-27b,llama-3.3-70b \
  --skills core,fs,git,browser,db,search \
  --output clawfirm.airgap.tar.zst

# Transfer media (data-diode / sneakernet) — typically 30–80 GB

# On the target air-gapped network:
clawfirm install --from-bundle clawfirm.airgap.tar.zst
```

The bundle includes:

- All container images (signed)
- All model weights (digest-pinned)
- All skills, plugins, MCP servers (signed)
- The k0sctl + Flux repo
- A signed manifest

Verify before installing:

```bash
clawfirm bundle verify clawfirm.airgap.tar.zst
```

## Default sandbox

Firecracker microVMs via Kata Containers. Every untrusted tool call gets a fresh microVM with default-deny egress and a wall-clock budget.

## Default governance

Tier 1 (smart approval), with **Tier 2 enforced** for any tool tagged `egress=true` or `write=true`. Tier 3 only via signed change-control ticket integrated with ServiceNow / Jira.

## SSO and federation

Authentik is the SSO front for all components. Federation:

- OIDC / SAML to your existing IdP (Okta, Azure AD, Google Workspace, etc.)
- SCIM provisioning supported
- Optional Keycloak federation for very large user bases (>100k)

## Air-gapped model updates

Detached Sigstore signatures. Operators bring updates via approved cross-domain solution (CDS). ClawFirm verifies signatures before applying.

```bash
clawfirm update apply --bundle quarterly-update-q3-2026.tar.zst --verify
```

## Compliance posture

- SOC 2 Type II preparation included as part of the Enterprise tier (Q3-26)
- FIPS-mode build profile (`--profile air-gap-fips`) available Q3-26
- Audit log is hash-chained, exportable, and retained per your policy
- All cross-component traffic stays on the Tailscale/Headscale overlay (or pure WireGuard for strict air-gap)

## Backup and DR

OpenBao Raft snapshots, Postgres logical+physical backups, and Memory Service event-log snapshots all flow into your S3-compatible object store on a configurable cadence. DR rehearsals are scripted via `clawfirm dr rehearse` (creates an ephemeral parallel cluster from latest backups; exercises memory recall and tool-call replay).
