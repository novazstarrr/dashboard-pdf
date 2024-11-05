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
	_ "tech-test/backend/docs" 
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
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	if err := godotenv.Load(); err != nil {
		logger.Warn("No .env file found")
	}

	cfg := config.NewConfig()

	app := &Application{
		config: cfg,
		logger: logger,
		router: mux.NewRouter(),
	}

	if err := app.setupDependencies(); err != nil {
		return nil, fmt.Errorf("failed to setup dependencies: %w", err)
	}

	return app, nil
}

func (app *Application) setupDependencies() error {
	db, err := database.SetupDB(app.config.Database)
	if err != nil {
		return fmt.Errorf("database setup failed: %w", err)
	}
	app.db = db

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	userRepo := memory.NewUserRepository()
	fileRepo := sqlite.NewFileRepository(db)

	uploadDir := filepath.Join(".", "uploads")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return fmt.Errorf("failed to create upload directory: %w", err)
	}

	userService := userService.NewService(userRepo, app.logger)
	fileService := fileService.NewService(
		fileRepo,
		app.logger,
		uploadDir,
	)

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

	protected.HandleFunc("/users/me", userHandler.GetCurrentUser).Methods(http.MethodGet, http.MethodOptions)

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
}

func (app *Application) run(ctx context.Context) error {
	app.httpServer = &http.Server{
		Addr:    ":" + app.config.Port,
		Handler: app.router,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,  
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	app.logger.Info("Starting server",
		zap.String("port", app.config.Port),
		zap.String("env", app.config.Environment),
		zap.String("version", "1.0.0"))

	serverErrors := make(chan error, 1)
	go func() {
		if err := app.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- fmt.Errorf("server error: %w", err)
		}
	}()

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		app.logger.Info("Shutdown signal received", zap.String("signal", sig.String()))

		shutdownCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
		defer cancel()

		if err := app.httpServer.Shutdown(shutdownCtx); err != nil {
			app.httpServer.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	app.logger.Info("Server stopped gracefully")
	return nil
}

func (app *Application) cleanup() {
	if app.db != nil {
		if err := database.CloseDB(app.db); err != nil {
			app.logger.Error("Failed to close database connection", zap.Error(err))
		}
	}

	
	_ = app.logger.Sync()
}

func (app *Application) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"version": "1.0.0",
	})
}

