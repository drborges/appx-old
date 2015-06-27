package appx_test

import "github.com/drborges/appx"

type Comment struct {
	appx.Model
	Content string
}

// KeyMetadata in conjunction with appx.Model implements
// appx.Entity interface making Comment compatible
// with appx.Datastore
//
// A Comment is saved under the kind "Comments" as defined
// below and its key is defined to use a parent key
//
// If no parent key is provided, operations with appx.Datastore
// and its extensions like appx.CachedDatastore will fail with
// an error when resolving the model's key
func (this Comment) KeyMetadata() *appx.KeyMetadata {
	return &appx.KeyMetadata{
		Kind:      "Comments",
		HasParent: true,
	}
}
