package api

import (
	"fmt"
	"os"
	"quotio-electron-go/backend/internal/config"
	"quotio-electron-go/backend/internal/proxy"
	
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	db     *gorm.DB
	config *config.Config
	proxy  *proxy.Server
	router *gin.Engine
}

func NewServer(db *gorm.DB, cfg *config.Config) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// CORS configuration for Electron app
	// Restrict to localhost and file:// protocol (Electron renderer)
	allowedOrigins := getAllowedOrigins(cfg)
	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	server := &Server{
		db:     db,
		config: cfg,
		router: router,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	api := s.router.Group("/api")
	{
		// Health check
		api.GET("/health", s.handleHealth)
		
		// Dashboard
		api.GET("/dashboard", s.handleDashboard)
		
	// Providers
	api.GET("/providers", s.handleGetProviders)
	api.GET("/providers/status", s.handleGetProviderStatus)
	api.GET("/providers/detect", s.handleDetectProviderAccounts)
	api.GET("/providers/health", s.handleGetProviderHealth)
	api.POST("/providers/health/:id", s.handleCheckProviderHealth)
	api.POST("/providers", s.handleAddProvider)
	api.PUT("/providers/:id", s.handleUpdateProvider)
	api.DELETE("/providers/:id", s.handleDeleteProvider)

	// Quota
	api.GET("/quota", s.handleGetQuota)
	api.GET("/quota/history/:id", s.handleGetQuotaHistory)
	api.GET("/quota/failed", s.handleGetFailedRequests)
	api.POST("/quota/reset/:id", s.handleResetQuota)
	api.GET("/models", s.handleGetModels)

	// Routing
	api.POST("/routing-strategy", s.handleUpdateRoutingStrategy)
	api.GET("/rate-limits", s.handleGetRateLimits)

	// OAuth Detection
	api.GET("/providers/detect-oauth", s.handleDetectOAuthCredentials)
	api.POST("/providers/from-oauth", s.handleAddProviderFromOAuth)

	// Agents
	api.GET("/agents", s.handleGetAgents)
	api.POST("/agents/configure", s.handleConfigureAgent)
	api.POST("/agents/refresh", s.handleRefreshAgents)

	// Proxy
	api.POST("/proxy/start", s.handleStartProxy)
	api.POST("/proxy/stop", s.handleStopProxy)
	api.GET("/proxy/status", s.handleProxyStatus)

	// Settings
	api.GET("/settings", s.handleGetSettings)
	api.PUT("/settings", s.handleUpdateSettings)
	}
}

func (s *Server) Start() error {
	return s.router.Run(fmt.Sprintf(":%d", s.config.Port))
}

// getAllowedOrigins returns a restricted list of allowed CORS origins
// For Electron apps, this includes localhost on the configured port and file:// protocol
func getAllowedOrigins(cfg *config.Config) []string {
	origins := []string{
		// Localhost with configured port (for dev/testing)
		fmt.Sprintf("http://localhost:%d", cfg.Port),
		fmt.Sprintf("http://127.0.0.1:%d", cfg.Port),
		"http://localhost:5173",     // Common Vite dev server port
		"http://127.0.0.1:5173",     // For testing with local dev server
		"http://localhost:3000",     // Electron/React dev server port
		"http://127.0.0.1:3000",     // For testing with local dev server on port 3000
	}

	// Allow custom origins via environment variable (for advanced/testing scenarios)
	if customOrigins := os.Getenv("QUOTIO_ALLOWED_ORIGINS"); customOrigins != "" {
		// Note: This should only be used for development. In production, keep the default list.
		// Format: "http://example.com,http://another.com"
		// This is intentionally restrictive - Electron apps should use localhost
	}

	return origins
}

