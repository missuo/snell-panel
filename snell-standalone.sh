#!/bin/bash

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

SNELL_VERSION="v5.0.1"
INSTALL_DIR="/etc/snell"
CONF_FILE="${INSTALL_DIR}/snell-server.conf"
BIN_FILE="/usr/local/bin/snell-server"
SERVICE_FILE="/etc/systemd/system/snell.service"

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
    if ! command -v unzip &> /dev/null; then
        apt-get update && apt-get install -y unzip wget curl
    fi
}

# Enable TCP Fast Open
enable_tfo() {
    echo -e "${GREEN}Enabling TCP Fast Open (TFO)...${NC}"

    TFO_SETTING="net.ipv4.tcp_fastopen = 3"
    SYSCTL_CONF="/etc/sysctl.conf"

    # Check if TFO is already configured
    if grep -q "^net.ipv4.tcp_fastopen" "$SYSCTL_CONF" 2>/dev/null; then
        # Update existing setting
        sed -i 's/^net.ipv4.tcp_fastopen.*/'"$TFO_SETTING"'/' "$SYSCTL_CONF"
        echo -e "${GREEN}Updated existing TFO setting in $SYSCTL_CONF${NC}"
    else
        # Append new setting
        echo "$TFO_SETTING" >> "$SYSCTL_CONF"
        echo -e "${GREEN}Added TFO setting to $SYSCTL_CONF${NC}"
    fi

    # Apply the setting
    sysctl -p > /dev/null 2>&1
    echo -e "${GREEN}TFO enabled successfully.${NC}"
}

# Architecture detection
get_download_url() {
    ARCH=$(uname -m)
    case "$ARCH" in
        x86_64)
            echo "https://dl.nssurge.com/snell/snell-server-${SNELL_VERSION}-linux-amd64.zip"
            ;;
        aarch64)
            echo "https://dl.nssurge.com/snell/snell-server-${SNELL_VERSION}-linux-aarch64.zip"
            ;;
        *)
            echo -e "${RED}Unsupported architecture: $ARCH${NC}"
            exit 1
            ;;
    esac
}

install_snell() {
    install_dependencies
    enable_tfo

    DOWNLOAD_URL=$(get_download_url)

    # Download and Install
    echo -e "${GREEN}Downloading Snell Server ${SNELL_VERSION}...${NC}"
    mkdir -p ${INSTALL_DIR}
    wget -q --show-progress "$DOWNLOAD_URL" -O snell.zip
    unzip -o snell.zip -d ${INSTALL_DIR}
    rm snell.zip
    mv ${INSTALL_DIR}/snell-server ${BIN_FILE}
    chmod +x ${BIN_FILE}

    # Configuration
    echo -e "${GREEN}Configuring Snell Server...${NC}"

    # Port selection
    read -p "Please enter the port (default: random): " PORT
    if [ -z "$PORT" ]; then
        PORT=$(shuf -i 1024-65535 -n 1)
    fi

    # PSK Generation
    PSK=$(openssl rand -base64 16)

    # Generate Random Node Name
    RANDOM_SUFFIX=$(cat /dev/urandom | tr -dc 'A-Z0-9' | fold -w 4 | head -n 1)
    NODE_NAME="Snell-${RANDOM_SUFFIX}"

    # Write Config
    cat > ${CONF_FILE} <<EOF
[snell-server]
listen = 0.0.0.0:${PORT}
psk = ${PSK}
ipv6 = true
EOF

    # Systemd Service
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
ExecStart=${BIN_FILE} -c ${CONF_FILE}
Restart=on-failure

[Install]
WantedBy=multi-user.target
EOF

    # Start Service
    systemctl daemon-reload
    systemctl enable snell
    systemctl restart snell

    # Output Configuration
    IP=$(curl -s -4 http://checkip.amazonaws.com)
    echo -e "${GREEN}Installation Complete!${NC}"
    echo -e "${YELLOW}------------------------------------------------------${NC}"
    echo -e "${NODE_NAME} = snell, ${IP}, ${PORT}, psk = ${PSK}, version = 5"
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
    DOWNLOAD_URL=$(get_download_url)
    
    systemctl stop snell
    
    wget -q --show-progress "$DOWNLOAD_URL" -O snell.zip
    unzip -o snell.zip -d ${INSTALL_DIR}
    rm snell.zip
    mv ${INSTALL_DIR}/snell-server ${BIN_FILE}
    chmod +x ${BIN_FILE}
    
    systemctl start snell
    echo -e "${GREEN}Snell Server upgraded to ${SNELL_VERSION}.${NC}"
}

start_snell() {
    systemctl start snell
    echo -e "${GREEN}Snell Server started.${NC}"
}

stop_snell() {
    systemctl stop snell
    echo -e "${GREEN}Snell Server stopped.${NC}"
}

restart_snell() {
    systemctl restart snell
    echo -e "${GREEN}Snell Server restarted.${NC}"
}

show_menu() {
    echo -e "${BLUE}=== Snell Server Management ===${NC}"
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
