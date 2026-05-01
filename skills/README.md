# Skills registry

ClawFirm uses Sigstore-signed skill packages for the EdgeClaw shell, Dify plugins, and OpenHands tools. Every skill in `registry/index.json` is verifiable via `cosign verify-blob` against its `.sig` and `.crt` files in `registry/signed/`.

## Default policy

- Signature required (`SkillsConfig.RequireSig = true` in seed.yaml)
- Mandatory dangerous-code scan on install
- Tier 0 / Tier 1 limited to `core/` skills (read-only filesystem, web search, git diff)
- `--dangerously-force-unsafe-install` removed entirely on Enterprise tier

## Adding a skill

1. Open a PR against this directory adding the skill to `registry/index.json`.
2. Sign the skill artifact with Sigstore (the `release.yml` workflow does this for first-party skills).
3. Reviewer verifies the signature and the dangerous-code scan output before merging.
