package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"math"
	"net/http"
	"os"
	"web-dev-journey/cmd/web/config"
	"web-dev-journey/cmd/web/db"
)

type application struct {
	logger *slog.Logger
	dbConn *sql.DB
	config *config.Config
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

	app := &application{
		logger: logger,
	}

	app.loadConfig()
	app.connectDB()
	app.migrateDB(doMigrate, migrationTarget)

	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
		// TODO: add slog here
	}

	logger.Info("Starting server on %s", slog.Any("addr", *addr))
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

func (app *application) connectDB() {
	conn, connErr := db.OpenDB(app.config.DatabasePath())

	if connErr != nil {
		app.logger.Error("Error connecting to database: %v", connErr)
	}

	app.dbConn = conn
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
