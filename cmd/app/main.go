package main

import (
	"context"
	"embed"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"base_app/internal/adapter/auth/inmemory"
	"base_app/internal/adapter/repository/postgresql"
	"base_app/internal/config"
	apiHandler "base_app/internal/handler/http"
	v1 "base_app/internal/handler/http/v1"
	"base_app/internal/service"
	"base_app/internal/usecase"
	"base_app/pkg/logger"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/getsentry/sentry-go"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed web
var embeddedFiles embed.FS

const (
	appName = "base_app"
)

var (
	mode       string
	configPath string
)

func init() {
	flag.StringVar(&mode, "mode", "Start", "Application run mode. Use 'Prepare' to run migrations.")
	flag.StringVar(&configPath, "config", "configs/config.yaml", "Path to the configuration file.")
}

func main() {
	flag.Parse()

	cfg := config.MustLoad(configPath)

	var log *slog.Logger
	if !cfg.Logger.Enabled {
		log = logger.NewDiscardLogger()
	} else {
		var err error
		log, err = logger.New(cfg.Logger.Level, cfg.Logger.Destination)
		if err != nil {
			log = slog.New(slog.NewJSONHandler(os.Stdout, nil))
			log.Error("failed to initialize file logger, falling back to stdout", "error", err)
		}
	}
	log.Info("starting application", slog.String("app", appName), slog.String("mode", mode))

	switch mode {
	case "Prepare":
		runMigrations(cfg.Postgres, log)
	case "Start":
		runApp(cfg, log)
	default:
		log.Error("invalid mode specified", slog.String("mode", mode))
		os.Exit(1)
	}
}

func runApp(cfg *config.Config, log *slog.Logger) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if cfg.Sentry.Enabled {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:              cfg.Sentry.Dsn,
			TracesSampleRate: 1.0,
		}); err != nil {
			log.Error("failed to initialize sentry", "error", err)
		}
		defer sentry.Flush(2 * time.Second)
		log.Info("sentry is enabled")
	} else {
		log.Info("sentry is disabled")
	}

	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.Secure = false
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode

	if cfg.Redis.Enabled {
		redisPool := &redis.Pool{
			MaxIdle: 10,
			Dial: func() (redis.Conn, error) {
				return redis.Dial("tcp", cfg.Redis.Host+":"+cfg.Redis.Port)
			},
		}
		sessionManager.Store = redisstore.New(redisPool)
		log.Info("redis is configured as the session store")
	} else {
		log.Info("redis is disabled, using in-memory session store")
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DBName, cfg.Postgres.SSLMode)
	pgClient, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Error("failed to connect to postgres", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pgClient.Close()
	log.Info("successfully connected to postgres")

	var authService usecase.AuthService
	switch cfg.Auth.Provider {
	case "inmemory":
		authService = inmemory.New(log)
		log.Info("using in-memory auth provider")
	case "postgres":
		repo := postgresql.NewRepo(pgClient, log)
		authService = service.NewAuthService(repo, log)
		log.Info("using postgres auth provider")
	default:
		log.Error("invalid auth provider specified", "provider", cfg.Auth.Provider)
		os.Exit(1)
	}

	repo := postgresql.NewRepo(pgClient, log)
	dataService := service.NewDataService(repo, log)
	catalogService := service.NewCatalogService(repo, log)
	authUsecase := usecase.NewAuthUsecase(authService, log)
	dataUsecase := usecase.NewDataUsecase(dataService, log)
	catalogUsecase := usecase.NewCatalogUsecase(catalogService, log)

	contentFS, err := fs.Sub(embeddedFiles, "web")
	if err != nil {
		log.Error("failed to create sub-filesystem for embedded files", "error", err)
		os.Exit(1)
	}

	handler := apiHandler.NewHandler(authUsecase, dataUsecase, catalogUsecase, sessionManager, contentFS)

	ogenServer, err := v1.NewServer(handler, handler)
	if err != nil {
		log.Error("failed to create ogen server", "error", err)
		os.Exit(1)
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(sessionManager.LoadAndSave)

	router.Mount("/api/v1", ogenServer)
	router.Get("/*", handler.ServeHTTP)

	server := &http.Server{
		Addr:         cfg.HTTP.Host + ":" + cfg.HTTP.Port,
		Handler:      router,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		url := fmt.Sprintf("%s://%s:%s", "http", cfg.HTTP.Host, cfg.HTTP.Port)
		log.Info("http server starting", slog.String("addr", url))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("http server error", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	<-stop

	log.Info("shutting down server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("server shutdown failed", slog.String("error", err.Error()))
	} else {
		log.Info("server gracefully stopped")
	}
}

func runMigrations(cfg config.PostgresConfig, log *slog.Logger) {
	log.Info("running database migrations")
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)
	m, err := migrate.New("file://contracts/pgsql/migrations", dsn)
	if err != nil {
		log.Error("failed to create migrate instance", slog.String("error", err.Error()))
		os.Exit(1)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Error("failed to apply migrations", slog.String("error", err.Error()))
		os.Exit(1)
	}
	log.Info("migrations applied successfully")
}
