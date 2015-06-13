package dsx

import (
	"github.com/drborges/ds"
	"reflect"
)

type IntIDExtractor struct {
	Metadata *ds.KeyMetadata
}

func (this IntIDExtractor) Accept(f reflect.StructField) bool {
	if f.Tag.Get("ds") == "" {
		return false
	}

	switch f.Type.Kind() {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		return true
	default:
		return false
	}
}

func (this IntIDExtractor) Extract(e ds.Entity, f reflect.StructField, v reflect.Value) error {
	value := v.Int()
	if value == 0 {
		return ErrMissingIntID
	}
	this.Metadata.IntID = value
	return nil
}
