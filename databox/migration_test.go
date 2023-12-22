package databox_test

import (
	"testing"

	"github.com/dominik-da-rocha/go-toolbox/databox"
	"github.com/golang-migrate/migrate"
	"github.com/stretchr/testify/assert"
)

const currentMigVersion = 2

func newFakeMigration() *migrate.Migrate {
	db, name := databox.OpenMemoryDb()
	mig := databox.NewMigration(db, name, "file://")
	return mig
}

func Test_NewMigration(t *testing.T) {
	assert.NotPanics(t, func() {
		mig := newFakeMigration()
		assert.NotNil(t, mig)
	})
}

func Test_AutoMigration(t *testing.T) {
	assert.NotPanics(t, func() {
		db, name := databox.OpenMemoryDb()
		databox.AutoMigration(db, name, "file://")
	})
}

func Test_GetMigrationVersion(t *testing.T) {
	mig := newFakeMigration()
	ver, dirty := databox.GetMigrationVersion(mig)
	assert.Equal(t, uint(0), ver)
	assert.Equal(t, false, dirty)
}

func Test_MigrateUp(t *testing.T) {
	mig := newFakeMigration()
	databox.MigrateUp(mig)
	ver, dirty := databox.GetMigrationVersion(mig)
	assert.Equal(t, uint(currentMigVersion), ver)
	assert.Equal(t, false, dirty)
}
func Test_MigrateDown(t *testing.T) {
	mig := newFakeMigration()
	databox.MigrateUp(mig)
	ver, dirty := databox.GetMigrationVersion(mig)
	assert.Equal(t, uint(currentMigVersion), ver)
	assert.Equal(t, false, dirty)

	databox.MigrateDown(mig)
	ver, dirty = databox.GetMigrationVersion(mig)
	assert.Equal(t, uint(0), ver)
	assert.Equal(t, false, dirty)
}
func Test_MigrateTo(t *testing.T) {
	mig := newFakeMigration()
	databox.MigrateTo(mig, 1)
	ver, dirty := databox.GetMigrationVersion(mig)
	assert.Equal(t, uint(1), ver)
	assert.Equal(t, false, dirty)
}
