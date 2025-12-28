# Quotio Electron Go

A cross-platform desktop application for managing multiple AI provider subscriptions (Claude, ChatGPT, Gemini, Antigravity, Copilot, etc.) with any coding AI agent (Claude Code, OpenCode, Droid, etc.). The app runs locally, tracks quota usage, and automatically switches accounts when rate limits are hit.

## Features

- **Multi-Provider Support**: Connect accounts from Claude, OpenAI, Gemini, Antigravity, GitHub Copilot, Qwen, Vertex AI, iFlow, and Kiro
- **Proxy Server**: Local HTTP proxy that routes requests to appropriate providers
- **Quota Tracking**: Real-time monitoring of quota usage per account
- **Smart Routing**: Round Robin or Fill First strategies for account selection
- **Auto-Failover**: Automatically switches accounts when rate limits are hit
- **Agent Configuration**: Auto-detect and configure AI coding tools
- **Cross-Platform**: Works on Windows and macOS

## Architecture

- **Frontend**: Electron + React + Vite
- **Backend**: Go HTTP server + Proxy server
- **Database**: SQLite for local data storage
- **Communication**: HTTP REST API between Electron and Go backend

## Project Structure

```
quotio-electron-go/
├── backend/                # Go HTTP server + Proxy server
│   ├── cmd/server/         # Main entry point
│   └── internal/
│       ├── api/            # HTTP REST API handlers (Gin router on :8080)
│       ├── proxy/          # HTTP proxy server for request routing
│       ├── storage/        # SQLite database layer (GORM ORM)
│       ├── quota/          # Quota tracking per provider
│       ├── providers/      # Provider integrations (OpenAI, Claude, etc.)
│       ├── agents/         # Agent detection & configuration
│       ├── config/         # Configuration management
│       └── notifications/  # Notification system
├── frontend/               # Electron + React + Vite application
│   └── src/
│       ├── main/           # Electron main process
│       ├── preload/        # Electron preload scripts
│       └── renderer/       # React UI (components, pages, services)
└── AGENTS.md               # Developer guide (build, architecture, code style)
```

**Database**: SQLite (local file storage via GORM ORM)

## Development Setup

### Prerequisites

- Go 1.20 or later
- Node.js 18 or later
- npm or yarn

### Quick Start

**Backend** (Go server on port 8080):
```bash
cd backend
go run cmd/server/main.go
```

**Frontend** (Vite dev server on port 5173):
```bash
cd frontend
npm install
npm run dev
```

**Electron App** (desktop application):
```bash
cd frontend
npm run electron:dev
```

## Building

### Build Go Backend
```bash
cd backend
go build -o quotio-server cmd/server/main.go
```

### Build Electron App
```bash
cd frontend
npm run build
```

This will create platform-specific installers in the `release/` directory.

## Development Workflow

For detailed build commands, testing, architecture info, and code style guidelines, see [AGENTS.md](./AGENTS.md).

## Usage

1. **Start the Application**: Launch Quotio from the application menu
2. **Add Providers**: Go to Providers tab and add your AI provider accounts
3. **Configure Agents**: Go to Agents tab and configure your coding agents to use the proxy
4. **Start Proxy**: Click "Start Proxy" on the Dashboard
5. **Monitor Quota**: Check the Quota tab to see usage across all accounts

## API Endpoints

- `GET /api/health` - Health check
- `GET /api/dashboard` - Dashboard statistics
- `GET /api/providers` - List providers
- `POST /api/providers` - Add provider
- `PUT /api/providers/:id` - Update provider
- `DELETE /api/providers/:id` - Delete provider
- `GET /api/quota` - Get quota status
- `GET /api/agents` - List agents
- `POST /api/agents/configure` - Configure agent
- `POST /api/proxy/start` - Start proxy
- `POST /api/proxy/stop` - Stop proxy
- `GET /api/proxy/status` - Get proxy status
- `GET /api/settings` - Get settings
- `PUT /api/settings` - Update settings

## License

MIT

