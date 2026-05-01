# Enterprise tier — Flux GitOps

ClawFirm Enterprise uses Flux for declarative configuration. The `clawfirm init --tier enterprise --gitops flux` command emits a starter GitOps repository structure that you commit to a Git server reachable from the cluster (a private GitHub repo, Gitea, internal GitLab, etc.).

## Bootstrap

```bash
flux bootstrap git \
  --url=ssh://git@github.com/your-org/clawfirm-cluster.git \
  --branch=main \
  --path=clusters/clawfirm-prod
```

`bootstrap.yaml` is the seed configuration that Flux applies first; it pulls the rest of the stack via `Kustomization` resources.
