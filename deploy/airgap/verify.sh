#!/bin/sh
# Wrapper around `clawfirm bundle verify`. See README.md.
set -eu
exec clawfirm bundle verify "$@"
