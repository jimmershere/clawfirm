#!/bin/sh
# Wrapper around `clawfirm bundle build`. See README.md.
set -eu
exec clawfirm bundle build "$@"
