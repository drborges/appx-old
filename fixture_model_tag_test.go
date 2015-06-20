package appx_test

import "github.com/drborges/appx"

type Tag struct {
	appx.Model
	Name  string
	Owner string
}

func (this Tag) KeyMetadata() *appx.KeyMetadata {
	return &appx.KeyMetadata{
		Kind:     "Tags",
		StringID: this.Name,
	}
}

func (this Tag) CacheID() string {
	return this.Name
}
