package databox

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/dominik-da-rocha/go-toolbox/toolbox"
)

var ErrFirstFieldMustBeAutoIncrInteger = errors.New("Expected first field in structure to be an int64")

type ScanFunc func(row any) []any
type IdSetter func(id int64, row any)
type FactoryFunc func() any

type Store interface {
	Insert(row any) error
	InsertAll(rows ...any) error
	Update(row any) error
	UpdateAll(row ...any) error
	SelectFirst(where string, args ...any) (any, error)
	SelectById(id int64) (any, error)
	SelectAll(where string, args ...any) ([]any, error)
	SelectAllBy(column string, key any) ([]any, error)
	SelectAllByPage(offset *int64, limit *int64, orderBy *string, sortDir *string, search *string) ([]any, int64, error)
	DeleteById(id int64) (bool, error)
	Delete(where string, args ...any) (int64, error)
	DeleteAllBy(column string, key any) (int64, error)
	DeleteFromBy(table string, column string, key any) (int64, error)
	CountAll() (int64, error)
	Count(where string, args ...any) (int64, error)
	CountBy(column string, key any) (int64, error)
	Exists(where string, args ...any) (bool, error)
	ExistsById(id int64) (bool, error)
	ExistsBy(column string, key any) (bool, error)
	IncrRev() (int64, error)
	GetRev() (int64, error)
	Close() error
}

type StoreAdapter struct {
	Db          *sql.DB
	Dialect     Dialect
	Table       string
	Scan        ScanFunc
	Factory     FactoryFunc
	RevisionKey string
}

type store struct {
	db          *sql.DB
	dialect     Dialect
	table       string
	columns     []string
	placeholder []string
	scan        ScanFunc
	factory     FactoryFunc
	schema      RevisionStore
	revisionKey string
}

func NewStore(config StoreAdapter) Store {
	cols := toolbox.GetStructTags(config.Factory(), "db")
	placeholder := []string{}
	columns := []string{}
	for _, col := range cols {
		placeholder = append(placeholder, "?")
		columns = append(columns, config.Dialect.Qt(col))
	}
	toolbox.MustBeTrue(config.Db != nil, "Db must not be empty")
	toolbox.MustBeTrue(config.Factory != nil, "Factory must not be empty")
	toolbox.MustBeTrue(config.RevisionKey != "", "Revision key must not be empty")
	toolbox.MustBeTrue(config.Scan != nil, "Scan must not be empty")
	toolbox.MustBeTrue(config.Table != "", "Table must not be empty")

	return &store{
		db:          config.Db,
		dialect:     config.Dialect,
		table:       config.Dialect.Qt(config.Table),
		columns:     columns,
		placeholder: placeholder,
		scan:        config.Scan,
		factory:     config.Factory,
		schema:      NewRevisionStore(config.Db),
		revisionKey: config.RevisionKey,
	}
}

func (s *store) Insert(row any) error {
	columns := strings.Join(s.columns[1:], ",")
	placeholder := strings.Join(s.placeholder[1:], ",")
	values := s.scan(row)[1:]
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s)", s.table, columns, placeholder)
	res, err := s.db.Exec(query, values...)
	if err != nil {
		return err
	}
	rowsAff, err := res.RowsAffected()
	if err != nil {
		return err
	} else if rowsAff != 1 {
		return errors.New(fmt.Sprintf("Expected 1 affected row got %d", rowsAff))
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return err
	}

	idPtr, ok := s.scan(row)[0].(*int64)
	if ok {
		*idPtr = lastId
	}

	_, err = s.IncrRev()
	return err
}

func (s *store) InsertAll(rows ...any) error {
	for _, row := range rows {
		err := s.Insert(row)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *store) Update(row any) error {
	idName := s.columns[0]
	cols := s.columns[1:]
	setters := []string{}
	for _, col := range cols {
		setters = append(setters, fmt.Sprintf("%s=?", col))
	}
	set := strings.Join(setters, ", ")
	allValues := s.scan(row)
	values := allValues[1:]
	idPtr, ok := allValues[0].(*int64)
	if !ok {
		return ErrFirstFieldMustBeAutoIncrInteger
	}
	id := *idPtr
	values = append(values, id)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s=?", s.table, set, idName)
	res, err := s.db.Exec(query, values...)

	rowsAff, err := res.RowsAffected()
	if err != nil {
		return err
	} else if rowsAff != 1 {
		return errors.New(fmt.Sprintf("Expected 1 affected row got %d", rowsAff))
	}

	_, err = s.IncrRev()
	return err
}

func (s *store) UpdateAll(rows ...any) error {
	for _, row := range rows {
		err := s.Update(row)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *store) SelectFirst(where string, args ...any) (any, error) {
	result := s.factory()
	selectCols := strings.Join(s.columns, ",")
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE %s LIMIT 1`, selectCols, s.table, where)
	values := s.scan(result)
	err := s.db.QueryRow(query, args...).Scan(values...)
	return result, err
}

func (s *store) SelectById(id int64) (any, error) {
	idName := s.columns[0]
	where := fmt.Sprintf("%s=?", idName)
	return s.SelectFirst(where, id)
}

func (s *store) SelectAll(where string, args ...any) ([]any, error) {
	result := []any{}
	selectCols := strings.Join(s.columns, ",")
	query := fmt.Sprintf(`SELECT %s FROM %s %s`, selectCols, s.table, where)
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		item := s.factory()
		err = rows.Scan(s.scan(item)...)
		if err != nil {
			return result, err
		}
		result = append(result, item)
	}

	return result, nil
}

func (s *store) SelectAllBy(column string, key any) ([]any, error) {
	idName := s.columns[0]
	column = s.dialect.Qt(column)
	where := fmt.Sprintf(`WHERE %s = ? ORDER BY %s ASC`, column, idName)
	return s.SelectAll(where, key)
}

func (s *store) SelectAllByPage(offsetPtr *int64, limitPtr *int64, orderByPtr *string, sortDirPtr *string, searchPtr *string) ([]any, int64, error) {
	items := []any{}
	total := int64(0)
	offset := toolbox.ToInt64(offsetPtr, 0)
	limit := toolbox.ToInt64(limitPtr, 100)
	orderBy := toolbox.ToAnyStringOf(s.dialect.QtPtr(orderByPtr), s.columns, s.columns[0])
	sortDir := toolbox.ToAnyStringOf(sortDirPtr, []string{"asc", "desc"}, "asc")
	sortDir = strings.ToUpper(sortDir)
	columns := strings.Join(s.columns, ",")
	idName := s.columns[0]

	countQuery := ""
	selectQuery := ""
	args := []any{}

	if searchPtr == nil {
		selectQuery = fmt.Sprintf(`SELECT %s FROM %s ORDER BY %s %s LIMIT %d OFFSET %d `, columns, s.table, orderBy, sortDir, limit, offset)
		countQuery = fmt.Sprintf(`SELECT COUNT(%s) FROM %s`, orderBy, s.table)
	} else {
		wheres := []string{}
		likeExp := "%" + *searchPtr + "%"
		for _, col := range s.columns {
			like := fmt.Sprintf(`%s LIKE ?`, col)
			wheres = append(wheres, like)
			args = append(args, likeExp)
		}
		where := strings.Join(wheres, " OR ")
		selectQuery = fmt.Sprintf(`SELECT %s FROM %s WHERE %s ORDER BY %s LIMIT %d OFFSET %d`, columns, s.table, where, orderBy, limit, offset)
		countQuery = fmt.Sprintf(`SELECT COUNT(%s) FROM %s WHERE %s`, idName, s.table, where)
	}
	rows, err := s.db.Query(selectQuery, args...)
	if err != nil {
		return items, total, err
	}
	defer rows.Close()

	for rows.Next() {
		item := s.factory()
		err = rows.Scan(s.scan(item)...)
		if err != nil {
			return items, total, err
		}
		items = append(items, item)
	}

	err = s.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return items, total, err
	}

	return items, total, err
}

func (s *store) DeleteById(id int64) (bool, error) {
	idName := s.columns[0]
	query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", s.table, idName)
	res, err := s.db.Exec(query, id)
	if err != nil {
		return false, err
	}

	rowsAff, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	if rowsAff != 1 {
		return false, nil
	}

	_, err = s.IncrRev()

	return true, err
}

func (s *store) Delete(where string, args ...any) (int64, error) {

	query := fmt.Sprintf("DELETE FROM %s WHERE %s", s.table, where)
	res, err := s.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	rowsAff, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	_, err = s.IncrRev()
	return rowsAff, err
}

func (s *store) DeleteAllBy(column string, key any) (int64, error) {
	column = s.dialect.Qt(column)
	where := fmt.Sprintf("%s=?", column)
	return s.Delete(where, key)
}

func (s *store) CountAll() (int64, error) {
	count := int64(0)
	idName := s.columns[0]
	query := fmt.Sprintf("SELECT COUNT(%s) FROM %s", idName, s.table)
	err := s.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *store) DeleteFromBy(table string, column string, key any) (int64, error) {
	column = s.dialect.Qt(column)
	query := fmt.Sprintf("DELETE FROM %s WHERE %s=?", table, column)
	slog.Debug("query", query, "values", key)
	res, err := s.db.Exec(query, key)
	if err != nil {
		return 0, err
	}
	rowsAff, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	_, err = s.IncrRev()
	return rowsAff, err
}

func (s *store) Count(where string, args ...any) (int64, error) {
	count := int64(0)
	idName := s.columns[0]
	query := fmt.Sprintf("SELECT COUNT(%s) FROM %s WHERE %s", idName, s.table, where)
	err := s.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *store) CountBy(column string, key any) (int64, error) {
	where := s.dialect.Qt(column)
	return s.Count(where, key)
}

func (s *store) Exists(where string, args ...any) (bool, error) {
	count := int64(0)
	idName := s.columns[0]
	query := fmt.Sprintf("SELECT COUNT(%s) FROM %s WHERE %s LIMIT 1", idName, s.table, where)
	err := s.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *store) ExistsById(id int64) (bool, error) {
	idName := s.columns[0]
	where := fmt.Sprintf("%s=?", idName)
	return s.Exists(where, id)
}

func (s *store) ExistsBy(column string, key any) (bool, error) {
	column = s.dialect.Qt(column)
	where := fmt.Sprintf("%s=?", column)
	return s.Exists(where, key)
}

func (s *store) Close() error {
	return s.db.Close()
}

func (s *store) IncrRev() (int64, error) {
	return s.schema.IncrRev(s.revisionKey)
}

func (s *store) GetRev() (int64, error) {
	return s.schema.GetRev(s.revisionKey)
}
