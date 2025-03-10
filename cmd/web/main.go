package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/thisisjab/snippetbox-go/cmd/web/config"
	"github.com/thisisjab/snippetbox-go/cmd/web/db"
	"github.com/thisisjab/snippetbox-go/internal/model"
	"html/template"
	"log/slog"
	"math"
	"net/http"
	"os"
	"time"
)

type application struct {
	config         *config.Config
	dbConn         *sql.DB
	formDecoder    *form.Decoder
	logger         *slog.Logger
	sessionManager *scs.SessionManager
	snippets       model.SnippetModelInterface
	templateCache  map[string]*template.Template
	users          model.UserModelInterface
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	useTLS := flag.Bool("tls", false, "Connection uses TLS if true. Corresponding key and certificate path must be set as env vars.")
	doMigrate := flag.Bool("doMigrate", false, "Run migrations")
	migrationTarget := flag.Int("migrationTarget", 0, "Migrations target: Negative values mean downgrade.")

	flag.Parse()

	app := &application{}

	app.setupLogger()
	app.loadConfig()
	app.connectDBModels()
	app.migrateDB(doMigrate, migrationTarget)
	app.setupSessionManager()
	app.loadTemplates()
	app.setupFormDecoder()

	tlsConfig := &tls.Config{
		// There is no browser that supports TLS 1.3 and does not support SameSite cookies.
		MinVersion: tls.VersionTLS13,
	}

	srv := &http.Server{
		Addr:         *addr,
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		TLSConfig:    tlsConfig,
	}

	app.logger.Info("Starting server", slog.Any("tls", *useTLS), slog.Any("addr", *addr))
	var serverErr error

	if *useTLS {
		serverErr = srv.ListenAndServeTLS(app.config.TLSCertPath(), app.config.TLSKeyPath())
	} else {
		serverErr = srv.ListenAndServe()
	}

	app.logger.Error(serverErr.Error())
	os.Exit(1)
}

func (app *application) loadConfig() {
	c, configErr := config.LoadConfig()

	if configErr != nil {
		app.logger.Error("Error loading config", "error", configErr)
	}

	app.config = c
}

func (app *application) connectDBModels() {
	conn, connErr := db.OpenDB(app.config.DatabasePath())

	if connErr != nil {
		app.logger.Error("Error connecting to database", "error", connErr)
	}

	app.dbConn = conn
	app.snippets = &model.SnippetModel{DB: conn}
	app.users = &model.UserModel{DB: conn}
}

func (app *application) migrateDB(doMigrate *bool, target *int) {
	if *doMigrate {

		ms := &db.MigrationSet{}
		ms.LoadMigrations(app.config.MigrationsPath())

		upgrade := *target >= 0

		target := int(math.Abs(float64(*target)))

		migrationErr := ms.RunMigrations(app.dbConn, target, upgrade)

		if migrationErr != nil {
			app.logger.Error("Error running migrations", "error", migrationErr)
		}

		app.logger.Info("Applied migrations")
	}

}

func (app *application) setupSessionManager() {
	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(app.dbConn)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	app.sessionManager = sessionManager
}

func (app *application) setupLogger() {
	loggerHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})
	logger := slog.New(loggerHandler)

	app.logger = logger
}

func (app *application) loadTemplates() {
	templateCache, tcErr := newTemplateCache()

	if tcErr != nil {
		app.logger.Error(tcErr.Error())
		os.Exit(1)
	}

	app.templateCache = templateCache
}

func (app *application) setupFormDecoder() {
	fd := form.NewDecoder()

	app.formDecoder = fd
}
