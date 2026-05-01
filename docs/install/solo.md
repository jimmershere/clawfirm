# Solo tier install

For prosumers, indie developers, single-machine setups. Runs on a laptop or a single VPS.

## Hardware floor

- 8-core CPU
- 16 GB RAM
- 100 GB SSD
- Optional: 8–24 GB GPU. Apple Silicon M2/M3/M4 supported via Ollama (and vllm-mlx in Q3-26).

You can run on less, but expect slow first-token times on the chat model.

## Install

```bash
curl -fsSL https://get.clawfirm.io | sh
clawfirm init --tier solo
clawfirm up
```

`clawfirm init --tier solo` writes `~/.clawfirm/config.yaml`, generates self-signed certs for local TLS, generates an OpenBao root token (sealed), and pulls the default models (`qwen3.5-4b` for the policy-judge and chat fallback; `qwen3-coder-30b-a3b` if RAM permits).

`clawfirm up` brings up the Compose stack from `deploy/solo/docker-compose.yml`.

## What runs by default

- OpenClaw / EdgeClaw shell on `127.0.0.1:7878` (web UI only)
- ClawSecure dashboard on `127.0.0.1:3188`
- Langfuse on `127.0.0.1:3000`
- Authentik on `127.0.0.1:9000` (single-user mode)
- Dify, n8n, LangGraph runtime as backends to the shell
- Ollama on `127.0.0.1:11434`
- Postgres + pgvector for shared state
- OpenBao single-shard for secrets
- ClawRails as the inference router
- The MCP gateway, memory service, approval shim, and policy judge

All bound to `127.0.0.1`. To access from another device on your network, install Tailscale on both ends:

```bash
clawfirm net enable --provider tailscale
```

## What's disabled by default

- All chat-app channels (WhatsApp, Telegram, Slack, Discord, iMessage, Matrix, Teams)
- Frontier API egress (no API keys configured)
- Skill marketplace (only signed core skills)
- Tier 2 / Tier 3 governance — you start at Tier 0 (chat-only)

Promote when you're ready:

```bash
clawfirm tier promote 1   # Tier 1 = smart approval
```

## Hardware footprint at idle

Approximately 3.5 GB RAM, <5% CPU. Most of the RAM is the loaded Qwen3.5-4B model.

## Common operations

```bash
clawfirm status          # show all services + tier + governance
clawfirm logs <service>  # tail logs
clawfirm doctor          # health check
clawfirm upgrade         # pull latest signed images, verify, restart
clawfirm down            # stop everything (state preserved)
```

## Troubleshooting

If `clawfirm up` fails, run `clawfirm doctor` first — it checks Docker socket, port conflicts, disk space, and SELinux/AppArmor profiles.

If you're on macOS, ensure Docker Desktop is using the VZ framework (not the legacy hypervisor) for DifySandbox to work.

If you're on a Raspberry Pi 5 or similar ARM SBC, run `clawfirm init --tier solo --profile minimal` to skip Dify and n8n (they're heavy on memory). You'll lose visual builders but keep the shell, ClawSecure, ClawRails, and direct LLM chat.
