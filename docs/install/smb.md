# SMB tier install

For small teams running on a single VPS or a 1–3 node mini-cluster.

## Hardware floor

Per node:

- 16 vCPU
- 64 GB RAM
- 500 GB SSD
- 1× 24 GB GPU (RTX 4090, L40S, A5000, or equivalent) on the inference node

Reserve roughly 2 vCPU / 8 GB RAM for ClawFirm overhead per node.

## Install on the first node

```bash
clawfirm init --tier smb --domain ai.acme.local
clawfirm up
```

`clawfirm init --tier smb` brings up k3s in single-node mode (embedded SQLite), installs the SMB Kustomize overlay from `deploy/smb/k3s/`, and starts the Coolify control plane on `:8000`.

## Add worker nodes

```bash
# On the master node:
clawfirm node token

# On each worker:
clawfirm node join <master-ip> <token>
```

## SSO

```bash
clawfirm sso configure --provider okta --client-id ... --client-secret ...
```

This wires Authentik as the front for every component (Dify, n8n, Langfuse, ClawSecure, the OpenClaw shell), federated to your IdP.

## Default governance

Tier 1 (smart approval). The policy-judge auto-approves `safe` tool calls; everything else prompts via the ClawSecure approval queue. Slack / Discord / Telegram approval clients can be configured for mobile approval workflows.

```bash
clawfirm approval channel add slack --webhook ...
```

## What runs

In addition to the Solo-tier services:

- vLLM with shared chat model (default: Qwen3.6-27B in BF16 or Llama 3.3 70B Q4)
- Ollama as fallback for small / specialty models and the policy-judge
- gVisor as the default sandbox for tool execution
- Headscale control plane on the master node
- Coolify for one-click app management of any extra services you want to run

## Backup

```bash
clawfirm backup create --target s3://backups/clawfirm/
```

Backs up Postgres (Dify + n8n + Langfuse + memory service), Qdrant snapshots if enabled, OpenBao seal config, ClawSecure SQLite, and your Authentik config.

## Upgrade

```bash
clawfirm upgrade
```

Performs a rolling upgrade. State migrations are auto-applied; non-reversible migrations require `--allow-irreversible`.

## Cost

A single GPU SMB node on a typical VPS provider runs ~$300–600/month all-in for inference + storage. Frontier-API costs are separately attributable per user / job / feature via ClawRails.
