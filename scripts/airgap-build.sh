#!/bin/sh
# Build a self-contained air-gap bundle for ClawFirm.
#
#   ./scripts/airgap-build.sh \
#       --tier enterprise \
#       --profile permissive-only \
#       --models qwen3.5-4b,qwen3.6-27b,llama-3.3-70b-q4 \
#       --skills core,fs,git,browser,db,search \
#       --output ./clawfirm.airgap.tar.zst
#
# Wraps `clawfirm bundle build` and produces a signed bundle.

set -eu

if ! command -v clawfirm >/dev/null 2>&1; then
  echo "clawfirm CLI not on PATH; install via scripts/install.sh first" >&2
  exit 1
fi

clawfirm bundle build "$@"
