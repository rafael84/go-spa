package database

import (
	"reflect"
	"strings"
)

const tagName = "db"

type dbTag struct {
	fieldname  string
	omitfield  bool
	autofilled bool
	pk         bool
}

func extractTag(field reflect.StructField) *dbTag {
	values := strings.Split(field.Tag.Get(tagName), ",")
	tag := dbTag{omitfield: true}
	for _, value := range values {
		if value == "autofilled" {
			tag.autofilled = true
		} else if value == "pk" {
			tag.pk = true
		} else if value != "-" {
			tag.fieldname = value
			tag.omitfield = false
		}
	}
	return &tag
}
