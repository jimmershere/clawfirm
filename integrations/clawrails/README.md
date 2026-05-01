# Integration: ClawRails

This directory contains ClawFirm's glue for [`jimmershere/clawrails`](https://github.com/jimmershere/clawrails) — the inference router and cost layer.

## Where it sits

ClawRails runs as a container in the ClawFirm stack (`deploy/solo/docker-compose.yml`, `deploy/smb/k3s/clawrails.yaml`). Every model call from every engine — Dify, n8n, LangGraph, OpenHands, the OpenClaw shell — points at `http://clawrails:4180/v1` as its OpenAI-compatible endpoint. ClawRails then decides:

1. Run the local **policy-judge** to classify the call.
2. Route to **Ollama** (small / single-user), **vLLM** (shared / GPU), or a **frontier API** (gated by ClawSecure policy).
3. Attribute cost to a `(job, customer, feature)` triple.
4. Apply the kill-switch if a budget is breached.
5. Emit OpenTelemetry GenAI spans to **Langfuse**.

## Files

- `config/routes.yaml.example` — the per-tier default routing rules.
- `hooks/pre-egress.js` — JS hook called before any frontier-API call. Posts to ClawSecure for policy check.
- `hooks/post-completion.js` — JS hook called after every completion. Posts cost + token metrics to ClawSecure for the audit timeline and to Langfuse for the trace.

## What ClawFirm extends in ClawRails

- Adds the local policy-judge as a routing input (sensitivity + complexity score).
- Wires hard kill-switches to ClawSecure rules so an operator can halt egress in one click.
- Pre-configured per-tier route presets matching the governance ladder.
