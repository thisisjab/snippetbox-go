package main

import (
	"database/sql"
	"flag"
	"log"
	"math"
	"net/http"
	"os"
	"web-dev-journey/cmd/web/config"
	"web-dev-journey/cmd/web/db"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	dbConn   *sql.DB
	config   *config.Config
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	doMigrate := flag.Bool("doMigrate", false, "Run migrations")
	migrationTarget := flag.Int("migrationTarget", 0, "Migrations target: Negative values mean downgrade.")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.LUTC)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.LUTC|log.Llongfile)

	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
	}

	app.loadConfig()
	app.connectDB()
	app.migrateDB(doMigrate, migrationTarget)

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()

	errorLog.Fatal(err)
}

func (app *application) loadConfig() {
	c, configErr := config.LoadConfig()

	if configErr != nil {
		app.errorLog.Fatalf("Error loading config: %v", configErr)
	}

	app.config = c
}

func (app *application) connectDB() {
	conn, connErr := db.OpenDB(app.config.DatabasePath())

	if connErr != nil {
		app.errorLog.Fatalf("Error connecting to database: %v", connErr)
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
			app.errorLog.Fatalf("Error running migrations: %v", migrationErr)
		}

		app.infoLog.Print("Applied migrations")
	}

}
