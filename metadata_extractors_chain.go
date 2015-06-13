package ds

import (
	"reflect"
)

type TaggedModelExtractorsChain []MetadataExtractor

func NewTaggedModelExtractorsChain(metadata *KeyMetadata) TaggedModelExtractorsChain {
	return TaggedModelExtractorsChain{
		StringIdExtractor{metadata},
	}
}

func (this TaggedModelExtractorsChain) ExtractFrom(e Entity) error {
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
