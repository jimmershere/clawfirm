# ADR-0006: Bundle both Ollama and vLLM behind a unified router

- **Status:** Accepted
- **Date:** 2026-04-30
- **Deciders:** @jimmershere

## Context

Local model serving has roughly three serious open options in 2026: Ollama, llama.cpp/llama-server, vLLM. Each is best at something different.

- **Ollama** — best ergonomics, hot model swap, OpenAI-compat API, native tool-calling. Plateaus at ~22 RPS on H100 because no continuous batching; FIFO queue.
- **llama.cpp** — lowest single-user latency on CPU and Apple Silicon. Same scaling limitation as Ollama for multi-user.
- **vLLM** — 3.3–35× higher RPS than Ollama at high concurrency on H100/A100 (Red Hat 2025 benchmark). PagedAttention, continuous batching, FlashAttention 4 in v0.17. Heavier install; CUDA/ROCm only.

ClawFirm has to serve solo prosumers (single-user, no GPU or one consumer GPU) AND SMB/Enterprise (multi-user, shared GPU). No single backend is right for both.

## Decision

**Bundle both.** Ollama for solo / single-user / model-swap; vLLM for SMB/Enterprise GPU concurrency. Expose both behind a unified OpenAI-compatible router (LiteLLM, used by ClawRails internally).

Tier defaults:

- **Solo:** Ollama only. Models: `qwen3.5-4b`, `qwen3-coder-30b-a3b` (if RAM permits). Apple Silicon: same.
- **SMB:** Ollama + vLLM both running. Ollama owns small models for the policy-judge and chat fallback. vLLM owns the shared chat model (e.g., Llama 3.3 70B Q4 or Qwen3.6-27B).
- **Enterprise:** vLLM with tensor parallelism owns shared inference. Ollama optional for one-off models.

The cost router (ClawRails) decides per-request which backend to hit.

## Consequences

**Easier:**
- Solo gets a 5-minute install with no GPU concerns.
- SMB/Enterprise gets the throughput they actually need.
- Operators don't have to choose — the router decides per request.

**Harder:**
- Two inference stacks to maintain images and configs for. Acceptable — both are standardized OCI images with stable upstreams.
- Memory pressure on SMB nodes that run both — mitigated by recommending separate model cache volumes and by the Solo profile shipping with Ollama only.

## Alternatives considered

- **Ollama only.** Rejected — RPS plateau is a real ceiling for Enterprise workloads.
- **vLLM only.** Rejected — too heavy for solo prosumers; CUDA/ROCm requirement excludes Apple Silicon.
- **TGI (Text Generation Inference).** Comparable to vLLM on a subset of models; ecosystem is smaller, install is heavier. Available as an opt-in alternative.
- **vllm-mlx for Apple Silicon Enterprise nodes.** Adopted in addition — the Q3-26 roadmap adds vllm-mlx as the default for Mac-class enterprise nodes.

## Cost router rule of thumb

EdgeClaw's empirical data suggests 60–80% of agent traffic can be served cheaply by local small models when a judge model classifies requests well. ClawFirm assumes the same and pins the policy-judge to Qwen3.5-4B as the cost-vs-quality sweet spot. We re-benchmark every 6 months.
