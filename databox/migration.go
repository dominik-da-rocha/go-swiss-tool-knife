package databox

import (
	"database/sql"
	"log/slog"

	"github.com/dominik-da-rocha/go-toolbox/toolbox"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/sqlite3"

	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func NewMigration(db *sql.DB, databaseName string, migrationUrl string) *migrate.Migrate {
	slog.Info("New Migration")
	slog.Info("Configure migration driver sqlite")
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	toolbox.Uups(err)

	slog.Info("New migration with database", "db", databaseName, "mig-url", migrationUrl)
	m, err := migrate.NewWithDatabaseInstance(migrationUrl, databaseName, driver)
	toolbox.Uups(err)

	return m
}

func AutoMigration(db *sql.DB, databaseName string, migrationUrl string) {
	slog.Info("Auto migration")
	mig := NewMigration(db, databaseName, migrationUrl)

	version, isDirty := GetMigrationVersion(mig)

	if isDirty && version > 0 {
		panic("database version is dirty abort migration")
	}

	MigrateUp(mig)
}

func GetMigrationVersion(mig *migrate.Migrate) (uint, bool) {
	version, isDirty, err := mig.Version()
	if err == migrate.ErrNilVersion {
		slog.Warn("get mig version", "msg", err)
	} else if err != nil {
		toolbox.Uups(err)
	}
	status := "clean"
	if isDirty {
		status = "dirty"
	}
	slog.Info("migration version", "v", version, "s", status)
	return version, isDirty
}

func MigrateUp(mig *migrate.Migrate) {
	slog.Info("migration up")
	err := mig.Up()
	if err == migrate.ErrNoChange {
		slog.Warn("mig up", "msg", err.Error())
	} else if err != nil {
		toolbox.Uups(err)
	} else {
		GetMigrationVersion(mig)
	}
}

func MigrateDown(mig *migrate.Migrate) {
	slog.Info("migration down")
	err := mig.Down()
	if err == migrate.ErrNoChange {
		slog.Warn("mig down", "msg", err.Error())
	} else if err != nil {
		toolbox.Uups(err)
	} else {
		GetMigrationVersion(mig)
	}
}

func MigrateTo(mig *migrate.Migrate, version uint) {
	slog.Info("migration to", "v", version)
	err := mig.Migrate(version)
	toolbox.Uups(err)
}

func NewMigratedMemoryDB(pathToMigrations string) *sql.DB {
	name := ":memory:"
	db, err := sql.Open("sqlite3", name)
	toolbox.Uups(err)
	AutoMigration(db, name, "file://"+pathToMigrations)
	return db
}
