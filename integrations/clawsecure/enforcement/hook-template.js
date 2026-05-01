// ClawSecure enforcement hook template.
//
// Drop into any Node engine's external-hook directory. Posts every
// candidate tool call to ClawSecure /api/evaluate and respects the response.

const CLAWSECURE_URL = process.env.CLAWSECURE_URL || "http://clawsecure:3188";

async function preToolCall({ tool, args, identity, source, governanceTier }) {
  const r = await fetch(`${CLAWSECURE_URL}/api/evaluate`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      tool,
      args,
      identity,
      source,
      governance_tier: governanceTier,
    }),
  });

  if (!r.ok) {
    // Fail closed: if ClawSecure is unreachable, block.
    throw new Error(`ClawSecure unreachable: ${r.status}`);
  }
  const decision = await r.json();
  if (decision.action === "block") {
    throw Object.assign(new Error(`Blocked by ClawSecure: ${decision.reason}`), {
      code: "CLAWFIRM_BLOCKED",
    });
  }
  if (decision.action === "require_approval") {
    // Caller is responsible for surfacing the approval prompt to the operator
    // and waiting for the queue ID to resolve.
    return { needsApproval: true, queueId: decision.queue_id };
  }
  return { needsApproval: false };
}

module.exports = { preToolCall };
