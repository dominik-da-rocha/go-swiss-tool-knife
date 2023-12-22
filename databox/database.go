package databox

import (
	"database/sql"
	"log/slog"

	"github.com/dominik-da-rocha/go-toolbox/toolbox"
)

func OpenDb(name string) *sql.DB {
	slog.Info("Open database", "name", name)
	db, err := sql.Open("sqlite3", name)
	toolbox.Uups(err)

	return db
}

func OpenMemoryDb() (*sql.DB, string) {
	name := ":memory:"
	slog.Info("Open database", "name", name)
	db, err := sql.Open("sqlite3", name)
	toolbox.Uups(err)
	return db, name
}
