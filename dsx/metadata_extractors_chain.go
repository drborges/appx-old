package dsx

import (
	"github.com/drborges/ds"
	"reflect"
)

type MetadataExtractor interface {
	Accept(reflect.StructField) bool
	Extract(ds.Entity, reflect.StructField, reflect.Value) error
}

type TaggedModelExtractorsChain []MetadataExtractor

func NewTaggedModelExtractorsChain(metadata *ds.KeyMetadata) TaggedModelExtractorsChain {
	return TaggedModelExtractorsChain{
		KindExtractor{metadata},
		IntIDExtractor{metadata},
		StringIDExtractor{metadata},
	}
}

func (this TaggedModelExtractorsChain) ExtractFrom(e ds.Entity) error {
	elemType := reflect.TypeOf(e).Elem()
	elemValue := reflect.ValueOf(e).Elem()

	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		value := elemValue.Field(i)
		for _, extractor := range this {
			if extractor.Accept(field) {
				if err := extractor.Extract(e, field, value); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
