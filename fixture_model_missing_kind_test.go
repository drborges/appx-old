package appx_test

import "github.com/drborges/appx"

type ModelMissingKind struct {
	appx.Model
}

// KeyMetadata in conjunction with appx.Model implement
// appx.Persistable interface making Tag compatible with
// appx.Datastore
//
// This particular KeyMetadata definition is invalid since
// Kind field is not specified. As of now, this field is
// required.
func (this ModelMissingKind) KeyMetadata() *appx.KeyMetadata {
	return &appx.KeyMetadata{} // missing required Kind field
}
