package appx_test

import (
	"appengine/datastore"
	"github.com/drborges/appx"
)

type Account struct {
	appx.Model
	Id    int64
	Token string
	Name  string
}

// KeyMetadata in conjunction with appx.Model implements
// appx.Persistable interface which makes an Account compatible
// with appx.Datastore
//
// An Account is saved under the kind "Accounts" as defined below
// and it is defined to use its id field in the IntID component of
// the datastore key
func (this Account) KeyMetadata() *appx.KeyMetadata {
	return &appx.KeyMetadata{
		Kind:  "Accounts",
		IntID: this.Id,
	}
}

// CacheID implements appx.Cacheable interface which makes
// an Account compatible with appx.CachedDatastore
func (this Account) CacheID() string {
	return this.Token
}

// CacheMissQuery implements appx.CacheMissQueryable interface
// allowing appx.CachedDatastore to fall back to the provided query
// when fetching data from datastore in case of a cache miss
func (this Account) CacheMissQuery() *datastore.Query {
	return datastore.NewQuery(this.KeyMetadata().Kind).Filter("Token=", this.Token).Limit(1)
}
