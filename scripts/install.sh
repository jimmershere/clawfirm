#!/bin/sh
# ClawFirm one-line installer.
#
#   curl -fsSL https://get.clawfirm.io | sh
#
# Detects OS+arch, downloads the latest signed clawfirm CLI, verifies the
# Sigstore signature, and installs to /usr/local/bin.

set -eu

CLAWFIRM_VERSION="${CLAWFIRM_VERSION:-latest}"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
GITHUB_OWNER="jimmershere"
GITHUB_REPO="clawfirm"

err() { printf 'install.sh: %s\n' "$*" >&2; exit 1; }

# --- detect ---
os="$(uname -s | tr '[:upper:]' '[:lower:]')"
arch="$(uname -m)"
case "$arch" in
  x86_64|amd64) arch=amd64 ;;
  aarch64|arm64) arch=arm64 ;;
  *) err "unsupported architecture: $arch" ;;
esac
case "$os" in
  linux|darwin) ;;
  *) err "unsupported OS: $os" ;;
esac

# --- resolve version ---
if [ "$CLAWFIRM_VERSION" = "latest" ]; then
  CLAWFIRM_VERSION="$(curl -fsSL "https://api.github.com/repos/${GITHUB_OWNER}/${GITHUB_REPO}/releases/latest" \
    | grep '"tag_name"' | head -n1 | cut -d'"' -f4)"
fi
[ -n "$CLAWFIRM_VERSION" ] || err "could not resolve latest version"

base="https://github.com/${GITHUB_OWNER}/${GITHUB_REPO}/releases/download/${CLAWFIRM_VERSION}"
binary="clawfirm-${os}-${arch}"
url_bin="${base}/${binary}.tar.gz"
url_sig="${base}/${binary}.tar.gz.sig"
url_crt="${base}/${binary}.tar.gz.crt"

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

printf 'Downloading clawfirm %s for %s/%s ...\n' "$CLAWFIRM_VERSION" "$os" "$arch"
curl -fsSL "$url_bin" -o "$tmp/clawfirm.tar.gz"
curl -fsSL "$url_sig" -o "$tmp/clawfirm.tar.gz.sig" || err "signature download failed"
curl -fsSL "$url_crt" -o "$tmp/clawfirm.tar.gz.crt" || err "certificate download failed"

# --- verify ---
if command -v cosign >/dev/null 2>&1; then
  printf 'Verifying Sigstore signature ...\n'
  cosign verify-blob \
    --certificate "$tmp/clawfirm.tar.gz.crt" \
    --signature "$tmp/clawfirm.tar.gz.sig" \
    --certificate-identity-regexp "https://github.com/${GITHUB_OWNER}/${GITHUB_REPO}/.github/workflows/release\\.yml@.*" \
    --certificate-oidc-issuer https://token.actions.githubusercontent.com \
    "$tmp/clawfirm.tar.gz" || err "signature verification failed"
else
  printf 'WARNING: cosign not installed; skipping signature verification.\n' >&2
  printf '         Install cosign and re-run for verified install: https://docs.sigstore.dev/cosign/installation\n' >&2
fi

# --- install ---
tar -xzf "$tmp/clawfirm.tar.gz" -C "$tmp"
if [ -w "$INSTALL_DIR" ]; then
  mv "$tmp/clawfirm" "$INSTALL_DIR/clawfirm"
else
  printf 'Installing to %s requires sudo ...\n' "$INSTALL_DIR"
  sudo mv "$tmp/clawfirm" "$INSTALL_DIR/clawfirm"
fi
chmod +x "$INSTALL_DIR/clawfirm" 2>/dev/null || sudo chmod +x "$INSTALL_DIR/clawfirm"

printf '\nInstalled: %s\n' "$INSTALL_DIR/clawfirm"
"$INSTALL_DIR/clawfirm" --version || true
printf '\nNext: run `clawfirm init --tier solo` to set up.\n'
