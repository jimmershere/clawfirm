# Solo tier deployment

Single-node Docker Compose stack. Driven by `clawfirm up` / `clawfirm down`.

## Files

- `docker-compose.yml` — the main stack.
- `docker-compose.gpu.yml` — NVIDIA GPU override + optional vLLM service.
- `Caddyfile` — local TLS + path-based routing for the dashboards.
- `.env.example` — template; `clawfirm init --tier solo` writes a real `.env`.
- `initdb/` — Postgres bootstrap SQL.

## Bring it up manually

```bash
cp .env.example .env
# edit .env to set strong random secrets, OR run:
clawfirm secrets generate --tier solo --output .env

docker compose up -d
```

For GPU:

```bash
docker compose -f docker-compose.yml -f docker-compose.gpu.yml up -d
```

## Service map

| Service | Internal port | Local port | Purpose |
|---|---|---|---|
| caddy | 80, 443 | 80, 443 | Ingress, TLS termination |
| edgeclaw | 7878 | 7878 | Assistant shell |
| dify-web | 3000 | 5001 | Visual agent builder |
| n8n | 5678 | 5678 | Workflow engine |
| langfuse | 3000 | 3000 | Observability |
| clawsecure | 3188 | 3188 | Approvals + audit |
| authentik-server | 9000 | 9000 | SSO |
| openbao | 8200 | 8200 | Secrets |
| ollama | 11434 | 11434 | LLM serving |
| clawrails | 4180 | 4180 | Inference router |
| mcp-gateway | 8443 | 8443 | Tool bus |

All ports listed as "Local port" are bound to `127.0.0.1` only.
