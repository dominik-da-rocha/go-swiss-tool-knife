package databox

import (
	"database/sql"
	"errors"
	"fmt"
)

type revisionStore struct {
	db *sql.DB
}

type RevisionStore interface {
	GetRev(key string) (int64, error)
	IncrRev(key string) (int64, error)
	Close()
}

func NewRevisionStore(db *sql.DB) RevisionStore {
	s := revisionStore{
		db: db,
	}
	return &s
}

func (s *revisionStore) setInt(key string, value int64) error {
	res, err := s.db.Exec(`
	INSERT INTO "revision" ("key", "value") VALUES(?,?)
	ON CONFLICT ("key") DO 
	UPDATE SET "value"=? WHERE "key"=?`, key, value, value, key)

	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if affected != 1 {
		return errors.New(fmt.Sprintf("Expected 1 affected row got %d", affected))
	}
	return nil
}

func (s *revisionStore) GetRev(key string) (int64, error) {
	value := int64(0)
	err := s.db.QueryRow(`SELECT "value" FROM "revision" WHERE "key"=?`, key).Scan(&value)
	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return value, err
	}
	return value, nil
}

func (s *revisionStore) IncrRev(key string) (int64, error) {
	value, err := s.GetRev(key)
	if err == sql.ErrNoRows {
		value = 0
	} else if err != nil {
		return value, err
	}
	value++
	err = s.setInt(key, value)
	return value, err
}

func (s *revisionStore) Close() {
	s.db.Close()
}
