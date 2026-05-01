# Integration tests

Focused tests for a single service or contract:

- `mcp-gateway/` — round-trip a tool call through the gateway with a mock MCP server; verify allowlist + audit + token minting.
- `memory-service/` — hash chain integrity across crash/restart; projection rebuild from event log.
- `approval-shim/` — adapter contract for each framework (LangGraph, n8n, Dify, OpenHands).

Implemented as Go tests (`go test ./...`) or pytest as appropriate per service.
