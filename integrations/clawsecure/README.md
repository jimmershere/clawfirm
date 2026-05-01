# Integration: ClawSecure

This directory contains ClawFirm's glue for [`jimmershere/clawsecure`](https://github.com/jimmershere/clawsecure) — the local-first dashboard for events, approvals, policy rules, and artifact controls.

## Where it sits

ClawSecure runs as a container in the ClawFirm stack (`deploy/solo/docker-compose.yml`, `deploy/smb/k3s/clawsecure.yaml`). It is the **single source of truth for**:

- Event timeline (every tool call, every model call, every approval decision)
- Approval queue (driven by `services/approval-shim`)
- Policy rules (priority-ordered; higher priority wins; ties break by newest rule ID)
- Artifact redaction
- Preset bundles (`balanced-default`, `strict-lockdown`, plus ClawFirm-branded `clawfirm-solo` / `clawfirm-smb` / `clawfirm-enterprise`)

## Files

- `presets/clawfirm-solo.bundle.json`        — preset matching Solo tier defaults.
- `presets/clawfirm-smb.bundle.json`         — preset matching SMB tier defaults.
- `presets/clawfirm-enterprise.bundle.json`  — preset matching Enterprise tier defaults.
- `enforcement/hook-template.js`             — emitter template that other engines (n8n, OpenHands, etc.) use to post events to ClawSecure.

## Apply a preset

```bash
curl -X POST http://localhost:3188/api/bundles/clawfirm-smb/apply
```

## What ClawFirm extends in ClawSecure

ClawSecure today is excellent at visibility, triage, and policy management. ClawFirm closes the gap to **enforcement**:

- The MCP gateway, approval shim, and ClawRails all post to ClawSecure's REST API and respect its decisions BEFORE executing.
- The hash-chained audit log becomes mandatory at SMB+ tiers.
- The `openclaw-emitter.js` pattern is generalized into a unified MCP-gateway event sink so every tool call (not just OpenClaw skills) flows through the same timeline.
