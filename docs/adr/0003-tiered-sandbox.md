# ADR-0003: Tiered sandbox (DifySandbox + gVisor + Firecracker)

- **Status:** Accepted
- **Date:** 2026-04-30
- **Deciders:** @jimmershere

## Context

Tool execution from agents must be sandboxed. Options range across containers (Docker rootless, Podman), user-space kernels (gVisor), microVMs (Firecracker, Kata, Cloud Hypervisor), library-OS approaches (LiteBox), and WebAssembly (Wasmtime, WasmEdge).

No single sandbox is right for every workload-and-tier combination. A tier-1 solo user on a laptop has different cost/security tradeoffs than a Tier-3 enterprise autonomous job.

## Decision

ClawFirm ships a **tiered sandbox**, with the active sandbox selected by the (tier × tool-class × ClawFirm-tier) triple:

| Use | Sandbox |
|---|---|
| Python / Node code execution inside a Dify-built workflow (any tier) | **DifySandbox** — seccomp whitelist, very low overhead, language-aware |
| Arbitrary tool execution at SMB tier | **gVisor (runsc)** — user-space kernel, ~50–100 ms cold start |
| Arbitrary tool execution + Tier 3 autonomous workloads at Enterprise | **Firecracker** microVM — hardware boundary via KVM, ~125 ms boot, <5 MiB memory overhead |
| OpenHands coding-agent runtime | **DockerWorkspace** (the OpenHands-bundled sandbox), upgradable to Firecracker via Kata for Enterprise |
| Solo on macOS | **DifySandbox + Docker Desktop's VZ** — gVisor and Firecracker don't run natively on macOS |

The MCP gateway is the single decision point: it knows the tier, the tool class, and the policy, and routes each tool call to the right sandbox.

## Consequences

**Easier:**
- Each tier gets a sandbox proportional to its threat model and its host capabilities.
- We get to inherit DifySandbox's mature Python/Node story for free.
- Tier 3 (autonomous) becomes safe to offer because of Firecracker's hardware boundary.

**Harder:**
- Three sandbox runtimes to maintain images for and test against. Mitigation: each has a stable upstream and a small, well-defined integration surface.
- Cold-start cost increases as you go up the tier ladder. Acceptable — Tier 3 is for batch / background jobs where 125 ms is invisible.

## Alternatives considered

- **One sandbox for all tiers.** Considered Firecracker only — too heavy for solo prosumer hardware. Considered gVisor only — insufficient boundary for Tier 3 enterprise autonomous jobs.
- **WebAssembly (Wasmtime / WasmEdge) as the universal sandbox.** Rejected for now — WASI capability set is still catching up; many tool integrations require POSIX. Future ADR may add WASM as an option for portable signed skills.
- **Kata Containers as the only microVM.** Considered — and we DO use Kata as the integration surface for Firecracker-on-k8s. The choice between Kata-with-Firecracker and Kata-with-Cloud-Hypervisor is a tunable.
- **E2B / Modal hosted sandboxes.** Rejected as the default — violates the "lowest cost to run" priority. Available as opt-in routing destinations via ClawRails for customers who want them.
