# ClawBoot overlays

Overlays are merged on top of a base seed (`seeds/<tier>.seed.yaml`) when a special install mode applies. Currently:

- `airgap.overlay.yaml` — applied when the install media contains a `clawfirm.airgap.tar.zst` bundle. Disables network-dependent bootstrap steps (Sigstore verification falls back to local Rekor; model pulls come from the bundle; no upstream skills registry).
