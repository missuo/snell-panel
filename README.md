# Snell Panel for Surge

## Overview

Snell Panel is a comprehensive management system for Snell proxy nodes that provides unified node management and automatic subscription link generation. The system consists of a backend API server and multiple frontend interfaces (Web UI and iOS App) for seamless node administration.

**Key Components:**
- **Backend Server**: RESTful API server that manages nodes and generates subscription links
- **Web UI**: Browser-based management interface accessible at [snell-panel.owo.nz](http://snell-panel.owo.nz)
- **iOS App**: Native mobile application available through TestFlight for enhanced mobile experience

## Features

### Node Management
- **Multi-node Support**: Unified management of multiple Snell proxy nodes
- **Node Operations**: Add, delete, and modify node configurations
- **Node Renaming**: Customize node names for better organization
- **Relay Nodes**: Support for adding sub-nodes (relay/transit nodes) for advanced routing
- **Real-time Monitoring**: Track node status and performance

### Subscription Management
- **Automatic Generation**: Generate subscription links compatible with Surge and other proxy clients
- **Dynamic Updates**: Subscription links automatically reflect node changes
- **Multiple Formats**: Support for various subscription formats

### API & Integration
- **RESTful API**: Complete API endpoints for programmatic node management
- **Cross-platform**: API serves both Web UI and iOS App with consistent functionality
- **Secure Access**: Token-based authentication for API security

### User Interface
- **Web UI**: Feature-rich browser interface for desktop management
- **iOS App**: Native mobile app with optimized touch interface
- **Responsive Design**: Consistent experience across all devices

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

**Web UI Features:**
- **Dashboard**: Overview of all managed nodes and their status
- **Node Management**: Add, delete, and configure Snell proxy nodes
- **Node Customization**: Rename nodes for better organization and identification
- **Relay Configuration**: Set up sub-nodes (relay/transit nodes) for advanced routing scenarios
- **Subscription Links**: Generate and copy subscription URLs for various proxy clients
- **Real-time Updates**: Live status monitoring and configuration changes

You can get the subscription link from the Web UI.

### Alternative: iOS App

You can also use the iOS App instead of the WebUI for a better mobile experience. Since the app is not yet available on the App Store, you must use TestFlight to install it:

**TestFlight Beta Download:**
[https://testflight.apple.com/join/wKvw64P6](https://testflight.apple.com/join/wKvw64P6)

**iOS App Features:**
- **Native Interface**: Optimized touch interface for iOS devices
- **Full Functionality**: Complete node management capabilities matching the Web UI
- **Node Operations**: Add, delete, rename, and configure nodes on-the-go
- **Relay Management**: Configure sub-nodes and routing policies
- **Subscription Sharing**: Easy copy and share subscription links
- **Offline Access**: View node configurations even when offline

**Installation Notes:** 
- The TestFlight beta may have limited slots available
- You need to install [TestFlight](https://apps.apple.com/app/testflight/id899247664) first on your iOS device
- The iOS app provides the same functionality as the web interface with native mobile optimization

## API Documentation

Snell Panel provides a comprehensive RESTful API for programmatic node management and integration.

**Base URL:** `https://your-panel-domain.com`

**Authentication:** All API requests require a `token` query parameter:
```
GET /entries?token=your_api_token_here
```

### Endpoints

#### 1. Welcome Message
```
GET /
```
Returns a welcome message with basic API information.

**Response:**
```json
{
  "status": "success",
  "message": "Welcome to Snell Panel. Please use the API to manage the entries.\n https://github.com/missuo/snell-panel"
}
```

#### 2. Create Node Entry
```
POST /entry?token=your_token
```

**Request Body:**
```json
{
  "ip": "example.com",
  "port": 443,
  "psk": "your_psk_here",
  "node_name": "Custom Node Name",
  "version": "4"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Entry created successfully",
  "data": {
    "id": 1,
    "ip": "example.com",
    "port": 443,
    "psk": "your_psk_here",
    "country_code": "US",
    "isp": "Example ISP",
    "asn": 12345,
    "node_id": "uuid-generated-string",
    "node_name": "Custom Node Name",
    "version": "4"
  }
}
```

#### 3. List All Nodes
```
GET /entries?token=your_token
```

**Response:**
```json
{
  "status": "success",
  "message": "Entries retrieved successfully",
  "data": [
    {
      "id": 1,
      "ip": "example.com",
      "port": 443,
      "psk": "your_psk_here",
      "country_code": "US",
      "isp": "Example ISP",
      "asn": 12345,
      "node_id": "uuid-string",
      "node_name": "Custom Node Name",
      "version": "4"
    }
  ]
}
```

#### 4. Delete Node by IP
```
DELETE /entry/:ip?token=your_token
```

**Example:** `DELETE /entry/192.168.1.1?token=your_token`

**Response:**
```json
{
  "status": "success",
  "message": "Entry deleted successfully"
}
```

#### 5. Delete Node by Node ID
```
DELETE /entry/node/:node_id?token=your_token
```

**Example:** `DELETE /entry/node/uuid-string?token=your_token`

**Response:**
```json
{
  "status": "success",
  "message": "Entry deleted successfully"
}
```

#### 6. Generate Subscription Link
```
GET /subscribe?token=your_token
```

**Response:** Plain text subscription content compatible with Surge:
```
🇺🇸 Custom Node Name = snell, example.com, 443, psk = your_psk_here, version = 4
🇯🇵 JP Node = snell, jp.example.com, 443, psk = another_psk, version = 4
```

#### 7. Modify Node
```
PUT /modify/:node_id?token=your_token
```

**Request Body:**
```json
{
  "node_name": "New Node Name",
  "ip": "new.example.com"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Node updated successfully"
}
```

### Data Models

#### Entry Model
```json
{
  "id": 1,
  "ip": "string",
  "port": 443,
  "psk": "string",
  "country_code": "string",
  "isp": "string", 
  "asn": 12345,
  "node_id": "string",
  "node_name": "string",
  "version": "string"
}
```

#### API Response Model
```json
{
  "status": "success|error|warning",
  "message": "string",
  "data": "object|array (optional)"
}
```

### Notes
- The `node_id` is automatically generated when creating entries
- IP addresses can be domains or direct IPs - geolocation info is automatically resolved
- Default version is "4" if not specified
- All authenticated endpoints return 401 if token is invalid
- 404 responses are returned for non-existent resources

## TODO

- [x] **Backend API Server**: RESTful API with authentication
- [x] **Web UI**: Full-featured browser interface
- [x] **iOS App**: Native mobile application
- [x] **Node Management**: Add, delete, rename, and configure nodes
- [x] **Relay Nodes**: Support for sub-nodes and transit routing
- [x] **Subscription Generation**: Automatic link generation and updates
- [ ] **Node Health Monitoring**: Real-time health checks and alerts
- [ ] **Advanced Analytics**: Usage statistics and performance metrics
- [ ] **Android App**: Native Android application
- [ ] **Multi-user Support**: User accounts and permission management

## Web UI Source Code

The Snell Panel project consists of multiple components:

- **Backend API Server**: Open source (this repository) - Go-based RESTful API server
- **Web UI Frontend**: Closed source - Browser-based management interface
- **iOS App**: Closed source - Native iOS application

**We are currently not considering open-sourcing the frontend code (Web UI and iOS App) at this time.** The backend API server remains fully open source and provides complete programmatic access to all functionality.

## Contributing

We welcome contributions to the Snell Panel backend server! Here's how you can help:

### Development Setup

**Prerequisites:**
- Go 1.24+ installed
- PostgreSQL database (local or remote)
- Git

**Local Development:**

1. **Clone the repository**
   ```bash
   git clone https://github.com/missuo/snell-panel.git
   cd snell-panel
   ```

2. **Set up environment**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Install dependencies**
   ```bash
   go mod tidy
   ```

4. **Run in development mode**
   ```bash
   go run .
   ```

### Testing

```bash
go mod tidy
go test ./...
```

### Building

**For local platform:**
```bash
go mod tidy
go build .
```

**Cross-platform builds:**
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o snell-panel-linux .

# macOS
GOOS=darwin GOARCH=amd64 go build -o snell-panel-macos .

# Windows  
GOOS=windows GOARCH=amd64 go build -o snell-panel-windows.exe .
```

### Contributing Guidelines

- Fork the repository and create a feature branch
- Write clear commit messages
- Add tests for new functionality
- Ensure all tests pass before submitting
- Update documentation as needed
- Submit a pull request with a clear description

## License

This project is licensed under GPL-3.0.