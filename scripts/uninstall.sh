#!/bin/sh
# Remove the clawfirm CLI binary. Does NOT touch ~/.clawfirm/ data — pass
# --purge to also remove that.

set -eu
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
purge=0
[ "${1:-}" = "--purge" ] && purge=1

if [ -f "$INSTALL_DIR/clawfirm" ]; then
  rm -f "$INSTALL_DIR/clawfirm" 2>/dev/null || sudo rm -f "$INSTALL_DIR/clawfirm"
  echo "Removed $INSTALL_DIR/clawfirm"
fi

if [ "$purge" = "1" ]; then
  rm -rf "$HOME/.clawfirm"
  echo "Removed ~/.clawfirm"
fi
