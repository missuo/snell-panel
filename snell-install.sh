#!/bin/bash

SNELL_VERSION="v4.1.0"
INSTALL_DIR="$HOME/snell-server"
ARCH=$(uname -m)

if [ "$ARCH" == "x86_64" ]; then
    DOWNLOAD_URL="https://dl.nssurge.com/snell/snell-server-v4.1.0-linux-amd64.zip"
elif [ "$ARCH" == "aarch64" ]; then
    DOWNLOAD_URL="https://dl.nssurge.com/snell/snell-server-v4.1.0-linux-aarch64.zip"
else
    echo "Unsupported architecture: $ARCH"
    exit 1
fi

install_snell() {
    echo "Starting Snell server installation..."
    echo

    # Create installation directory
    mkdir -p $INSTALL_DIR
    cd $INSTALL_DIR

    echo "Downloading Snell server..."
    wget $DOWNLOAD_URL -O snell-server.zip
    unzip snell-server.zip
    rm snell-server.zip
    echo "Download complete."
    echo

    # Grant execution permission
    chmod +x snell-server

    # Generate random PSK and port
    PSK=$(openssl rand -base64 16)
    PORT=$(shuf -i 60000-65535 -n 1)

    echo "Generating configuration file..."
    # Create pre-configured snell-server.conf file
    cat > snell-server.conf <<EOL
[snell-server]
listen = 0.0.0.0:$PORT
psk = $PSK
ipv6 = false
EOL
    echo "Configuration file created."
    echo

    # Get public IP
    IP=$(curl -4 ip.sb)

    echo "Sending data to API..."
    # Send data to API with token
    curl -X POST "$API_URL/entry?token=$TOKEN" -H "Content-Type: application/json" -d "{\"ip\":\"$IP\",\"port\":$PORT,\"psk\":\"$PSK\"}"
    echo
    echo "API update complete."
    echo

    echo "Creating systemd service file..."
    # Create systemd service file
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
    echo

    echo "Enabling and starting Snell service..."
    # Enable and start service
    sudo systemctl daemon-reload
    sudo systemctl enable snell
    sudo systemctl start snell
    echo "Snell service enabled and started."
    echo

    echo "Snell server installation completed successfully."
    echo
    echo "Installation summary:"
    echo "---------------------"
    echo "Installation directory: $INSTALL_DIR"
    echo "Server IP: $IP"
    echo "Server Port: $PORT"
    echo "PSK: $PSK"
    echo
}

uninstall_snell() {
    echo "Starting Snell server uninstallation..."
    echo

    echo "Stopping and disabling Snell service..."
    # Stop and disable service
    sudo systemctl stop snell
    sudo systemctl disable snell
    echo "Snell service stopped and disabled."
    echo

    echo "Removing Snell files..."
    # Remove files
    rm -rf $INSTALL_DIR
    sudo rm /etc/systemd/system/snell.service
    echo "Snell files removed."
    echo

    echo "Deleting entry from API..."
    # Delete entry from API with token
    IP=$(curl -4 ip.sb)
    curl -X DELETE "$API_URL/entry/$IP?token=$TOKEN"
    echo
    echo "API entry deleted."
    echo

    sudo systemctl daemon-reload

    echo "Snell server uninstallation completed successfully."
    echo
}

# Main logic
if [ $# -lt 3 ]; then
    echo "Usage: $0 {install|uninstall} API_URL TOKEN"
    echo
    echo "Parameters:"
    echo "  {install|uninstall}  - Action to perform"
    echo "  API_URL              - URL of the API"
    echo "  TOKEN                - Authentication token for API access"
    echo
    exit 1
fi

ACTION=$1
API_URL=$2
TOKEN=$3

if [ "$ACTION" == "install" ]; then
    install_snell
elif [ "$ACTION" == "uninstall" ]; then
    uninstall_snell
else
    echo "Invalid action. Usage: $0 {install|uninstall} API_URL TOKEN"
    echo
    exit 1
fi