# AGENTS.md - Quotio Electron Go

## Build & Test
**Backend**: `cd backend && go build -o quotio-server cmd/server/main.go` | Run: `go run cmd/server/main.go` | Test: `go test ./...` | Single: `go test ./internal/quota -v` | Lint: `go fmt ./...`
**Frontend**: `cd frontend && npm run dev` (Vite :5173) | Electron: `npm run electron:dev` | Build: `npm run build`

## Architecture
**Stack**: Go (Gin + SQLite) backend + Electron + React + Vite frontend
**Backend**: `cmd/server/main.go` → `internal/{api,proxy,storage,quota,providers,agents,config,notifications}`
**Frontend**: `src/{main,preload,renderer}` with React Router, Axios, Zustand; API: `http://localhost:8080/api/*`
**Database**: SQLite (GORM ORM)

## Code Style
**Go**: CamelCase (exported), error returns, `log.Fatalf()` on startup fail, handlers: `func (s *Server) handleXxx(c *gin.Context)`, routes in `api/server.go`
**JS**: `.jsx` for components, functional + hooks, Zustand state, Axios for API calls
**Naming**: Go exports CamelCase; JS components PascalCase; describe purpose
**Imports**: Full paths `quotio-electron-go/backend/internal/...`, use Gin/GORM/React/Axios
**API**: `/api/*` endpoints, JSON, CORS enabled for Electron
