# ClawSecure enforcement hooks

ClawSecure today provides excellent visibility, triage, and policy management. ClawFirm closes the gap to actual blocking enforcement by pre-checking with ClawSecure before tool execution.

## Hook pattern

Every engine that can call tools is configured with an external hook that posts to ClawSecure's `/api/evaluate` endpoint and respects the decision before executing. The pattern is identical regardless of engine; only the wiring differs.

- **n8n** — `EXTERNAL_HOOK_FILES=/etc/n8n/hooks/clawsecure-emit.js` (see Solo compose)
- **OpenHands** — `OPENHANDS_APPROVAL_BACKEND=grpc://approval-shim:50051`, which posts to ClawSecure
- **Dify** — code nodes use the ClawFirm Dify plugin that wraps every tool call with a ClawSecure check
- **LangGraph** — `interrupt_before` nodes call out to the approval-shim
- **OpenClaw / EdgeClaw** — uses `openclaw-emitter.js` from the upstream ClawSecure repo + an exec-approval wrapper

`hook-template.js` is the canonical template; per-engine adapters follow this pattern.
