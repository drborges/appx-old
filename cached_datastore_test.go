package ds_test

import (
	"appengine/aetest"
	"github.com/drborges/ds"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"appengine/memcache"
	"time"
)

func TestCachedDatastore(t *testing.T) {
	c, _ := aetest.NewContext(nil)
	defer c.Close()

	Convey("CachedDatastore", t, func() {
		Convey("Load", func() {
			Convey("Given I have a cached model", func() {
				tag := &Tag{Name: "golang", Owner: "Borges"}
				ds.ResolveKey(c, tag)
				memcache.JSON.Set(c, &memcache.Item{
					Key: tag.CacheID(),
					Object: ds.CacheableEntity{tag, tag.Key()}},
				)

				Convey("When I load it with CachedDatastore", func() {
					tagFromCache := &Tag{Name:tag.Name}
					err := ds.NewCachedDatastore(c).Load(tagFromCache)

					Convey("Then it succeeds", func() {
						So(err, ShouldBeNil)
					})

					Convey("Then it loads the model's data", func() {
						So(tagFromCache, ShouldResemble, tag)
					})
				})
			})

			Convey("Given I have a not cached model in datastore", func() {
				tag := &Tag{Name: "golang", Owner: "Borges"}
				ds.Datastore{c}.Create(tag)

				Convey("When I load it with CachedDatastore", func() {
					tagFromCache := &Tag{Name: tag.Name}
					err := ds.NewCachedDatastore(c).Load(tagFromCache)

					Convey("Then it successfully falls back to datastore look up by key", func() {
						So(err, ShouldBeNil)
					})

					Convey("Then it loads the model's data", func() {
						So(tagFromCache, ShouldResemble, tag)
					})
				})
			})

			Convey("Given I have a not cached queryable model in datastore", func() {
				account := &Account{Name: "Borges", Token: "my-auth-token"} // datastore key not resolved
				ds.Datastore{c}.Create(account)
				time.Sleep(1 * time.Second) // gives datastore some time to index the data before querying

				Convey("When I load it with CachedDatastore", func() {
					accountFromCache := &Account{Token: account.Token}
					err := ds.NewCachedDatastore(c).Load(accountFromCache)

					Convey("Then it successfully falls back to the provided CacheMissQuery", func() {
						So(err, ShouldBeNil)
					})

					Convey("Then it loads the model's data", func() {
						So(accountFromCache, ShouldResemble, account)
					})
				})
			})
		})

		Convey("Create", func() {
			Convey("Given I have a not cached entity", func() {
				tag := &Tag{Name: "golang", Owner: "Borges"}

				Convey("When I create it with CachedDatastore", func() {
					err := ds.NewCachedDatastore(c).Create(tag)

					Convey("Then the operation succeeds", func() {
						So(err, ShouldBeNil)
					})

					Convey("And I can load a cacheable entity from the cache", func() {
						cachableEntity := &ds.CacheableEntity{Cacheable: &Tag{Name: tag.Name}}
						memcache.JSON.Get(c, tag.CacheID(), cachableEntity)
						cachableEntity.Cacheable.SetKey(cachableEntity.Key)

						So(cachableEntity.Cacheable, ShouldResemble, tag)
					})
				})
			})

			Convey("Given I have a cached entity", func() {
				tag := &Tag{Name: "golang", Owner: "Borges"}
				ds.ResolveKey(c, tag)

				memcache.JSON.Set(c, &memcache.Item{
					Key: tag.CacheID(),
					Object: ds.CacheableEntity{tag, tag.Key()},
				})

				Convey("When I create the entity with CachedDatastore", func() {
					tag.Owner = "Diego"
					err := ds.NewCachedDatastore(c).Create(tag)

					Convey("Then the operation succeeds", func() {
						So(err, ShouldBeNil)
					})

					Convey("And I the cache information is overwritten", func() {
						cachedEntity := &ds.CacheableEntity{Cacheable: &Tag{Name: tag.Name}}
						memcache.JSON.Get(c, tag.CacheID(), cachedEntity)
						cachedEntity.Cacheable.SetKey(cachedEntity.Key)

						So(cachedEntity.Cacheable, ShouldResemble, tag)
					})
				})
			})
		})
	})
}
