package appx_test

import "github.com/drborges/appx"

type Post struct {
	appx.Model
	Description string
}

// KeyMetadata in conjunction with appx.Model implement
// appx.Persistable interface making Tag compatible with
// appx.Datastore
//
// Since IntID nor StringID fields are used in the KeyMetadata
// definition the key generated for a Post is incomplete a.k.a IntID == 0, StringID == ""
//
// Incomplete keys signal to datastore that a auto generated key
// needs to be created for the given entity upon creation and the actual
// key is then returned
func (this Post) KeyMetadata() *appx.KeyMetadata {
	return &appx.KeyMetadata{
		Kind: "Posts",
	}
}
