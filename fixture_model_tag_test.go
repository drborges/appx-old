package appx_test

import "github.com/drborges/appx"

type Tag struct {
	appx.Model
	Name      string
	Owner     string
	PostCount appx.Counter
}

// KeyMetadata in conjunction with appx.Model implement
// appx.Persistable interface making Tag compatible with
// appx.Datastore
//
// A tag key is defined to use its name as the StringID
// component in the datastore key
func (this Tag) KeyMetadata() *appx.KeyMetadata {
	return &appx.KeyMetadata{
		Kind:     "Tags",
		StringID: this.Name,
	}
}

// CacheID implements the appx.Cacheable interface
// making a Tag compatible with appx.CachedDatastore
//
// The Tag's name is used as the cache id and Fall backs
// to datastore are done by look ups using the model's key
func (this Tag) CacheID() string {
	return this.Name
}
