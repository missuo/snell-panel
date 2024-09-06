# Snell Panel for Surge

## Overview

This project provides unified management of Snell Nodes, automatically generating subscription links for users. The panel's web interface can be accessed via [snell-panel.owo.nz](http://snell-panel.owo.nz). You only need to enter your server’s API URL and token to start using it. 

## How to Use

1. **Start the Snell Panel server**

   Run the following command to start the server:

   ```bash
   ./snell-panel -token=your_token
   ```

   Use Docker Compose:

   ```bash
   docker compose up -d
   ```

2. **Install Snell Server on Your Node**

   Use the following command to install Snell Server on your node:

   ```bash
   bash <(curl -Ls https://raw.githubusercontent.com/missuo/snell-panel/main/snell-install.sh) install your_panel_url your_token

   or

   bash <(curl -Ls https://ssa.sx/sn) install your_panel_url your_token
   ```

   For Example:

   ```bash
   bash <(curl -Ls https://ssa.sx/sn) install http://snell.owo.nz helloworld
   ```

3. **Access the Web UI**

   Access the management Web UI using the following link:

   [http://snell-panel.owo.nz](http://snell-panel.owo.nz)

   ![Snell Panel](./screenshots/web.png)

   You can get the subscription link from the Web UI.

## Features

- Unified management of multiple Snell nodes
- Automatic subscription link generation
- Simple web interface for easy management

## TODO

- [x] Web UI implemented
- [ ] Node survival detection

## License

This project is licensed under GPL-3.0.