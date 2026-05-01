# Integration: ClawBoot

This directory contains ClawFirm's glue for [`jimmershere/clawboot`](https://github.com/jimmershere/clawboot) — the bootable USB / cloud-init installer.

## Where it sits

ClawBoot remasters Ubuntu Server 24.04 into a hands-free autoinstall ISO and runs a hardened first-boot script (`install.sh`) that brings up OpenClaw, SSH lockdown, UFW, and (optionally) NVIDIA drivers. ClawFirm extends this with **tier-aware seed files** that the first-boot script consumes to decide which Compose / k3s / k0s profile to bring up.

## Files

- `seeds/solo.seed.yaml`        — Solo tier seed: bring up Compose stack, Tier 0 governance.
- `seeds/smb.seed.yaml`         — SMB tier seed: bring up k3s, Tier 1 smart-approval.
- `seeds/enterprise.seed.yaml`  — Enterprise tier seed: deferred install pending k0sctl topology.
- `overlays/airgap.overlay.yaml` — overlay applied on top of any tier seed when the install media is an air-gap bundle.

## How it's wired

The ClawBoot first-boot script reads `/etc/clawfirm/seed.yaml` (placed by the autoinstall step). If present, it shells out to `clawfirm install --from-seed /etc/clawfirm/seed.yaml`.

ClawFirm's CLI generates these seeds via `clawfirm init --tier <t> --emit-seed > seed.yaml` so the same seed format is used by both the bare-metal first boot and the in-cluster CLI.

## What ClawFirm extends in ClawBoot

- Tier-aware seed format consumed by the first-boot script.
- Air-gap bundle ingestion mode.
- Sigstore-verified post-boot pull for skills, plugins, model digests.
- Cloud-init meta-data overlay for fleet provisioning via Ansible / Terraform / MAAS.
