#!/bin/bash
#
# Snell node installer for snell-panel — Snell v6 (default: official Surge v6.0.0b2).
#
# Binary source (SNELL_VARIANT):
#   official   (default)  Surge snell-server v6.0.0b2 — public, dl.nssurge.com
#   opensnell             OpenSnell v6 — the repo is private, so the binary is
#                         pulled from missuo/opensnell-v6 GitHub Releases with the
#                         GitHub CLI (`gh`). Requires `gh` installed AND
#                         authenticated on this host (`gh auth login`).
#
# Snell v6 config differs from v5: `obfs` is gone and `ipv6 = true/false` is
# replaced by `dns-ip-preference`. The PSK must be 16–255 bytes.
#
# Usage:
#   bash snell-install.sh install   API_URL TOKEN [NODE_NAME]
#   bash snell-install.sh uninstall API_URL TOKEN
#   bash snell-install.sh update   [API_URL TOKEN [NODE_NAME]]
#
# Knobs (environment variables):
#   SNELL_VARIANT       official | opensnell           (default: official)
#   SURGE_VERSION       official download version       (default: v6.0.0b2)
#   OPENSNELL_REPO      private repo for the gh variant (default: missuo/opensnell-v6)
#   OPENSNELL_TAG       gh release tag, or "latest"     (default: latest)
#   DNS_IP_PREFERENCE   default|prefer-ipv4|prefer-ipv6|ipv4-only|ipv6-only (default: default)

# Record whether the caller set SNELL_VARIANT explicitly (so `update` can honour
# it over the value saved at install time).
[ -n "${SNELL_VARIANT:-}" ] && SNELL_VARIANT_SET=1
SNELL_VARIANT="${SNELL_VARIANT:-official}"
SURGE_VERSION="${SURGE_VERSION:-v6.0.0b2}"
OPENSNELL_REPO="${OPENSNELL_REPO:-missuo/opensnell-v6}"
OPENSNELL_TAG="${OPENSNELL_TAG:-latest}"
DNS_IP_PREFERENCE="${DNS_IP_PREFERENCE:-default}"
SNELL_PROTOCOL_VERSION="6"   # POSTed to the panel + emitted as `version = 6`

INSTALL_DIR="${INSTALL_DIR:-$HOME/snell-server}"
CONFIG_FILE="$INSTALL_DIR/snell-server.conf"
BIN="$INSTALL_DIR/snell-server"
META_FILE="$INSTALL_DIR/.variant"
ARCH=$(uname -m)
RUN_USER=$(whoami)

# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------
get_version_info() {
    if [ -f "$BIN" ]; then
        local out
        out=$("$BIN" -v 2>&1 | grep "snell-server v")
        if [ -n "$out" ]; then
            echo "$out" | sed 's/.*\(snell-server v[^)]*)\).*/\1/'
        else
            echo "snell-server (version unavailable)"
        fi
    else
        echo "Snell server binary not found"
    fi
}

# A 24-char alphanumeric PSK (>= 16 bytes, no /+= so it is safe in configs/URLs).
gen_psk() {
    openssl rand -base64 24 | tr -dc 'A-Za-z0-9' | head -c 24
}

install_depend() {
    echo "Checking and installing dependencies..."
    local dependencies=("unzip" "wget" "curl" "openssl")
    for dep in "${dependencies[@]}"; do
        if ! command -v "$dep" &> /dev/null; then
            echo "$dep is not installed. Installing..."
            if command -v apt &> /dev/null; then
                sudo apt update && sudo apt install -y "$dep"
            elif command -v dnf &> /dev/null; then
                sudo dnf install -y "$dep"
            elif command -v yum &> /dev/null; then
                sudo yum install -y "$dep"
            elif command -v pacman &> /dev/null; then
                sudo pacman -Sy --noconfirm "$dep"
            elif command -v zypper &> /dev/null; then
                sudo zypper install -y "$dep"
            else
                echo "Unsupported package manager. Please install $dep manually."
                exit 1
            fi
        else
            echo "$dep is already installed."
        fi
    done
}

ensure_gh() {
    if ! command -v gh &> /dev/null; then
        echo "ERROR: the 'opensnell' variant needs the GitHub CLI (gh), which is not installed."
        echo "Install it from https://cli.github.com/ and run 'gh auth login', then retry."
        exit 1
    fi
    if ! gh auth status &> /dev/null; then
        echo "ERROR: gh is installed but not authenticated. Run 'gh auth login' (needs read"
        echo "access to ${OPENSNELL_REPO}), then retry."
        exit 1
    fi
}

enable_tfo() {
    echo "Enabling TCP Fast Open (TFO)..."
    local TFO_SETTING="net.ipv4.tcp_fastopen = 3"
    local SYSCTL_CONF="/etc/sysctl.conf"
    if grep -q "^net.ipv4.tcp_fastopen" "$SYSCTL_CONF" 2>/dev/null; then
        sudo sed -i 's/^net.ipv4.tcp_fastopen.*/'"$TFO_SETTING"'/' "$SYSCTL_CONF"
    else
        echo "$TFO_SETTING" | sudo tee -a "$SYSCTL_CONF" > /dev/null
    fi
    sudo sysctl -p > /dev/null 2>&1
    echo "TFO enabled successfully."
}

# Download the snell-server binary into $INSTALL_DIR per the selected variant.
download_binary() {
    mkdir -p "$INSTALL_DIR"
    cd "$INSTALL_DIR" || exit 1

    if [ "$SNELL_VARIANT" = "opensnell" ]; then
        ensure_gh
        local oa
        case "$ARCH" in
            x86_64)        oa="amd64" ;;
            aarch64|arm64) oa="arm64" ;;
            i386|i686)     oa="386"   ;;
            armv7l|armv7)  oa="armv7" ;;
            *) echo "Unsupported architecture for OpenSnell: $ARCH"; exit 1 ;;
        esac
        echo "Downloading OpenSnell v6 server (${OPENSNELL_REPO}, tag=${OPENSNELL_TAG}, linux-${oa}) via gh..."
        local tagarg=()
        [ "$OPENSNELL_TAG" != "latest" ] && tagarg=("$OPENSNELL_TAG")
        if ! gh release download "${tagarg[@]}" -R "$OPENSNELL_REPO" \
                -p "snell-server-linux-${oa}" -O "$BIN" --clobber; then
            echo "ERROR: gh release download failed. Check the tag/asset and your gh auth."
            exit 1
        fi
        chmod +x "$BIN"
    else
        local sa
        case "$ARCH" in
            x86_64)        sa="amd64"   ;;
            aarch64|arm64) sa="aarch64" ;;
            i386|i686)     sa="i386"    ;;
            *) echo "Surge snell-server ${SURGE_VERSION} is not available for $ARCH."; exit 1 ;;
        esac
        local url="https://dl.nssurge.com/snell/snell-server-${SURGE_VERSION}-linux-${sa}.zip"
        echo "Downloading official Surge snell-server ${SURGE_VERSION} (linux-${sa})..."
        wget -q "$url" -O snell-server.zip || { echo "Download failed: $url"; exit 1; }
        unzip -o snell-server.zip
        rm -f snell-server.zip
        chmod +x "$BIN"
    fi
    echo "$SNELL_VARIANT" > "$META_FILE"
}

# ExecStart line for the systemd unit. OpenSnell needs --v6b2 to speak v6; the
# official Surge server auto-negotiates the version the client asks for.
exec_start() {
    if [ "$SNELL_VARIANT" = "opensnell" ]; then
        echo "$BIN --v6b2 -c $CONFIG_FILE"
    else
        echo "$BIN -c $CONFIG_FILE"
    fi
}

# ---------------------------------------------------------------------------
# Actions
# ---------------------------------------------------------------------------
install_snell() {
    install_depend
    enable_tfo
    echo "Starting Snell v6 server installation (variant: ${SNELL_VARIANT})..."

    download_binary
    echo "Download complete."

    PSK=$(gen_psk)
    PORT=$(shuf -i 60000-65535 -n 1)

    if [ -f "$CONFIG_FILE" ]; then
        echo "Configuration file already exists. Reusing its PSK/port."
        PSK=$(grep -E '^psk' "$CONFIG_FILE" | head -1 | cut -d= -f2 | tr -d ' ')
        PORT=$(grep -E '^listen' "$CONFIG_FILE" | head -1 | sed -E 's/.*:([0-9]+).*/\1/')
    else
        echo "Generating Snell v6 configuration file..."
        # v6 config: no `obfs`; `ipv6` replaced by `dns-ip-preference`.
        cat > "$CONFIG_FILE" <<EOL
[snell-server]
listen = 0.0.0.0:$PORT
psk = $PSK
dns-ip-preference = $DNS_IP_PREFERENCE
EOL
        chmod 600 "$CONFIG_FILE"
        echo "Configuration file created."
    fi

    IP=$(curl -s -4 ip.sb)

    echo "Sending data to API..."
    API_DATA="{\"ip\":\"$IP\",\"port\":$PORT,\"psk\":\"$PSK\",\"version\":\"$SNELL_PROTOCOL_VERSION\""
    if [ -n "${NODE_NAME:-}" ]; then
        API_DATA="$API_DATA,\"node_name\":\"$NODE_NAME\""
    fi
    API_DATA="$API_DATA}"
    curl -s -X POST "$API_URL/entry?token=$TOKEN" -H "Content-Type: application/json" -d "$API_DATA"
    echo "API update complete."

    echo "Creating systemd service file..."
    sudo tee /etc/systemd/system/snell.service > /dev/null <<EOL
[Unit]
Description=Snell Proxy Service
After=network.target

[Service]
Type=simple
User=$RUN_USER
WorkingDirectory=$INSTALL_DIR
ExecStart=$(exec_start)
Restart=on-failure
LimitNOFILE=1048576

[Install]
WantedBy=multi-user.target
EOL

    echo "Enabling and starting Snell service..."
    sudo systemctl daemon-reload
    sudo systemctl enable snell
    sudo systemctl restart snell

    echo "Snell v6 server installation completed successfully."
    echo "Installation summary:"
    echo "---------------------"
    echo "Variant:                ${SNELL_VARIANT}"
    echo "Installation directory: $INSTALL_DIR"
    echo "Server IP:              $IP"
    echo "Server Port:            $PORT"
    echo "PSK:                    $PSK"
    echo "Snell protocol version: $SNELL_PROTOCOL_VERSION"
    echo "Installed binary:       $(get_version_info)"
}

uninstall_snell() {
    echo "Starting Snell server uninstallation..."
    sudo systemctl stop snell
    sudo systemctl disable snell
    rm -rf "$INSTALL_DIR"
    sudo rm -f /etc/systemd/system/snell.service

    echo "Deleting entry from API..."
    IP=$(curl -s -4 ip.sb)
    curl -s -X DELETE "$API_URL/entry/$IP?token=$TOKEN"
    echo "API entry deleted."

    sudo systemctl daemon-reload
    echo "Snell server uninstallation completed successfully."
}

update_snell() {
    echo "Starting Snell server update..."
    # Preserve the variant chosen at install time unless SNELL_VARIANT is set.
    if [ -z "${SNELL_VARIANT_SET:-}" ] && [ -f "$META_FILE" ]; then
        SNELL_VARIANT=$(cat "$META_FILE")
    fi
    local OLD_VERSION; OLD_VERSION=$(get_version_info)

    sudo systemctl stop snell
    download_binary
    sudo systemctl start snell

    local NEW_VERSION; NEW_VERSION=$(get_version_info)
    echo "Update completed (variant: ${SNELL_VARIANT}):"
    echo "From: $OLD_VERSION"
    echo "To:   $NEW_VERSION"
}

# ---------------------------------------------------------------------------
# Main
# ---------------------------------------------------------------------------
ACTION="${1:-}"
if [ -z "$ACTION" ]; then
    echo "Usage: $0 {install|uninstall|update} [API_URL TOKEN [NODE_NAME]]"
    exit 1
fi

case "$ACTION" in
    install|uninstall)
        if [ $# -lt 3 ]; then
            echo "Usage: $0 {install|uninstall} API_URL TOKEN [NODE_NAME]"
            exit 1
        fi
        API_URL=$2
        TOKEN=$3
        NODE_NAME=${4:-}
        ;;
    update)
        API_URL=${2:-}
        TOKEN=${3:-}
        NODE_NAME=${4:-}
        ;;
    *)
        echo "Invalid action. Usage: $0 {install|uninstall|update} [API_URL TOKEN [NODE_NAME]]"
        exit 1
        ;;
esac

case "$ACTION" in
    install)   install_snell ;;
    uninstall) uninstall_snell ;;
    update)    update_snell ;;
esac
