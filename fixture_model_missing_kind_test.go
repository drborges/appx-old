package appx_test

import "github.com/drborges/appx"

type ModelMissingKind struct {
	appx.Model
}

func (this ModelMissingKind) KeyMetadata() *appx.KeyMetadata {
	return &appx.KeyMetadata{}
}
