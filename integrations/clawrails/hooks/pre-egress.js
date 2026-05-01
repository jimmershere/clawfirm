// pre-egress hook — runs before any frontier-API call.
// Posts to ClawSecure for a policy check; aborts the call if rejected.
//
// Loaded by ClawRails via:
//   RAILS_HOOK_PRE_EGRESS=/etc/clawrails/hooks/pre-egress.js

module.exports = async function preEgress({ model, identity, request, sensitivity }) {
  const decision = await fetch(`${process.env.RAILS_CLAWSECURE_URL}/api/evaluate`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      kind: "frontier_egress",
      model,
      identity,
      sensitivity,
      tokens_in: request.estimated_tokens || 0,
    }),
  }).then((r) => r.json());

  if (decision.action === "block") {
    const err = new Error(`ClawSecure blocked frontier egress: ${decision.reason}`);
    err.code = "CLAWFIRM_EGRESS_BLOCKED";
    throw err;
  }
  if (decision.action === "downgrade") {
    return { override_model: decision.fallback_model };
  }
  return { ok: true };
};
