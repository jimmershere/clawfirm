# SMB tier — k3s deployment

Kustomize overlay for a 1-3 node k3s cluster. Driven by `clawfirm up` after `clawfirm init --tier smb`.

## Files

| File | Purpose | Status (90-day MVP) |
|---|---|---|
| `kustomization.yaml` | Top-level Kustomize manifest | ✅ |
| `namespace.yaml` | clawfirm + clawfirm-system namespaces with PSA labels | ✅ |
| `postgres.yaml` | Postgres + pgvector StatefulSet with multi-DB initdb | ✅ |
| `ollama.yaml` | Ollama deployment with GPU node selector | ✅ |
| `vllm.yaml` | vLLM deployment with GPU node selector | ✅ |
| Other service manifests | redis, openbao, authentik, dify, n8n, langgraph, openhands, clawrails, clawsecure, memory-service, mcp-gateway, approval-shim, policy-judge, edgeclaw, langfuse, headscale, ingress | 🚧 scaffold pointing to compose ref |

The scaffold manifests are placeholders during the MVP phase. The Solo tier `docker-compose.yml` is the canonical reference — every service has the same image, env vars, and inter-service wiring there. The `clawfirm` CLI's tier-aware installer will translate that into proper K8s manifests with HPA, pod anti-affinity, and service mesh policy as part of the Weeks 11-12 sprint of the [90-day roadmap](../../../ROADMAP.md).

## Apply manually

```bash
# Required: a k3s cluster with NVIDIA device plugin and at least one node
# labeled clawfirm.io/inference=true.

kubectl apply -k .

# Watch
kubectl -n clawfirm get pods --watch
```

## Helm path

For operators who prefer Helm over Kustomize, `../helm-values/` will contain values.yaml files for each upstream component. The `clawfirm` CLI generates them from the same source of truth.

## Sandbox

The SMB tier defaults to **gVisor** as the sandbox runtime for tool execution. Enable RuntimeClass on your k3s nodes:

```bash
# On each node:
sudo apt install runsc
echo 'runtime "runsc" { runtime_path = "/usr/local/bin/runsc" }' >> /etc/containerd/config.toml
sudo systemctl restart containerd
```

Then ClawFirm pods that execute untrusted code request the `gvisor` RuntimeClass automatically.
