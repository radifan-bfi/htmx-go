package main

import (
	"context"
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"htmx-go/handlers"
	"htmx-go/repositories"
	"htmx-go/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Secure())
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 60 * time.Second,
	}))

	// Initialize templates
	templates := utils.InitializeTemplates()
	renderer := utils.NewTemplateRenderer(templates)
	e.Renderer = renderer

	// Init DB connection
	db, err := sql.Open("sqlite3", "./forms.db")
	if err != nil {
		log.Fatal(err)
	}

	// Run migrations
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"sqlite3",
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	// Init repositories
	formSchemaRepository := repositories.NewFormSchemaRepository(db)
	formSubmissionRepository := repositories.NewFormSubmissionRepository(db)

	// Init views handlers
	handlers.NewFormSchemaViewHandler(e, formSchemaRepository, formSubmissionRepository)
	handlers.NewFormSubmissionViewHandler(e, formSubmissionRepository, formSchemaRepository)

	// Init API handlers
	apiV1 := e.Group("/api/v1")
	handlers.NewFormSchemaHandler(apiV1, formSchemaRepository)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Start server
	go func() {
		if err := e.Start(":1323"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
