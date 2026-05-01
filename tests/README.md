# Tests

- `e2e/` — end-to-end tests that bring up a full tier stack and exercise it via the public surface (CLI + HTTP + dashboards). Implemented in Bats.
- `integration/` — focused integration tests for individual services (MCP gateway round trip, memory-service hash chain, approval-shim adapter contract). Implemented in Go test or pytest.

Run e2e locally:

```bash
make solo-up
bats tests/e2e/solo-up-down.bats
```
