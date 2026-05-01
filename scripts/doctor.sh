#!/bin/sh
# Cross-platform health check. Run before/after `clawfirm up`.

set -eu

ok()   { printf '  \033[32mOK\033[0m   %s\n' "$1"; }
warn() { printf '  \033[33mWARN\033[0m %s\n' "$1"; }
fail() { printf '  \033[31mFAIL\033[0m %s\n' "$1"; failures=$((failures+1)); }

failures=0
echo "ClawFirm doctor:"

# Docker
if command -v docker >/dev/null 2>&1; then
  if docker info >/dev/null 2>&1; then ok "docker daemon reachable"
  else fail "docker installed but daemon not reachable"; fi
else fail "docker not installed"; fi

# Disk
free_gb=$(df -BG --output=avail "$HOME" 2>/dev/null | tail -n1 | tr -d 'G ' || echo 0)
if [ "${free_gb:-0}" -lt 50 ]; then
  warn "less than 50 GiB free in \$HOME (have ${free_gb} GiB) - models will eat space"
else ok "disk space (${free_gb} GiB free)"; fi

# Memory
mem_gb=$(awk '/MemTotal/ {printf "%.0f", $2/1024/1024}' /proc/meminfo 2>/dev/null || sysctl -n hw.memsize 2>/dev/null | awk '{printf "%.0f", $1/1024/1024/1024}')
if [ "${mem_gb:-0}" -lt 16 ]; then
  warn "less than 16 GiB RAM (have ${mem_gb} GiB) - Solo tier minimum"
else ok "RAM (${mem_gb} GiB)"; fi

# Ports
for p in 7878 3000 3188 5001 5678 8200 9000 11434; do
  if (echo > /dev/tcp/127.0.0.1/$p) >/dev/null 2>&1; then
    fail "port $p already in use - ClawFirm needs it"
  fi
done
[ "$failures" = "0" ] && ok "no required ports occupied"

# Cosign for signature verification
if command -v cosign >/dev/null 2>&1; then ok "cosign installed (signature verification available)"
else warn "cosign not installed - releases will install unverified"; fi

if [ "$failures" = "0" ]; then
  printf "\n\033[32mAll checks passed.\033[0m\n"
  exit 0
else
  printf "\n\033[31m%d check(s) failed.\033[0m\n" "$failures"
  exit 1
fi
