# Contributing

Thanks for being interested in ClawFirm.

## Quick orientation

ClawFirm is a meta-distribution. The intelligence lives in three sister repos and a curated upstream stack; this repo is the **glue** plus the **opinionated assembly**.

If you want to contribute:

- **CLI / installer / tier logic** → this repo, under `cmd/` and `internal/`.
- **Routing / cost / kill switches** → contribute to [`clawrails`](https://github.com/jimmershere/clawrails).
- **Bare-metal install / cloud-init / firstboot hardening** → contribute to [`clawboot`](https://github.com/jimmershere/clawboot).
- **Approvals / policy / audit / dashboard** → contribute to [`clawsecure`](https://github.com/jimmershere/clawsecure).
- **MCP gateway / memory service / approval shim / policy judge** → this repo, under `services/`.
- **Per-tier deployment manifests** → this repo, under `deploy/`.

If you're not sure where something belongs, open an issue here and we'll figure it out together.

---

## Development setup

```bash
git clone --recurse-submodules https://github.com/jimmershere/clawfirm.git
cd clawfirm

# Build the CLI
make build

# Bring up the solo stack against your dev checkout
make solo-up

# Run integration tests
make test
```

You'll want Go 1.22+, Python 3.12+, Docker 25+, and (optionally) k3s for SMB-tier work.

---

## Architecture decisions

Big design changes go through an ADR — see [`docs/adr/`](./docs/adr/) for the existing ones and `docs/adr/0000-template.md` for the format. Don't merge a PR that re-litigates an existing ADR; instead, write a new ADR that supersedes it and explains why.

---

## Coding style

- Go: `gofmt -s`, `golangci-lint`. Standard error wrapping (`fmt.Errorf("...: %w", err)`).
- Python: `ruff` + `black`. Type hints required on public functions.
- Shell: `shellcheck`. POSIX `sh` for installers; `bash` allowed for everything else.
- YAML: 2-space indent, no tabs.

CI enforces all of the above.

---

## Commit messages

Conventional Commits. Examples:

```
feat(mcp-gateway): add per-tool rate limiting
fix(memory-service): hash chain breaks on Postgres reconnect
docs(adr): supersede ADR 0003 with new sandbox tiering
chore(deps): bump LangGraph to 1.2.0
```

---

## License

By contributing, you agree your contributions are licensed under the MIT License (this repo). For contributions to ClawRails, ClawBoot, or ClawSecure, the same applies — they're all MIT.
