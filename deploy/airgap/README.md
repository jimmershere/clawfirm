# Air-gap bundle build & verify

ClawFirm Enterprise supports fully air-gapped install. The bundle contains every container image, model weight, signed skill, Helm chart, and the `clawfirm` CLI itself, with a top-level signed manifest.

## Build

```bash
./bundle.sh \
  --tier enterprise \
  --profile permissive-only \
  --models qwen3.5-4b,qwen3.6-27b,llama-3.3-70b-q4 \
  --skills core,fs,git,browser,db,search \
  --output ./clawfirm.airgap.tar.zst
```

## Verify

```bash
./verify.sh ./clawfirm.airgap.tar.zst
```

See [`docs/install/airgap.md`](../../docs/install/airgap.md) for the full workflow.
