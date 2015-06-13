package dsx

import (
	"github.com/drborges/ds"
	"reflect"
)

type StringIDExtractor struct {
	Metadata *ds.KeyMetadata
}

func (this StringIDExtractor) Accept(f reflect.StructField) bool {
	return f.Tag.Get("ds") != "" && f.Type.Kind() == reflect.String
}

func (this StringIDExtractor) Extract(e ds.Entity, f reflect.StructField, v reflect.Value) error {
	value := v.String()
	if value == "" {
		return ErrMissingStringID
	}
	this.Metadata.StringID = value
	return nil
}
