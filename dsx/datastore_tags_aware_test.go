package dsx_test

import (
	"appengine/aetest"
	"appengine/datastore"
	"github.com/drborges/ds"
	"github.com/drborges/ds/dsx"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestLoadModel(t *testing.T) {
	c, _ := aetest.NewContext(nil)
	defer c.Close()

	Convey("dsx.DatastoreTagsAware", t, func() {
		type Account struct {
			ds.Model
			Id   int64
			Name string `ds:"id"`
		}

		Convey("Load", func() {
			Convey("Given I have a model with StringID key saved in datastore", func() {
				account := &Account{Name: "Borges", Id: 123}
				key := ds.NewKey(c, dsx.TaggedModel{account})
				datastore.Put(c, key, account)
				account.SetKey(key)

				Convey("When I load it using dsx.DatastoreTagsAware", func() {
					loadedAccount := Account{Name: "Borges"}
					err := dsx.DatastoreTagsAware{c}.Load(&loadedAccount)

					Convey("Then it succeeds", func() {
						So(err, ShouldBeNil)
					})

					Convey("Then the data is loaded from datastore", func() {
						So(&loadedAccount, ShouldResemble, account)
					})

					Convey("Then the model has its key resolved", func() {
						So(loadedAccount.Key(), ShouldNotBeNil)
					})
				})
			})
		})
	})
}
