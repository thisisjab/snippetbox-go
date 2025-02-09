package main

import (
	"database/sql"
	"flag"
	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"html/template"
	"log/slog"
	"math"
	"net/http"
	"os"
	"time"
	"web-dev-journey/cmd/web/config"
	"web-dev-journey/cmd/web/db"
	"web-dev-journey/internal/models"
)

type application struct {
	config         *config.Config
	dbConn         *sql.DB
	formDecoder    *form.Decoder
	logger         *slog.Logger
	sessionManager *scs.SessionManager
	snippets       *models.SnippetModel
	templateCache  map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	doMigrate := flag.Bool("doMigrate", false, "Run migrations")
	migrationTarget := flag.Int("migrationTarget", 0, "Migrations target: Negative values mean downgrade.")

	flag.Parse()

	loggerHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})
	logger := slog.New(loggerHandler)

	templateCache, tcErr := newTemplateCache()
	if tcErr != nil {
		logger.Error(tcErr.Error())
		os.Exit(1)
	}

	fd := form.NewDecoder()

	app := &application{
		logger:        logger,
		formDecoder:   fd,
		templateCache: templateCache,
	}

	app.loadConfig()
	app.connectDBModels()
	app.migrateDB(doMigrate, migrationTarget)
	app.setupSessionManager()

	srv := &http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}

	logger.Info("Starting server", slog.Any("addr", *addr))
	err := srv.ListenAndServe()

	app.logger.Error(err.Error())
	os.Exit(1)
}

func (app *application) loadConfig() {
	c, configErr := config.LoadConfig()

	if configErr != nil {
		app.logger.Error("Error loading config: %v", configErr)
	}

	app.config = c
}

func (app *application) connectDBModels() {
	conn, connErr := db.OpenDB(app.config.DatabasePath())

	if connErr != nil {
		app.logger.Error("Error connecting to database: %v", connErr)
	}

	app.dbConn = conn
	app.snippets = &models.SnippetModel{DB: conn}
}

func (app *application) migrateDB(doMigrate *bool, target *int) {
	if *doMigrate {

		ms := &db.MigrationSet{}
		ms.LoadMigrations(app.config.MigrationsPath())

		upgrade := *target >= 0

		target := int(math.Abs(float64(*target)))

		migrationErr := ms.RunMigrations(app.dbConn, target, upgrade)

		if migrationErr != nil {
			app.logger.Error("Error running migrations: %v", migrationErr)
		}

		app.logger.Info("Applied migrations")
	}

}

func (app *application) setupSessionManager() {
	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(app.dbConn)
	sessionManager.Lifetime = 12 * time.Hour

	app.sessionManager = sessionManager
}
