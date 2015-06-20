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

func (this Account) KeyMetadata() *appx.KeyMetadata {
	return &appx.KeyMetadata{
		Kind:  "Accounts",
		IntID: this.Id,
	}
}

func (this Account) CacheID() string {
	return this.Token
}

func (this Account) CacheMissQuery() *datastore.Query {
	return datastore.NewQuery(this.KeyMetadata().Kind).Filter("Token=", this.Token).Limit(1)
}
