# ADR-0002: k3s for SMB tier, k0s for Enterprise tier

- **Status:** Accepted
- **Date:** 2026-04-30
- **Deciders:** @jimmershere

## Context

The SMB and Enterprise tiers need a Kubernetes distribution. Options: k3s, k0s, MicroK8s, RKE2, Talos+k3s, OpenShift, vanilla kubeadm.

## Decision

- **SMB:** k3s. Smallest binary (~65 MB), embedded SQLite as the datastore by default, fewest dependencies on the host, easiest single-node-to-multi-node growth path.
- **Enterprise:** k0s. Zero host-OS dependencies (single binary, no systemd, no docker required), built-in HA via k0sctl, namespaces from day 1, FIPS mode available. Good fit for the appliance use case and for air-gap.

Talos Linux is supported as an opinionated host OS for an enterprise appliance image (Q3-26 roadmap), composed with k0s on top.

## Consequences

**Easier:**
- SMB operators get a single binary they can `curl | sh` and have a working cluster in minutes.
- Enterprise operators get a host-independent install they can drop on RHEL, Ubuntu, Debian, or a custom appliance OS without rework.
- Both are CNCF-friendly and don't lock anyone in.

**Harder:**
- Two K8s flavors to test against. Mitigation: deployment manifests are pure Kustomize/Helm, validated on both in CI.

## Alternatives considered

- **Vanilla kubeadm.** Rejected — too much operator burden for SMB.
- **MicroK8s.** Rejected — Snap dependency is awkward for non-Ubuntu hosts.
- **RKE2.** Considered — strong security defaults (CIS-aligned), and a future ADR may swap k0s → RKE2 for regulated-industry customers.
- **OpenShift.** Out of scope — too heavy and proprietary for the open-source ClawFirm distribution. Customers running OpenShift can deploy ClawFirm via Helm onto an existing cluster.
- **k3s for both tiers.** Rejected for Enterprise — k3s embedded etcd via SQLite doesn't scale for HA writes; switching k3s to external etcd negates its main advantage.
