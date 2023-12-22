package databox_test

import (
	"testing"

	"github.com/dominik-da-rocha/go-toolbox/databox"
	"github.com/stretchr/testify/assert"
)

func newRevisionStore() databox.RevisionStore {
	database := databox.NewMigratedMemoryDB("")
	store := databox.NewRevisionStore(database)
	return store
}

func Test_Revision_IncrRev(t *testing.T) {
	store := newRevisionStore()
	defer store.Close()
	rev, err := store.IncrRev("test")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), rev)

	val, err := store.GetRev("test")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), val)
}
