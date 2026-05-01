# Request Routing Logic

When a user message arrives, ClawFirm classifies it, applies governance, and dispatches it to the right engine. This document walks through the path.

## Pseudocode

```text
on user_message(msg, channel, user):
  trace = langfuse.start_trace(user, channel)

  # 1. Identity & policy
  identity = Authentik.verify(user)
  policy   = ClawSecure.policy_for(identity)

  # 2. Local PII / sensitivity scan
  scan = LocalRedactor.scan(msg)            # Presidio + custom rules
  if scan.contains_secrets and policy.deny_secret_egress:
      msg = scan.redact()

  # 3. LLM-as-judge classification (Qwen3.5-4B local, ~50 ms on CPU)
  klass = PolicyJudge.classify(msg, history, scan)
  # klass ∈ {chat, dify_app, n8n_workflow, langgraph_agent, openhands_code, mcp_tool}

  # 4. Approval ladder enforcement
  if policy.tier == 0 and klass != "chat":
      return ask_user("Tier 0 = chat-only. Approve upgrade?")
  if klass in {n8n_workflow, openhands_code, mcp_tool} and policy.tier <= 2:
      ClawSecure.require_approval(...)

  # 5. Cost routing (ClawRails)
  model_choice = ClawRails.route(msg, klass, policy, sensitivity=scan.score)
  # → local-small | local-large | frontier-api (gated)

  # 6. Dispatch
  match klass:
    case "chat":
        return Inference.complete(msg, model=model_choice)
    case "dify_app":
        return Dify.run_app(matched_app, msg, model=model_choice)
    case "n8n_workflow":
        return n8n.exec(matched_wf, msg)              # via n8n MCP
    case "langgraph_agent":
        return LangGraph.invoke(matched_graph, msg, checkpoint=user_thread)
    case "openhands_code":
        return OpenHands.spawn(workspace="firecracker", task=msg)
    case "mcp_tool":
        return MCPGateway.dispatch(tool, args, identity)

  # 7. Memory write (event-sourced; both projections updated)
  MemoryService.append(user, msg, response, policy.memory_scope)
  trace.end(response)
```

## When does each engine run?

| Request shape | Routes to | Why |
|---|---|---|
| Pure chat / Q&A on local data | OpenClaw shell + local LLM | Cheapest, lowest latency |
| Visual / low-code agent app a non-developer built | **Dify** | Dify owns the prompt-IDE, RAG, plugin marketplace |
| Event-driven / cron / multi-SaaS-step workflow | **n8n** | 400+ nodes, native scheduler, mature webhook handling |
| Long-running, deterministic, multi-step reasoning needing checkpoints / human-in-the-loop / replay | **LangGraph 1.x** | First-class durable state and time-travel debugging |
| Coding / software-engineering task that needs a sandboxed dev environment | **OpenHands V1 SDK** with DockerWorkspace | Best SWE-Bench score, clean SDK |
| Anything calling tools | Through the **MCP gateway** | Universal contract; one place to apply policy |
| Untrusted code execution | **DifySandbox / gVisor / Firecracker** by tier | Defense in depth |

## PolicyJudge confidence

The Qwen3.5-4B classifier returns `{class, confidence, suggested_route, sensitivity}`. If `confidence < 0.6`, ClawFirm falls back to a curated regex+keyword classifier and prompts the user for clarification rather than guessing. This avoids the most common failure mode where small judge models hallucinate a confident but wrong route.

## Cost router decisions are cached

ClawRails caches routing decisions by SHA-256 of (msg, identity, policy version) for 5 minutes. This keeps repeated identical questions from re-running the judge model.
