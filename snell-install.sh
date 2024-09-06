#!/bin/bash

API_URL="http://localhost:59999"
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
    # 创建安装目录
    mkdir -p $INSTALL_DIR
    cd $INSTALL_DIR

    # 下载并解压 Snell
    wget $DOWNLOAD_URL -O snell-server.zip
    unzip snell-server.zip
    rm snell-server.zip

    # 给予执行权限
    chmod +x snell-server

    # 生成随机 PSK 和端口
    PSK=$(openssl rand -base64 16)
    PORT=$(shuf -i 60000-65535 -n 1)

    # 创建预配置的 snell-server.conf 文件
    cat > snell-server.conf <<EOL
[snell-server]
listen = 0.0.0.0:$PORT
psk = $PSK
ipv6 = false
EOL

    # 获取公网 IP
    IP=$(curl -4 ip.sb)

    # 发送数据到 API
    curl -X POST "$API_URL/entry" -H "Content-Type: application/json" -d "{\"ip\":\"$IP\",\"port\":$PORT,\"psk\":\"$PSK\"}"

    # 创建 systemd 服务文件
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

    # 启用并启动服务
    sudo systemctl daemon-reload
    sudo systemctl enable snell
    sudo systemctl start snell

    echo "Snell server installed and started successfully."
    echo "Installation directory: $INSTALL_DIR"
    echo "Server IP: $IP"
    echo "Server Port: $PORT"
    echo "PSK: $PSK"
}

uninstall_snell() {
    # 停止并禁用服务
    sudo systemctl stop snell
    sudo systemctl disable snell

    # 删除文件
    rm -rf $INSTALL_DIR
    sudo rm /etc/systemd/system/snell.service

    # 从 API 删除条目
    IP=$(curl -4 ip.sb)
    curl -X DELETE "$API_URL/entry/$IP"

    sudo systemctl daemon-reload

    echo "Snell server uninstalled successfully."
}

# 主逻辑
if [ "$1" == "install" ]; then
    install_snell
elif [ "$1" == "uninstall" ]; then
    uninstall_snell
else
    echo "Usage: $0 {install|uninstall}"
    exit 1
fi