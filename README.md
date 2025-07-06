# Snell Panel for Surge

## Overview

This project provides unified management of Snell Nodes, automatically generating subscription links for users. The panel's web interface can be accessed via [snell-panel.owo.nz](http://snell-panel.owo.nz). You only need to enter your server's API URL and token to start using it. 

## Features

- Unified management of multiple Snell nodes
- Automatic subscription link generation
- Simple web interface for easy management

## How to Use

### Method 1: Direct Binary Execution

1. **Configure environment variables**

   Create a `.env` file or set environment variables:

   ```bash
   export API_TOKEN=your_token_here
   export DATABASE_URL=your_database_url_here  # e.g., postgres://user:pass@host:port/dbname
   export PORT=8080  # Optional, defaults to 8080
   ```

2. **Start the Snell Panel server**

   ```bash
   ./snell-panel
   ```

### Method 2: Using Docker Compose (Recommended)

1. **Configure the compose.yaml file**

   Edit the `compose.yaml` file and update the environment variables:

   ```yaml
   environment:
     - API_TOKEN=your_token_here
     - DATABASE_URL=your_database_url_here
   ```

2. **Start the service**

   ```bash
   docker-compose up -d
   ```

   The service will be available on port 9997.

### Method 3: Using Docker

1. **Build the Docker image**

   ```bash
   docker build -t snell-panel .
   ```

2. **Run the container**

   ```bash
   docker run -d \
     --name snell-panel \
     -p 8080:8080 \
     -e API_TOKEN=your_token_here \
     -e DATABASE_URL=your_database_url_here \
     snell-panel
   ```

### Method 4: Vercel Serverless Deployment

1. **Set up Supabase Database**

   - Go to [supabase.com](https://supabase.com) and create a new project
   - Navigate to Settings > Database and copy the connection string
   - The connection string format should be: `postgresql://postgres:[YOUR-PASSWORD]@[YOUR-HOST]:[YOUR-PORT]/postgres`

   **Detailed Steps:**
   - Create a Supabase account and new project
   - Go to Project Settings → Database
   - Under "Connection string", select "URI" 
   - Copy the connection string (it will look like `postgresql://postgres:[YOUR-PASSWORD]@db.xxxxx.supabase.co:5432/postgres`)
   - Replace `[YOUR-PASSWORD]` with your actual database password

2. **Deploy to Vercel**

   Click the button below to deploy directly to Vercel:

   [![Deploy with Vercel](https://vercel.com/button)](https://vercel.com/new/clone?repository-url=https://github.com/missuo/snell-panel)

   Or manually deploy:

   ```bash
   # Clone the repository
   git clone https://github.com/missuo/snell-panel.git
   cd snell-panel

   # Install Vercel CLI
   npm i -g vercel

   # Deploy to Vercel
   vercel --prod
   ```

3. **Configure Environment Variables in Vercel**

   In your Vercel dashboard, go to your project settings and add the following environment variables:

   ```
   API_TOKEN=your_token_here
   DATABASE_URL=your_supabase_database_url
   ```

   Example Supabase DATABASE_URL:
   ```
   postgresql://postgres:your_password@db.abcdefghijklmnop.supabase.co:5432/postgres
   ```

4. **Access your deployment**

   Your Snell Panel will be available at `https://your-project-name.vercel.app`

   **Advantages of Vercel Deployment:**
   - Serverless architecture with automatic scaling
   - Global edge network for faster response times
   - Automatic HTTPS and custom domain support
   - Zero server maintenance required
   - Free tier available for personal projects

## Install Snell Server

Use the following command to **install** Snell Server:

```bash
bash <(curl -Ls https://ssa.sx/sn) install your_panel_url your_token custom_node_name
```

Use the following command to **uninstall** Snell Server:

```bash
bash <(curl -Ls https://ssa.sx/sn) uninstall your_panel_url your_token custom_node_name
```

`custom_node_name` is optional. If your node name contains spaces, please use quotes. For example:

```bash
bash <(curl -Ls https://ssa.sx/sn) install your_panel_url your_token "My Node Name"
```

Use the following command to **update** Snell Server:
```bash
bash <(curl -Ls https://ssa.sx/sn) update
```

## Access the Web UI

Access the management Web UI using the following link:

[https://snell-panel.owo.nz](https://snell-panel.owo.nz)

![Snell Panel](./screenshots/web.png)

You can get the subscription link from the Web UI.

## TODO

- [x] Web UI implemented
- [ ] Node survival detection

## Web UI Source Code

**We are not considering open-sourcing any code for the Web UI at this time.**

## Contributing

### Test

```bash
go mod tidy
go run .
```

### Build

```bash
go mod tidy
go build .
```

## License

This project is licensed under GPL-3.0.