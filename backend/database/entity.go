package database

import (
	"fmt"
	"reflect"
)

type Entity interface {
	Table() string
}

func NewEntity(e Entity) Entity {
	mustBeStructPtr(e)
	entityType := reflect.TypeOf(e).Elem()
	newEntity := reflect.New(entityType)
	return newEntity.Interface().(Entity)
}

func mustBeStructPtr(e Entity) {
	v := reflect.ValueOf(e)
	t := v.Type()
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		panic(fmt.Errorf("entity must be pointer to struct; got %T", v))
	}
}
