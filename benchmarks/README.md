# Benchmarks

ClawFirm tracks two benchmark suites that drive routing thresholds and tier defaults.

- `routing/` — quality of the policy-judge classifier across rubrics (routing, sensitivity, complexity). Re-run quarterly when models change.
- `inference/` — throughput and latency of the Ollama vs vLLM backends across model sizes and concurrency levels.

Re-benchmark every 6 months at minimum; sooner when a major upstream model drop happens. Results inform the defaults in `internal/seed/seed.go` and the routing rules in `integrations/clawrails/config/routes.yaml.example`.
