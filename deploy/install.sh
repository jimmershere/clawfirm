#!/usr/bin/env bash
set -Eeuo pipefail

LOG_FILE="/var/log/clawboot-firstboot.log"
STATE_DIR="/var/lib/clawboot"
CONFIG_DIR="/etc/clawboot"
CONFIG_FILE="$CONFIG_DIR/clawboot.env"
OPENCLAW_USER="${SUDO_USER:-${USER:-clawboot}}"
INSTALL_OPENCLAW="true"
INSTALL_NANOCLAW="false"
INSTALL_CLAWFIRM="true"
CLAWFIRM_TIER="solo"
CLAWFIRM_REPO="https://github.com/jimmershere/clawfirm.git"
TELEGRAM_TOKEN=""
NONINTERACTIVE="${NONINTERACTIVE:-0}"
NODE_MAJOR="22"

mkdir -p "$STATE_DIR" "$CONFIG_DIR"
touch "$LOG_FILE"
exec > >(tee -a "$LOG_FILE") 2>&1

log() {
  printf '[%s] %s\n' "$(date '+%Y-%m-%d %H:%M:%S')" "$*"
}

fail() {
  log "ERROR: $*"
  exit 1
}

require_root() {
  if [[ ${EUID:-$(id -u)} -ne 0 ]]; then
    fail "Run as root."
  fi
}

load_config() {
  if [[ -f "$CONFIG_FILE" ]]; then
    # shellcheck disable=SC1090
    source "$CONFIG_FILE"
    INSTALL_OPENCLAW="${INSTALL_OPENCLAW:-$INSTALL_OPENCLAW}"
    INSTALL_NANOCLAW="${INSTALL_NANOCLAW:-$INSTALL_NANOCLAW}"
    TELEGRAM_TOKEN="${TELEGRAM_TOKEN:-$TELEGRAM_TOKEN}"
    OPENCLAW_USER="${OPENCLAW_USER:-$OPENCLAW_USER}"
    NONINTERACTIVE="${NONINTERACTIVE:-$NONINTERACTIVE}"
  fi
}

is_interactive() {
  [[ "$NONINTERACTIVE" != "1" && -t 0 && -t 1 ]]
}

prompt_choices() {
  if ! is_interactive; then
    log "Noninteractive mode detected; using config/defaults."
    return 0
  fi

  read -r -p "Install OpenClaw? [Y/n] " ans_openclaw || true
  case "${ans_openclaw:-Y}" in
    n|N) INSTALL_OPENCLAW="false" ;;
    *)   INSTALL_OPENCLAW="true" ;;
  esac

  read -r -p "Install NanoClaw scaffold too? [y/N] " ans_nanoclaw || true
  case "${ans_nanoclaw:-N}" in
    y|Y) INSTALL_NANOCLAW="true" ;;
    *)   INSTALL_NANOCLAW="false" ;;
  esac

  if [[ -z "$TELEGRAM_TOKEN" ]]; then
    read -r -p "Telegram bot token (optional, press Enter to skip): " TELEGRAM_TOKEN || true
  fi
}

apt_install() {
  export DEBIAN_FRONTEND=noninteractive
  apt-get update
  apt-get install -y \
    ca-certificates curl gnupg lsb-release software-properties-common \
    git jq unzip build-essential python3 python3-pip htop pciutils
}

install_node() {
  if command -v node >/dev/null 2>&1 && command -v npm >/dev/null 2>&1; then
    log "Node already installed: $(node -v), npm: $(npm -v)"
    return 0
  fi

  log "Installing Node.js ${NODE_MAJOR}.x"
  mkdir -p /etc/apt/keyrings

  if curl -4 -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key \
    | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg; then
    echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_${NODE_MAJOR}.x nodistro main" \
      > /etc/apt/sources.list.d/nodesource.list
    apt-get update || true
    apt-get install -y nodejs npm || true
  fi

  if ! command -v node >/dev/null 2>&1 || ! command -v npm >/dev/null 2>&1; then
    log "Nodesource path incomplete; falling back to direct Node tarball"
    local node_tar="/tmp/node-v${NODE_MAJOR}.14.0-linux-x64.tar.xz"
    curl -4 -fsSL "https://nodejs.org/dist/v${NODE_MAJOR}.14.0/node-v${NODE_MAJOR}.14.0-linux-x64.tar.xz" -o "$node_tar"
    tar -xf "$node_tar" -C /tmp
    cp -r "/tmp/node-v${NODE_MAJOR}.14.0-linux-x64/bin/"* /usr/local/bin/
    cp -r "/tmp/node-v${NODE_MAJOR}.14.0-linux-x64/lib/"* /usr/local/lib/
    cp -r "/tmp/node-v${NODE_MAJOR}.14.0-linux-x64/include/"* /usr/local/include/ 2>/dev/null || true
    hash -r
  fi

  command -v node >/dev/null 2>&1 || fail "Node install failed"
  command -v npm >/dev/null 2>&1 || fail "npm install failed"
  log "Installed Node: $(node -v), npm: $(npm -v)"
}

detect_gpu() {
  log "Detecting GPUs"
  # lspci lives in /usr/bin/; export full PATH for sudo environments
  export PATH="/usr/bin:/usr/sbin:/bin:/sbin:$PATH"

  if command -v lspci >/dev/null 2>&1; then
    lspci | grep -Ei 'vga|3d|display' || true
  else
    log "lspci not found — skipping GPU detection"
    return 0
  fi

  if lspci | grep -qi nvidia; then
    log "NVIDIA GPU detected"
    install_nvidia_stack_placeholder
  else
    log "No NVIDIA GPU detected"
  fi
}

install_nvidia_stack_placeholder() {
  log "Installing ubuntu-drivers-common for recommended NVIDIA driver flow"
  apt-get install -y ubuntu-drivers-common || true

  if command -v ubuntu-drivers >/dev/null 2>&1; then
    ubuntu-drivers devices || true
    # Draft behavior only: do not force driver install blindly on every box yet.
    log "Driver auto-install intentionally deferred in draft. Use: ubuntu-drivers autoinstall"
  fi

  if apt-cache show nvidia-utils-535 >/dev/null 2>&1; then
    apt-get install -y nvidia-utils-535 || true
  fi

  if command -v nvidia-smi >/dev/null 2>&1; then
    nvidia-smi || true
  else
    log "nvidia-smi not present yet; expected on some fresh installs until driver install/reboot"
  fi
}

install_openclaw() {
  [[ "$INSTALL_OPENCLAW" == "true" ]] || { log "Skipping OpenClaw install"; return 0; }

  log "Installing OpenClaw"
  npm install -g openclaw
  log "OpenClaw version: $(openclaw --version 2>/dev/null || echo unknown)"
}

install_nanoclaw_scaffold() {
  [[ "$INSTALL_NANOCLAW" == "true" ]] || { log "Skipping NanoClaw scaffold"; return 0; }

  local target_dir="/opt/nanoclaw"
  log "Preparing NanoClaw scaffold at $target_dir"
  mkdir -p "$target_dir"
  cat > "$target_dir/README.txt" <<'EOF'
NanoClaw scaffold placeholder.

Next steps (draft):
1. Clone NanoClaw repo here
2. Install container/runtime dependencies
3. Configure models, skills, and swarm settings
EOF
}

configure_telegram_token() {
  if [[ -z "$TELEGRAM_TOKEN" ]]; then
    log "No Telegram token provided. Skipping token setup."
    return 0
  fi

  mkdir -p /home/"$OPENCLAW_USER"/.openclaw/secrets || true
  printf '%s\n' "$TELEGRAM_TOKEN" > /home/"$OPENCLAW_USER"/.openclaw/secrets/telegram-bot-token.txt
  chown -R "$OPENCLAW_USER":"$OPENCLAW_USER" /home/"$OPENCLAW_USER"/.openclaw || true
  chmod 600 /home/"$OPENCLAW_USER"/.openclaw/secrets/telegram-bot-token.txt || true
  log "Telegram token file written for $OPENCLAW_USER"
}

write_status_note() {
  cat > /etc/motd <<'EOF'
ClawBoot first-boot provisioning completed.

Logs:
  /var/log/clawboot-firstboot.log

If Telegram token was skipped, add it later to:
  ~/.openclaw/secrets/telegram-bot-token.txt

Recommended next steps:
  openclaw status
  openclaw gateway status
EOF
}

install_clawfirm() {
  [[ "$INSTALL_CLAWFIRM" == "true" ]] || { log "Skipping ClawFirm install"; return 0; }
  log "Installing ClawFirm ${CLAWFIRM_TIER} tier..."

  # Install Docker if not present
  if ! command -v docker &>/dev/null; then
    log "Installing Docker..."
    curl -fsSL https://get.docker.com | sh
    usermod -aG docker "$OPENCLAW_USER" 2>/dev/null || true
  fi

  # Clone repo
  CLAWFIRM_DIR="/opt/clawfirm"
  if [[ -d "$CLAWFIRM_DIR" ]]; then
    git -C "$CLAWFIRM_DIR" pull --ff-only || true
  else
    git clone --depth=1 "$CLAWFIRM_REPO" "$CLAWFIRM_DIR"
  fi

  # Generate .env from example
  cd "$CLAWFIRM_DIR/deploy/$CLAWFIRM_TIER"
  if [[ ! -f .env ]]; then
    cp .env.example .env
    sed -i "s|CHANGEME_POSTGRES_PASSWORD|$(openssl rand -hex 16)|g" .env
    sed -i "s|CHANGEME_AUTHENTIK_SECRET_KEY|$(openssl rand -hex 32)|g" .env
    sed -i "s|CHANGEME_LANGFUSE_SECRET|$(openssl rand -hex 16)|g" .env
    sed -i "s|CHANGEME_LANGFUSE_SALT|$(openssl rand -hex 16)|g" .env
    sed -i "s|CHANGEME_N8N_ENCRYPTION_KEY|$(openssl rand -hex 16)|g" .env
    sed -i "s|CHANGEME_DIFY_SANDBOX_KEY|$(openssl rand -hex 16)|g" .env
    log "Generated .env with random secrets"
  fi

  # Pull third-party images and start stack
  log "Starting ClawFirm stack (--ignore-pull-failures for beta)..."
  docker compose pull --ignore-pull-failures 2>/dev/null || true
  docker compose up -d --remove-orphans 2>/dev/null || true

  # Systemd service for boot persistence
  cat > /etc/systemd/system/clawfirm.service <<'SVC'
[Unit]
Description=ClawFirm Stack
After=docker.service
Requires=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/opt/clawfirm/deploy/solo
ExecStart=/usr/bin/docker compose up -d --remove-orphans
ExecStop=/usr/bin/docker compose down

[Install]
WantedBy=multi-user.target
SVC
  systemctl enable clawfirm.service
  log "ClawFirm installed → /opt/clawfirm — UI at http://localhost:7878"
}

mark_complete() {
  touch "$STATE_DIR/firstboot-complete"
  systemctl disable clawboot-firstboot.service >/dev/null 2>&1 || true
  log "ClawBoot first-boot provisioning complete"
}

main() {
  require_root
  load_config
  prompt_choices

  log "Starting ClawBoot first-boot provisioning"
  apt_install
  install_node
  detect_gpu
  install_openclaw
  install_nanoclaw_scaffold
  install_clawfirm
  configure_telegram_token
  write_status_note
  mark_complete
}

main "$@"
