// main.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/joho/godotenv"

	"tech-test/backend/internal/config"
	"tech-test/backend/internal/database"
	"tech-test/backend/internal/handler"
	"tech-test/backend/internal/middleware"
	"tech-test/backend/internal/repository/memory"
	"tech-test/backend/internal/repository/sqlite"
	userService "tech-test/backend/internal/service/user"
	fileService "tech-test/backend/internal/service/file"
	_ "tech-test/backend/docs" // Swagger docs
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Application struct {
	config     *config.Config
	logger     *zap.Logger
	httpServer *http.Server
	router     *mux.Router
	db         *gorm.DB
}

func main() {
	// Initialize context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app, err := initializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer app.cleanup()

	if err := app.run(ctx); err != nil {
		app.logger.Fatal("Failed to run application", zap.Error(err))
	}
}

func initializeApp() (*Application, error) {
	// Initialize logger first for proper error reporting
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Load configuration
	if err := godotenv.Load(); err != nil {
		logger.Warn("No .env file found")
	}

	cfg := config.NewConfig()

	// Initialize application struct
	app := &Application{
		config: cfg,
		logger: logger,
		router: mux.NewRouter(),
	}

	// Setup dependencies
	if err := app.setupDependencies(); err != nil {
		return nil, fmt.Errorf("failed to setup dependencies: %w", err)
	}

	return app, nil
}

func (app *Application) setupDependencies() error {
	// Setup database with proper error handling
	db, err := database.SetupDB(app.config.Database)
	if err != nil {
		return fmt.Errorf("database setup failed: %w", err)
	}
	app.db = db

	// Verify database connection
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	// Setup repositories with dependency injection
	userRepo := memory.NewUserRepository()
	fileRepo := sqlite.NewFileRepository(db)

	// Setup upload directory with proper permissions
	uploadDir := filepath.Join(".", "uploads")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Setup services with proper dependency injection
	userService := userService.NewService(userRepo, app.logger)
	fileService := fileService.NewService(
		fileRepo,
		app.logger,
		uploadDir,
	)

	// Setup handlers with dependencies
	app.setupRoutes(
		handler.NewAuthHandler(userService),
		handler.NewFileHandler(
			fileService,
			app.config.File,
		),
		handler.NewUserHandler(userService),
	)

	return nil
}

func (app *Application) setupRoutes(
	authHandler *handler.AuthHandler,
	fileHandler *handler.FileHandler,
	userHandler *handler.UserHandler,
) {
	app.router.Use(middleware.CORS(app.logger))
	app.router.Use(middleware.RateLimiterMiddleware())
	app.router.Use(middleware.RequestLogger(app.logger))
	app.router.Use(middleware.SecurityHeaders())

	app.router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	app.router.HandleFunc("/health", app.healthCheck).Methods(http.MethodGet)
	app.router.HandleFunc("/api/register", middleware.ValidateRegister(authHandler.Register)).Methods(http.MethodPost, http.MethodOptions)
	app.router.HandleFunc("/api/login", authHandler.Login).Methods(http.MethodPost, http.MethodOptions)
	app.router.HandleFunc("/shared/{shareId}", fileHandler.GetSharedFile).Methods(http.MethodGet, http.MethodOptions)

	protected := app.router.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	files := protected.PathPrefix("/files").Subrouter()
	files.HandleFunc("/upload", fileHandler.Upload).Methods(http.MethodPost, http.MethodOptions)
	files.HandleFunc("/search", fileHandler.SearchFiles).Methods(http.MethodGet, http.MethodOptions)
	files.HandleFunc("/my", fileHandler.GetUserFiles).Methods(http.MethodGet, http.MethodOptions)
	files.HandleFunc("/{id}/download", fileHandler.Download).Methods(http.MethodGet, http.MethodOptions)
	files.HandleFunc("/{id}/view", fileHandler.View).Methods(http.MethodGet, http.MethodOptions)
	files.HandleFunc("/{id}", fileHandler.GetByID).Methods(http.MethodGet, http.MethodOptions)
	files.HandleFunc("/{id}", fileHandler.Delete).Methods(http.MethodDelete, http.MethodOptions)
	files.HandleFunc("", fileHandler.List).Methods(http.MethodGet, http.MethodOptions)
	files.HandleFunc("/{id}/share", fileHandler.GenerateShareableLink).Methods(http.MethodPost, http.MethodOptions)

	users := protected.PathPrefix("/users").Subrouter()
	users.HandleFunc("", userHandler.GetAllUsers).Methods(http.MethodGet, http.MethodOptions)
	users.HandleFunc("", userHandler.CreateUser).Methods(http.MethodPost, http.MethodOptions)
	users.HandleFunc("/{id}", userHandler.GetUser).Methods(http.MethodGet, http.MethodOptions)
	users.HandleFunc("/{id}", userHandler.UpdateUser).Methods(http.MethodPut, http.MethodOptions)
	users.HandleFunc("/{id}", userHandler.DeleteUser).Methods(http.MethodDelete, http.MethodOptions)
	users.HandleFunc("/me", userHandler.GetCurrentUser).Methods(http.MethodGet, http.MethodOptions)
}

func (app *Application) run(ctx context.Context) error {
	// Enhanced server configuration
	app.httpServer = &http.Server{
		Addr:    ":" + app.config.Port,
		Handler: app.router,
		// Enhanced timeouts
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,  // Add header timeout
	}

	// Create shutdown channel
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Server startup logging
	app.logger.Info("Starting server",
		zap.String("port", app.config.Port),
		zap.String("env", app.config.Environment),
		zap.String("version", "1.0.0"))

	// Start server in goroutine
	serverErrors := make(chan error, 1)
	go func() {
		if err := app.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- fmt.Errorf("server error: %w", err)
		}
	}()

	// Wait for shutdown signal or server error
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		app.logger.Info("Shutdown signal received", zap.String("signal", sig.String()))

		// Create shutdown context with timeout
		shutdownCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
		defer cancel()

		// Graceful shutdown
		if err := app.httpServer.Shutdown(shutdownCtx); err != nil {
			// Force shutdown if graceful shutdown fails
			app.httpServer.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	app.logger.Info("Server stopped gracefully")
	return nil
}

func (app *Application) cleanup() {
	// Cleanup database connection
	if app.db != nil {
		if err := database.CloseDB(app.db); err != nil {
			app.logger.Error("Failed to close database connection", zap.Error(err))
		}
	}

	// Ignore sync errors on cleanup
	_ = app.logger.Sync()
}

func (app *Application) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"version": "1.0.0",
	})
}

