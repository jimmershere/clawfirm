// post-completion hook — runs after every model call.
// Posts cost + tokens to ClawSecure for the audit timeline.

module.exports = async function postCompletion({ model, identity, usage, latency_ms, route }) {
  const event = {
    kind: "completion",
    severity: "info",
    source: "clawrails",
    model,
    identity,
    route,
    tokens_in: usage.prompt_tokens,
    tokens_out: usage.completion_tokens,
    cost_usd: usage.cost_usd ?? 0,
    latency_ms,
  };
  await fetch(`${process.env.RAILS_CLAWSECURE_URL}/api/events`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(event),
  }).catch((e) => console.error("clawsecure post failed:", e.message));
};
