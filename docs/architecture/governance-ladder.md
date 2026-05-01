# Governance Ladder

ClawFirm enforces approvals via a four-tier ladder. The same ladder applies whether the request is being handled by Dify, n8n, LangGraph, OpenHands, Goose, or the OpenClaw shell — because every framework's native approval primitive is adapted to a single backend (`approval-shim` → ClawSecure REST).

| Tier | Behavior | Sandbox required | Recommended for |
|---|---|---|---|
| **0 — Chat-only** | No tools accessible. Pure local LLM Q&A. All channels disabled until paired. | n/a | First-run default. Untrusted environments. Demo / kiosk mode. |
| **1 — Smart approval** | Local **policy-judge** (Qwen3.5-4B) classifies each tool call as `safe` / `risky` / `dangerous`. `safe` (read-only file ops, `git diff`, search, `rg`, web search) auto-approved; `risky` and `dangerous` prompt the operator. | DifySandbox or stronger | Solo and SMB normal use. |
| **2 — Manual approval** | Every tool call prompts. Compatible with Slack / Discord / Telegram approval clients via the ClawSecure approval queue. | gVisor or stronger | High-stakes work. Enterprise default for any tool tagged `egress=true` or `write=true`. |
| **3 — Autonomous** | No prompts. **Only inside a fresh Firecracker microVM** with default-deny egress and a wall-clock budget. Explicit `--i-understand-the-risk` opt-in. Two-person rule for change tickets at Enterprise tier. | Firecracker | Background batch jobs in tightly scoped sandboxes. |

## Why a ladder, not a switch

The OpenClaw family ships with a binary "approve everything / approve nothing" model that pushes operators into autonomy because manual approval is too noisy. Smart approval (Tier 1) is the productive middle: a local 4-billion-parameter model reads the actual tool call (target, args, identity, recent context) and only prompts on the things a human would actually want to see.

## How it's enforced uniformly

Every framework expresses approvals differently:

- **OpenHands** has a `Conversation` confirmation policy (`always`, `auto-confirm`, `confirm-on-action`).
- **Goose** has approval modes (`autonomous`, `smart-approval`, `manual-approval`, `chat`).
- **LangGraph** uses `interrupt_before=[<node>]` and human-in-the-loop checkpointing.
- **n8n** has dedicated approval nodes (Wait, Approve).
- **Dify** has human-input nodes inside workflows.
- **The OpenClaw exec-approvals socket** has a per-skill `autoAllowSkills` model.

ClawFirm's `approval-shim` provides one gRPC interface (`ApprovalService`) and ships per-framework adapters that translate it into each framework's native primitive. ClawSecure is the single source of truth for approval queue, policy rules, and audit log.

## ClawSecure preset bundles

Three ClawFirm-branded preset bundles are installable in one command:

- **`clawfirm-solo`** — Tier 0 default; whitelists web search, fs read, git diff for Tier 1 promotion.
- **`clawfirm-smb`** — Tier 1 default; blocks all egress except allowlisted model providers; requires approval for any `write=true` tool.
- **`clawfirm-enterprise`** — Tier 1 default with Tier 2 enforcement on `egress=true` and `write=true`; Tier 3 only via signed change-control ticket.

Install:

```bash
curl -X POST http://localhost:3188/api/bundles/clawfirm-smb/apply
```

(See `integrations/clawsecure/presets/`.)

## Promoting between tiers

Tier promotion is a deliberate operator action. The web UI shows the current tier and what's blocked at this tier with a one-click "promote with explanation." The promotion is logged to ClawSecure with a hash chain entry.

Demotion is also a deliberate action — but it's an emergency action: clicking "panic → Tier 0" stops all in-flight tool calls, drains the approval queue with auto-deny, and locks the system to chat-only until an admin re-promotes.

## What the policy-judge actually looks at

The Qwen3.5-4B classifier is given:

1. The proposed tool name and arguments (redacted).
2. The user identity and current tier.
3. The last N messages of context (truncated to a few hundred tokens).
4. The list of allowlisted tools.
5. A short rubric of `safe` / `risky` / `dangerous` examples (in `services/policy-judge/policy_judge/prompts/`).

It returns JSON: `{class, confidence, reason}`. ClawFirm rejects responses with `confidence < 0.6` and escalates to manual approval rather than guessing.
