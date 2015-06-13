package dsx

import (
	"reflect"
	"strings"
	"github.com/drborges/ds"
)

type KindExtractor struct {
	Metadata *ds.KeyMetadata
}

func (this KindExtractor) Accept(f reflect.StructField) bool {
	return f.Type.Name() == reflect.TypeOf(ds.Model{}).Name()
}

func (this KindExtractor) Extract(e ds.Entity, f reflect.StructField, v reflect.Value) error {
	elem := reflect.TypeOf(e).Elem()
	this.Metadata.Kind = elem.Name()

	kindMetadata := f.Tag.Get("ds")
	values := strings.Split(kindMetadata, ",")
	if strings.TrimSpace(values[0]) != "" {
		this.Metadata.Kind = values[0]
	}

	return nil
}
