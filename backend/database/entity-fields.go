package database

import (
	"database/sql"
	"reflect"
)

type EntityFields struct {
	Entity
}

func (ef EntityFields) extract() []reflect.StructField {
	mustBeStructPtr(ef.Entity)
	entityType := reflect.TypeOf(ef.Entity).Elem()
	fields := make([]reflect.StructField, 0)
	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)
		tag := extractTag(field)
		if !tag.omitfield {
			fields = append(fields, field)
		}
	}
	return fields
}

func (ef EntityFields) All() []string {
	fields := make([]string, 0)
	for _, field := range ef.extract() {
		tag := extractTag(field)
		fields = append(fields, tag.fieldname)
	}
	return fields
}

func (ef EntityFields) autofilled(autofilled bool) []string {
	fields := make([]string, 0)
	for _, field := range ef.extract() {
		tag := extractTag(field)
		if tag.autofilled == autofilled {
			fields = append(fields, tag.fieldname)
		}
	}
	return fields
}

func (ef EntityFields) autofilledValues(autofilled bool) []interface{} {
	entityValue := reflect.ValueOf(ef.Entity).Elem()
	values := make([]interface{}, 0)
	for _, field := range ef.extract() {
		tag := extractTag(field)
		if tag.autofilled == autofilled {
			value := entityValue.FieldByName(field.Name).Addr().Interface()
			values = append(values, value)
		}
	}
	return values
}

func (ef EntityFields) pk(pk bool) []string {
	fields := make([]string, 0)
	for _, field := range ef.extract() {
		tag := extractTag(field)
		if tag.pk == pk {
			fields = append(fields, tag.fieldname)
		}
	}
	return fields
}

func (ef EntityFields) pkValues(pk bool) []interface{} {
	entityValue := reflect.ValueOf(ef.Entity).Elem()
	values := make([]interface{}, 0)
	for _, field := range ef.extract() {
		tag := extractTag(field)
		if tag.pk == pk {
			value := entityValue.FieldByName(field.Name).Addr().Interface()
			values = append(values, value)
		}
	}
	return values
}

func (ef EntityFields) Updatable() []string {
	return ef.autofilled(false)
}

func (ef EntityFields) NotUpdatable() []string {
	return ef.autofilled(true)
}

func (ef EntityFields) Values() []interface{} {
	entityValue := reflect.ValueOf(ef.Entity).Elem()
	values := make([]interface{}, 0)
	for _, field := range ef.extract() {
		value := entityValue.FieldByName(field.Name).Addr().Interface()
		values = append(values, value)
	}
	return values
}

func (ef EntityFields) UpdatableValues() []interface{} {
	return ef.autofilledValues(false)
}

func (ef EntityFields) NotUpdatableValues() []interface{} {
	return ef.autofilledValues(true)
}

func (ef EntityFields) PK() []string {
	return ef.pk(true)
}

func (ef EntityFields) PKValues() []interface{} {
	return ef.pkValues(true)
}

func (ef EntityFields) Scan(rows *sql.Rows) error {
	values := ef.Values()
	return rows.Scan(values...)
}
