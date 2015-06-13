package dsx

import (
	"reflect"
	"github.com/drborges/ds"
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
		return ErrMissingStringId
	}
	this.Metadata.StringID = value
	return nil
}
