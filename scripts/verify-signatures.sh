#!/bin/sh
# Verify the Sigstore signature on a ClawFirm release artifact.

set -eu
artifact="${1:-}"
[ -n "$artifact" ] || { echo "usage: $0 <artifact.tar.gz>" >&2; exit 2; }

base="${artifact%.tar.gz}"
sig="${artifact}.sig"
crt="${artifact}.crt"

[ -f "$sig" ] || { echo "missing signature: $sig" >&2; exit 1; }
[ -f "$crt" ] || { echo "missing certificate: $crt" >&2; exit 1; }

cosign verify-blob \
  --certificate "$crt" \
  --signature   "$sig" \
  --certificate-identity-regexp 'https://github.com/jimmershere/clawfirm/.github/workflows/release\.yml@.*' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com \
  "$artifact"
