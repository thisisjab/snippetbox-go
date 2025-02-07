package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type MigrationSet struct {
	basePath   string
	migrations []migration
}

type migration struct {
	upFileName   string
	downFileName string
	name         string
	version      int
}

// RunMigrations takes a db connection, target migration version, and direction (upgrade/downgrade) and runs migrations.
// It's assumed that migration files are sorted based on their target version.
func (ms *MigrationSet) RunMigrations(db *sql.DB, targetVersion int, upgrade bool) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	for _, m := range ms.migrations {

		var migrationFilePath string
		if upgrade {
			if targetVersion < m.version {
				continue
			}

			migrationFilePath = filepath.Join(ms.basePath, m.upFileName)
		} else {
			if targetVersion >= m.version {
				continue
			}

			migrationFilePath = filepath.Join(ms.basePath, m.downFileName)
		}

		migrationSQL, err := os.ReadFile(migrationFilePath)
		if err != nil {
			return fmt.Errorf("reading migration file %s: %w", migrationFilePath, err)
		}

		_, err = tx.Exec(string(migrationSQL))
		if err != nil {
			return fmt.Errorf("executing migration %s: %w", migrationFilePath, err)
		}
	}

	return tx.Commit()
}

// LoadMigrations finds all migrations in given path.
func (ms *MigrationSet) LoadMigrations(path string) error {

	files, readDirErr := os.ReadDir(path)
	if readDirErr != nil {
		return readDirErr
	}

	migrations := make(map[int]migration, 0)

	re := regexp.MustCompile(`^(\d+)_.*(\.sql|_revert\.sql)$`)

	for _, file := range files {
		matches := re.FindStringSubmatch(file.Name())
		if matches != nil {
			version, _ := strconv.Atoi(matches[1])
			isRevert := strings.HasSuffix(file.Name(), "_revert.sql")
			upgrade := !isRevert

			if _, ok := migrations[version]; !ok {
				name := strings.TrimSuffix(file.Name(), ".sql")
				if upgrade {
					name = strings.TrimSuffix(name, "_revert")
					name = strings.TrimPrefix(name, strconv.Itoa(version)+"_")

					migrations[version] = migration{
						upFileName: file.Name(),
						version:    version,
						name:       name,
					}
				} else {
					name = strings.TrimPrefix(name, strconv.Itoa(version)+"_")

					migrations[version] = migration{
						downFileName: file.Name(),
						version:      version,
						name:         name,
					}
				}
			} else {
				temp := migrations[version]

				if upgrade {
					temp.upFileName = file.Name()
				} else {
					temp.downFileName = file.Name()
				}

				migrations[version] = temp
			}
		}
	}

	migrationsList := make([]migration, 0)
	for _, m := range migrations {
		migrationsList = append(migrationsList, m)
	}

	sort.Slice(migrationsList, func(i, j int) bool {
		return migrationsList[i].version < migrationsList[j].version
	})

	ms.basePath = path
	ms.migrations = migrationsList

	return nil
}
