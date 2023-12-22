package databox_test

import (
	"database/sql"
	"fmt"
	"log/slog"
	"reflect"

	"github.com/dominik-da-rocha/go-toolbox/databox"
	"github.com/dominik-da-rocha/go-toolbox/toolbox"
)

type FakeEntity struct {
	Id     int64    `db:"id"`
	Flag   bool     `db:"flag"`
	Real   float64  `db:"real"`
	Text   string   `db:"text"`
	Blob   []byte   `db:"blob"`
	Ignore []string `db:"-"`
}

func createFakeTable(db *sql.DB) {
	query := `
	CREATE TABLE "fake" (
		"id"     INTEGER PRIMARY KEY AUTOINCREMENT,
		"flag"	 BOOLEAN NOT NULL, 
		"real"	 REAL NOT NULL,
		"text"   TEST NOT NULL,
		"blob"   BLOB NOT NULL
	);`
	_, err := db.Exec(query)
	toolbox.Uups(err)
	slog.Info("created fake table")
}

func newFakeEntities(count int) []any {
	fakes := []any{}
	for i := 1; i <= count; i++ {
		f := newFakeEntity(i)
		fakes = append(fakes, &f)
	}
	return fakes
}

func FakeStore() databox.Store {
	db := databox.NewMigratedMemoryDB("")
	createFakeTable(db)
	dia := databox.NewDialect()
	adapter := databox.StoreAdapter{
		Db:      db,
		Dialect: dia,
		Table:   "fake",
		Scan: func(row any) []any {
			var ptr *FakeEntity
			switch d := row.(type) {
			case []*FakeEntity:
				ptr = d[0]
			case []FakeEntity:
				ptr = &d[0]
			case *FakeEntity:
				ptr = d
			}
			if ptr != nil {
				return []any{
					&ptr.Id,
					&ptr.Flag,
					&ptr.Real,
					&ptr.Text,
					&ptr.Blob,
				}
			}
			slog.Error("not a *FakeEntity got", "type", reflect.TypeOf(row))
			panic("not a *FakeEntity got")
		},
		Factory:     func() any { return &FakeEntity{} },
		RevisionKey: "fake",
	}
	store := databox.NewStore(adapter)
	return store
}

func newFakeEntity(i int) FakeEntity {
	f := FakeEntity{
		Id:   0,
		Flag: i%2 == 0,
		Real: 100.123 + float64(i),
		Text: fmt.Sprintf("Hello"),
		Blob: []byte(fmt.Sprintf("World")),
	}
	return f
}

func arrangeStoreAndData() (databox.Store, FakeEntity) {
	store := FakeStore()
	f := newFakeEntity(0)
	err := store.Insert(&f)
	toolbox.Uups(err)
	return store, f
}
