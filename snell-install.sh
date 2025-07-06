#!/bin/bash

SNELL_VERSION="v5.0.0b2"
INSTALL_DIR="$HOME/snell-server"
ARCH=$(uname -m)
USER=$(whoami)

case "$ARCH" in
    x86_64)
        DOWNLOAD_URL="https://dl.nssurge.com/snell/snell-server-${SNELL_VERSION}-linux-amd64.zip"
        ;;
    aarch64)
        DOWNLOAD_URL="https://dl.nssurge.com/snell/snell-server-${SNELL_VERSION}-linux-aarch64.zip"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

install_depend() {
    echo "Checking and installing dependencies..."

    dependencies=("unzip" "wget" "curl")

    for dep in "${dependencies[@]}"; do
        if ! command -v "$dep" &> /dev/null; then
            echo "$dep is not installed. Installing..."

            if command -v apt &> /dev/null; then
                sudo apt update && sudo apt install -y "$dep"
            elif command -v yum &> /dev/null; then
                sudo yum install -y "$dep"
            elif command -v dnf &> /dev/null; then
                sudo dnf install -y "$dep"
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

install_snell() {
    install_depend
    echo "Starting Snell server installation..."

    mkdir -p "$INSTALL_DIR"
    cd "$INSTALL_DIR" || exit

    echo "Downloading Snell server..."
    wget -q "$DOWNLOAD_URL" -O snell-server.zip
    unzip -o snell-server.zip
    rm snell-server.zip
    echo "Download complete."
    chmod +x snell-server

    # Generate random PSK and port
    PSK=$(openssl rand -base64 16)
    PORT=$(shuf -i 60000-65535 -n 1)

    CONFIG_FILE="$INSTALL_DIR/snell-server.conf"

    if [ -f "$CONFIG_FILE" ]; then
        echo "Configuration file already exists. Skipping creation."
    else
        echo "Generating configuration file..."
        # Create pre-configured snell-server.conf file
        cat > "$CONFIG_FILE" <<EOL
[snell-server]
listen = 0.0.0.0:$PORT
psk = $PSK
ipv6 = false
EOL
        echo "Configuration file created."
    fi

    IP=$(curl -s -4 ip.sb)

    # Extract major version number from SNELL_VERSION
    VERSION_WITHOUT_V=${SNELL_VERSION#v}
    MAJOR_VERSION=${VERSION_WITHOUT_V%%.*}

    echo "Sending data to API..."
    # Build API request data
    API_DATA="{\"ip\":\"$IP\",\"port\":$PORT,\"psk\":\"$PSK\",\"version\":\"$MAJOR_VERSION\""
    # If node_name is provided, add it to the API data
    if [ ! -z "$NODE_NAME" ]; then
        API_DATA="$API_DATA,\"node_name\":\"$NODE_NAME\""
    fi
    API_DATA="$API_DATA}"
    
    curl -s -X POST "$API_URL/entry?token=$TOKEN" -H "Content-Type: application/json" \
        -d "$API_DATA"
    echo "API update complete."

    echo "Creating systemd service file..."
    sudo tee /etc/systemd/system/snell.service > /dev/null <<EOL
[Unit]
Description=Snell Proxy Service
After=network.target

[Service]
Type=simple
User=$USER
WorkingDirectory=$INSTALL_DIR
ExecStart=$INSTALL_DIR/snell-server
Restart=on-failure
LimitNOFILE=1048576

[Install]
WantedBy=multi-user.target
EOL
    echo "Systemd service file created."

    echo "Enabling and starting Snell service..."
    # Enable and start service
    sudo systemctl daemon-reload
    sudo systemctl enable snell
    sudo systemctl start snell
    echo "Snell service enabled and started."

    echo "Snell server installation completed successfully."

    echo "Installation summary:"
    echo "---------------------"
    echo "Installation directory: $INSTALL_DIR"
    echo "Server IP: $IP"
    echo "Server Port: $PORT"
    echo "PSK: $PSK"
}

uninstall_snell() {
    echo "Starting Snell server uninstallation..."

    echo "Stopping and disabling Snell service..."
    sudo systemctl stop snell
    sudo systemctl disable snell
    echo "Snell service stopped and disabled."

    echo "Removing Snell files..."
    # Remove files
    rm -rf "$INSTALL_DIR"
    sudo rm /etc/systemd/system/snell.service
    echo "Snell files removed."

    echo "Deleting entry from API..."
    # Delete entry from API with token
    IP=$(curl -s -4 ip.sb)
    curl -s -X DELETE "$API_URL/entry/$IP?token=$TOKEN"
    echo "API entry deleted."

    sudo systemctl daemon-reload

    echo "Snell server uninstallation completed successfully."
}

update_snell() {
    echo "Starting Snell server update..."
    sudo systemctl stop snell

    cd "$INSTALL_DIR" || exit

    echo "Downloading Snell server..."
    wget -q "$DOWNLOAD_URL" -O snell-server.zip
    unzip -o snell-server.zip
    rm snell-server.zip
    echo "Download complete."

    chmod +x snell-server

    sudo systemctl start snell
    echo "Snell server updated and restarted."
}

# Main logic
ACTION=$1

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
        NODE_NAME=$4  # NODE_NAME is now optional
        ;;
    update)
        # API_URL, TOKEN, NODE_NAME are optional for update
        API_URL=$2
        TOKEN=$3
        NODE_NAME=$4
        ;;
    *)
        echo "Invalid action. Usage: $0 {install|uninstall|update} [API_URL TOKEN [NODE_NAME]]"
        exit 1
        ;;
esac

case "$ACTION" in
    install)
        install_snell
        ;;
    uninstall)
        uninstall_snell
        ;;
    update)
        update_snell
        ;;
esac
