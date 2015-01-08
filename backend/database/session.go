package database

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
)

var (
	ERecordNotFound  = errors.New("Record not found")
	EMultipleResults = errors.New("Unexpected multiple results from query")
)

type Session struct {
	DB *sql.DB
}

func NewSession() (*Session, error) {
	db, err := Connect()
	if err != nil {
		return nil, err
	}
	return &Session{db}, nil
}

func (s *Session) Create(entity Entity) error {
	ef := EntityFields{entity}

	updatableFields := ef.Updatable()
	updatableFieldsStr := strings.Join(updatableFields, ",")

	notUpdatableFields := ef.NotUpdatable()
	notUpdatableFieldsStr := strings.Join(notUpdatableFields, ",")

	placeholders := make([]string, 0)
	for index := range updatableFields {
		placeholders = append(placeholders, fmt.Sprintf("$%d", index+1))
	}
	placeholderList := strings.Join(placeholders, ",")

	query := fmt.Sprintf(
		"insert into %s (%s) values (%s) returning %s",
		entity.Table(), updatableFieldsStr, placeholderList, notUpdatableFieldsStr,
	)
	log.Debugf("SQL: %s", query)

	updatableValues := ef.UpdatableValues()
	notUpdatableValues := ef.NotUpdatableValues()

	err := s.DB.QueryRow(query, updatableValues...).Scan(notUpdatableValues...)
	switch {
	case err == sql.ErrNoRows:
		return fmt.Errorf("Could not insert row: %s", err)
	case err != nil:
		return fmt.Errorf("Failed to execute statement: %s", err)
	default:
		return nil
	}
}

func (s *Session) One(entity Entity, where string, whereParams ...interface{}) (Entity, error) {
	entities, err := s.Filter(entity, where, whereParams...)
	if err != nil {
		return nil, err
	}
	switch len(entities) {
	case 1:
		return entities[0], nil
	case 0:
		return nil, ERecordNotFound
	default:
		return nil, EMultipleResults
	}
}

func (s *Session) Filter(entity Entity, where string, whereParams ...interface{}) ([]Entity, error) {
	fields := strings.Join(EntityFields{entity}.All(), ",")

	var sql bytes.Buffer
	sql.WriteString(fmt.Sprintf("select %s from %s", fields, entity.Table()))

	if where != "" {
		sql.WriteString(" where ")
		sql.WriteString(where)
	}

	log.WithField(
		"params", fmt.Sprintf("%s", whereParams),
	).Debugf(
		"SQL: %s", sql.String(),
	)

	rows, err := s.DB.Query(sql.String(), whereParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]Entity, 0)
	for rows.Next() {
		instance := NewEntity(entity)
		err := EntityFields{instance}.Scan(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, instance)
	}
	return list, nil
}
