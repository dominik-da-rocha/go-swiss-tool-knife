package databox_test

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	_ "github.com/mattn/go-sqlite3"
)

func Test_Insert(t *testing.T) {
	// Arrange
	store := FakeStore()
	defer store.Close()
	f := newFakeEntity(0)
	// Act
	err := store.Insert(&f)
	// Assert
	assert.NoError(t, err, "Insert does not return any error")
	assert.Equal(t, int64(1), f.Id)
}

func Test_UpdateOneAndSelectOne(t *testing.T) {
	// Arrange
	store, fake := arrangeStoreAndData()
	defer store.Close()
	f := fake
	f.Flag = true
	f.Blob = []byte("foo")
	f.Real = 345.3
	f.Text = "bar"
	// Act
	err := store.Update(&f)
	assert.NoError(t, err, "UpdateById does not return any error")
	res, err := store.SelectById(f.Id)
	// Assert
	assert.NoError(t, err, "SelectById does not return any error")
	copy, ok := res.(*FakeEntity)
	assert.Equal(t, true, ok, "cast to *FakeEntity")
	assert.Equal(t, true, copy.Flag, "value of selected is matching")
	assert.Equal(t, []byte("foo"), copy.Blob, "value of selected is matching")
	assert.Equal(t, 345.3, copy.Real, "value of selected is matching")
	assert.Equal(t, "bar", copy.Text, "value of selected is matching")

}

func Test_InsertAndSelectAllBy(t *testing.T) {
	// Arrange
	store := FakeStore()
	defer store.Close()
	fakes := newFakeEntities(10)
	// Act
	err := store.InsertAll(fakes...)
	assert.NoError(t, err, "InsertAll no error expected")
	copies, err := store.SelectAllBy("text", "Hello")
	// Assert
	assert.NoError(t, err, "SelectAllBy no error expected")
	assert.Equal(t, 10, len(copies), "All inserted are found")
	for idx, copy := range copies {
		fake, ok := copy.(*FakeEntity)
		assert.Equal(t, true, ok, "cast to *FakeEntity is okay")
		assert.Equal(t, int64(idx+1), fake.Id, "all ordered by id starting with 1")
		assert.Equal(t, 101.123+float64(idx), fake.Real, "inserted values do match")
	}
}

func Test_UpdateAllAndSelectByPage(t *testing.T) {
	// Arrange
	store := FakeStore()
	defer store.Close()
	fakes := newFakeEntities(100)
	err := store.InsertAll(fakes...)
	assert.NoError(t, err, "InsertAll no error expected")
	for id, fake := range fakes {
		entity := fake.(*FakeEntity)
		entity.Real = 1000 - float64(id)/10
	}

	// Act
	err = store.UpdateAll(fakes...)
	assert.NoError(t, err, "UpdateAll no error expected")

	offset := int64(20)
	limit := int64(10)
	orderBy := "real"
	sortDir := "asc"

	copies, total, err := store.SelectAllByPage(&offset, &limit, &orderBy, &sortDir, nil)
	// Assert
	assert.NoError(t, err, "SelectAllByPage no error expected")
	assert.Equal(t, int64(100), total)
	assert.Equal(t, 10, len(copies))
	for idx, copy := range copies {
		fake := copy.(*FakeEntity)
		assert.Equal(t, int64(80-idx), fake.Id)
	}
}

func Test_SelectAllByPageWithLike(t *testing.T) {
	// Arrange
	store := FakeStore()
	defer store.Close()
	fakes := newFakeEntities(100)
	for id, fake := range fakes {
		entity := fake.(*FakeEntity)
		entity.Real = 1000 - float64(id)/10
		if id%2 == 0 {
			entity.Text = "even"
		} else {
			entity.Text = "odd"
		}
	}
	err := store.InsertAll(fakes...)
	assert.NoError(t, err, "InsertAll no error expected")

	offset := int64(20)
	limit := int64(10)
	orderBy := "id"
	sortDir := "asc"
	search := "even"
	// Act
	copies, total, err := store.SelectAllByPage(&offset, &limit, &orderBy, &sortDir, &search)
	// Assert
	assert.NoError(t, err, "SelectAllByPage no error expected")
	assert.Equal(t, int64(50), total)
	assert.Equal(t, 10, len(copies))
	even := int64(41)
	for _, copy := range copies {
		fake := copy.(*FakeEntity)
		assert.Equal(t, int64(even), fake.Id)
		even += 2
	}
}

func Test_DeleteById(t *testing.T) {
	// Arrange
	store := FakeStore()
	defer store.Close()
	fakes := newFakeEntities(100)
	err := store.InsertAll(fakes...)
	assert.NoError(t, err)
	// Act
	deleted, err := store.DeleteById(10)
	assert.NoError(t, err)
	assert.Equal(t, true, deleted)
	// Assert
	_, err = store.SelectById(10)
	assert.ErrorIs(t, sql.ErrNoRows, err)
}

func Test_DeleteAllBy(t *testing.T) {
	// Arrange
	store := FakeStore()
	defer store.Close()
	fakes := newFakeEntities(100)
	err := store.InsertAll(fakes...)
	assert.NoError(t, err)
	// Act
	deleted, err := store.DeleteAllBy("flag", true)
	assert.NoError(t, err)
	assert.Equal(t, int64(50), deleted)
	// Assert
	copies, err := store.SelectAllBy("flag", true)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(copies))
}

func Test_DeleteFromBy(t *testing.T) {
	// Arrange
	store := FakeStore()
	defer store.Close()
	fakes := newFakeEntities(100)
	err := store.InsertAll(fakes...)
	assert.NoError(t, err)
	// Act
	deleted, err := store.DeleteFromBy("fake", "flag", true)
	assert.NoError(t, err)
	assert.Equal(t, int64(50), deleted)
	// Assert
	copies, err := store.SelectAllBy("flag", true)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(copies))
}

func Test_CountAll(t *testing.T) {
	// Arrange
	store := FakeStore()
	defer store.Close()
	fakes := newFakeEntities(100)
	err := store.InsertAll(fakes...)
	assert.NoError(t, err)
	// Act
	count, err := store.CountAll()
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, int64(100), count)
}

func Test_CountBy(t *testing.T) {
	// Arrange
	store := FakeStore()
	defer store.Close()
	fakes := newFakeEntities(100)
	err := store.InsertAll(fakes...)
	assert.NoError(t, err)
	// Act
	count, err := store.CountBy("flag", true)
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, int64(50), count)
}

func Test_ExistsById(t *testing.T) {
	// Arrange
	store := FakeStore()
	defer store.Close()
	fakes := newFakeEntities(100)
	err := store.InsertAll(fakes...)
	assert.NoError(t, err)
	// Act
	exists, err := store.ExistsById(50)
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, true, exists)
}

func Test_ExistsBy(t *testing.T) {
	// Arrange
	store := FakeStore()
	defer store.Close()
	fakes := newFakeEntities(100)
	err := store.InsertAll(fakes...)
	assert.NoError(t, err)
	// Act
	exists, err := store.ExistsBy("flag", false)
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, true, exists)
}
