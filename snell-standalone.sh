#!/bin/bash
#
# Standalone Snell v6 server manager (no panel API). Default binary: the official
# Surge snell-server v6.0.0b2; OpenSnell v6 (private repo) is available via gh.

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

SURGE_VERSION="${SURGE_VERSION:-v6.0.0b2}"
OPENSNELL_REPO="${OPENSNELL_REPO:-missuo/opensnell-v6}"
OPENSNELL_TAG="${OPENSNELL_TAG:-latest}"
DNS_IP_PREFERENCE="${DNS_IP_PREFERENCE:-default}"

INSTALL_DIR="/etc/snell"
CONF_FILE="${INSTALL_DIR}/snell-server.conf"
BIN_FILE="/usr/local/bin/snell-server"
SERVICE_FILE="/etc/systemd/system/snell.service"
META_FILE="${INSTALL_DIR}/.variant"

# Check if running as root
check_root() {
    if [ "$EUID" -ne 0 ]; then
        echo -e "${RED}Please run as root${NC}"
        exit 1
    fi
}

# Install dependencies
install_dependencies() {
    echo -e "${GREEN}Checking dependencies...${NC}"
    if ! command -v unzip &> /dev/null || ! command -v wget &> /dev/null; then
        apt-get update && apt-get install -y unzip wget curl
    fi
}

# Enable TCP Fast Open
enable_tfo() {
    echo -e "${GREEN}Enabling TCP Fast Open (TFO)...${NC}"
    TFO_SETTING="net.ipv4.tcp_fastopen = 3"
    SYSCTL_CONF="/etc/sysctl.conf"
    if grep -q "^net.ipv4.tcp_fastopen" "$SYSCTL_CONF" 2>/dev/null; then
        sed -i 's/^net.ipv4.tcp_fastopen.*/'"$TFO_SETTING"'/' "$SYSCTL_CONF"
    else
        echo "$TFO_SETTING" >> "$SYSCTL_CONF"
    fi
    sysctl -p > /dev/null 2>&1
    echo -e "${GREEN}TFO enabled successfully.${NC}"
}

ensure_gh() {
    if ! command -v gh &> /dev/null; then
        echo -e "${RED}The 'opensnell' variant needs the GitHub CLI (gh).${NC}"
        echo -e "${YELLOW}Install it from https://cli.github.com/ and run 'gh auth login'.${NC}"
        exit 1
    fi
    if ! gh auth status &> /dev/null; then
        echo -e "${RED}gh is not authenticated. Run 'gh auth login' (needs read access to ${OPENSNELL_REPO}).${NC}"
        exit 1
    fi
}

# Download the snell-server binary to $BIN_FILE for the given variant ($1).
download_binary() {
    local variant="$1" ARCH
    ARCH=$(uname -m)
    mkdir -p "${INSTALL_DIR}"

    if [ "$variant" = "opensnell" ]; then
        ensure_gh
        local oa
        case "$ARCH" in
            x86_64)        oa="amd64" ;;
            aarch64|arm64) oa="arm64" ;;
            i386|i686)     oa="386"   ;;
            armv7l|armv7)  oa="armv7" ;;
            *) echo -e "${RED}Unsupported architecture for OpenSnell: $ARCH${NC}"; exit 1 ;;
        esac
        echo -e "${GREEN}Downloading OpenSnell v6 (${OPENSNELL_REPO}, tag=${OPENSNELL_TAG}, linux-${oa}) via gh...${NC}"
        local tagarg=()
        [ "$OPENSNELL_TAG" != "latest" ] && tagarg=("$OPENSNELL_TAG")
        gh release download "${tagarg[@]}" -R "$OPENSNELL_REPO" \
            -p "snell-server-linux-${oa}" -O "${BIN_FILE}" --clobber || {
                echo -e "${RED}gh release download failed.${NC}"; exit 1; }
        chmod +x "${BIN_FILE}"
    else
        local sa
        case "$ARCH" in
            x86_64)        sa="amd64"   ;;
            aarch64|arm64) sa="aarch64" ;;
            i386|i686)     sa="i386"    ;;
            *) echo -e "${RED}Surge snell-server ${SURGE_VERSION} is not available for $ARCH${NC}"; exit 1 ;;
        esac
        local url="https://dl.nssurge.com/snell/snell-server-${SURGE_VERSION}-linux-${sa}.zip"
        echo -e "${GREEN}Downloading official Surge snell-server ${SURGE_VERSION} (linux-${sa})...${NC}"
        wget -q --show-progress "$url" -O /tmp/snell.zip || { echo -e "${RED}Download failed: $url${NC}"; exit 1; }
        unzip -o /tmp/snell.zip -d "${INSTALL_DIR}"
        rm -f /tmp/snell.zip
        mv -f "${INSTALL_DIR}/snell-server" "${BIN_FILE}"
        chmod +x "${BIN_FILE}"
    fi
    echo "$variant" > "${META_FILE}"
}

exec_start() {
    if [ "$1" = "opensnell" ]; then
        echo "${BIN_FILE} --v6b2 -c ${CONF_FILE}"
    else
        echo "${BIN_FILE} -c ${CONF_FILE}"
    fi
}

install_snell() {
    install_dependencies
    enable_tfo

    # Variant selection
    echo -e "${BLUE}Choose a binary:${NC}"
    echo "  1) Official Surge snell-server ${SURGE_VERSION}  (default)"
    echo "  2) OpenSnell v6 (private repo, via gh CLI)"
    read -p "Variant [1]: " variant_choice
    local variant="official"
    [ "$variant_choice" = "2" ] && variant="opensnell"

    download_binary "$variant"

    echo -e "${GREEN}Configuring Snell v6 Server...${NC}"
    read -p "Please enter the port (default: random): " PORT
    if [ -z "$PORT" ]; then
        PORT=$(shuf -i 1024-65535 -n 1)
    fi

    # 24-char alphanumeric PSK (>= 16 bytes, required by v6).
    PSK=$(openssl rand -base64 24 | tr -dc 'A-Za-z0-9' | head -c 24)

    RANDOM_SUFFIX=$(cat /dev/urandom | tr -dc 'A-Z0-9' | fold -w 4 | head -n 1)
    NODE_NAME="Snell-${RANDOM_SUFFIX}"

    # v6 config: no `obfs`; `ipv6` replaced by `dns-ip-preference`.
    cat > ${CONF_FILE} <<EOF
[snell-server]
listen = 0.0.0.0:${PORT}
psk = ${PSK}
dns-ip-preference = ${DNS_IP_PREFERENCE}
EOF
    chmod 600 ${CONF_FILE}

    echo -e "${GREEN}Creating Systemd Service...${NC}"
    cat > ${SERVICE_FILE} <<EOF
[Unit]
Description=Snell Proxy Service
After=network.target

[Service]
Type=simple
User=root
Group=root
LimitNOFILE=32768
ExecStart=$(exec_start "$variant")
Restart=on-failure

[Install]
WantedBy=multi-user.target
EOF

    systemctl daemon-reload
    systemctl enable snell
    systemctl restart snell

    IP=$(curl -s -4 http://checkip.amazonaws.com)
    echo -e "${GREEN}Installation Complete!${NC}"
    echo -e "${YELLOW}------------------------------------------------------${NC}"
    echo -e "${NODE_NAME} = snell, ${IP}, ${PORT}, psk = ${PSK}, version = 6"
    echo -e "${YELLOW}------------------------------------------------------${NC}"
}

uninstall_snell() {
    echo -e "${YELLOW}Uninstalling Snell Server...${NC}"
    systemctl stop snell
    systemctl disable snell
    rm -f ${SERVICE_FILE}
    systemctl daemon-reload
    rm -f ${BIN_FILE}
    rm -rf ${INSTALL_DIR}
    echo -e "${GREEN}Snell Server uninstalled successfully.${NC}"
}

upgrade_snell() {
    echo -e "${GREEN}Upgrading Snell Server...${NC}"
    local variant="official"
    [ -f "${META_FILE}" ] && variant=$(cat "${META_FILE}")
    systemctl stop snell
    download_binary "$variant"
    systemctl start snell
    echo -e "${GREEN}Snell Server upgraded (variant: ${variant}).${NC}"
}

start_snell()   { systemctl start snell;   echo -e "${GREEN}Snell Server started.${NC}"; }
stop_snell()    { systemctl stop snell;    echo -e "${GREEN}Snell Server stopped.${NC}"; }
restart_snell() { systemctl restart snell; echo -e "${GREEN}Snell Server restarted.${NC}"; }

show_menu() {
    echo -e "${BLUE}=== Snell v6 Server Management ===${NC}"
    echo "1. Install Snell"
    echo "2. Uninstall Snell"
    echo "3. Upgrade Snell"
    echo "4. Start Snell"
    echo "5. Stop Snell"
    echo "6. Restart Snell"
    echo "0. Exit"
    echo -e "${BLUE}==============================${NC}"
    read -p "Please enter your choice [0-6]: " choice
    case $choice in
        1) install_snell ;;
        2) uninstall_snell ;;
        3) upgrade_snell ;;
        4) start_snell ;;
        5) stop_snell ;;
        6) restart_snell ;;
        0) exit 0 ;;
        *) echo -e "${RED}Invalid choice, please try again.${NC}" ;;
    esac
}

# Main execution
check_root

if [ -n "$1" ]; then
    case "$1" in
        install) install_snell ;;
        uninstall) uninstall_snell ;;
        upgrade) upgrade_snell ;;
        start) start_snell ;;
        stop) stop_snell ;;
        restart) restart_snell ;;
        *) echo "Usage: $0 {install|uninstall|upgrade|start|stop|restart}" ;;
    esac
else
    while true; do
        show_menu
        echo ""
        read -p "Press Enter to continue..."
    done
fi
