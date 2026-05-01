# Enterprise tier — Helm values

Per-component `values.yaml` overlays for the upstream Helm charts ClawFirm Enterprise composes:

- Authentik (operator-managed at scale)
- OpenBao (HA Raft)
- Postgres (CloudNative-PG)
- Qdrant (clustered)
- Headscale
- Langfuse Enterprise
- The four ClawFirm-original services as Helm charts

The `clawfirm` CLI generates a complete set of values.yaml files for your topology via:

```bash
clawfirm helm values --tier enterprise --domain ai.acme.internal --output ./helm-values/
```

These get committed into the Flux repo so Flux applies them declaratively. The MVP scaffold ships with one example per component; production charts land in the Q3-26 milestone.
