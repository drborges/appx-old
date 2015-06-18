package ds_test

import (
	"appengine/aetest"
	"github.com/drborges/ds"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"appengine/memcache"
	"time"
	"appengine/datastore"
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

		Convey("Update", func() {
			Convey("Given I have a cached entity", func() {
				tag := &Tag{Name: "golang", Owner: "Borges"}
				cds := ds.NewCachedDatastore(c)
				cds.Create(tag)

				Convey("When I update it with CachedDatastore", func() {
					tag.Owner = "Diego"
					err := cds.Update(tag)

					Convey("Then the operation succeeds", func() {
						So(err, ShouldBeNil)
					})

					Convey("And cache information is updated", func() {
						cachableEntity := &ds.CacheableEntity{Cacheable: &Tag{Name: tag.Name}}
						memcache.JSON.Get(c, tag.CacheID(), cachableEntity)
						cachableEntity.Cacheable.SetKey(cachableEntity.Key)

						So(cachableEntity.Cacheable, ShouldResemble, tag)
					})

					Convey("And datastore information is updated", func() {
						tagFromDatastore := &Tag{Name: tag.Name}
						ds.Datastore{c}.Load(tagFromDatastore)

						So(tagFromDatastore, ShouldResemble, tag)
					})
				})
			})

			Convey("Given I have a cached queryable entity", func() {
				account := &Account{Id: 12, Name: "Borges", Token: "my-auth-token"}
				cds := ds.NewCachedDatastore(c)
				cds.Create(account)

				Convey("When I update it with CachedDatastore", func() {
					account.Name = "Diego"
					err := cds.Update(account)

					Convey("Then the operation succeeds", func() {
						So(err, ShouldBeNil)
					})

					Convey("And cache information is updated", func() {
						cachableEntity := &ds.CacheableEntity{Cacheable: &Account{Token: account.Token}}
						memcache.JSON.Get(c, account.CacheID(), cachableEntity)
						cachableEntity.Cacheable.SetKey(cachableEntity.Key)

						So(cachableEntity.Cacheable, ShouldResemble, account)
					})

					Convey("And datastore information is updated", func() {
						accountFromDatastore := &Account{Id: 12}
						ds.Datastore{c}.Load(accountFromDatastore)

						So(accountFromDatastore, ShouldResemble, account)
					})
				})
			})
		})

		Convey("Delete", func() {
			Convey("Given I have an entity cached by its key", func() {
				cds := ds.NewCachedDatastore(c)
				tag := &Tag{Name: "golang", Owner: "Borges"}
				cds.Create(tag)

				Convey("When I delete the entity with CachedDatastore", func() {
					err := cds.Delete(tag)

					Convey("Then the operation succeeds", func() {
						So(err, ShouldBeNil)
					})

					Convey("And the data is deleted from the cache", func() {
						_, err := memcache.JSON.Get(c, tag.CacheID(), nil)
						So(err, ShouldEqual, memcache.ErrCacheMiss)
					})

					Convey("And the data is deleted from datastore", func() {
						err := ds.Datastore{c}.Load(tag)
						So(err, ShouldEqual, datastore.ErrNoSuchEntity)
					})
				})
			})

			Convey("Given I have a queryable entity cached", func() {
				cds := ds.NewCachedDatastore(c)
				account := &Account{Id: 321, Name: "Borges", Token: "my-auth-token"}
				cds.Create(account)
				account.Id = 0

				Convey("When I delete the entity with CachedDatastore", func() {
					err := cds.Delete(account)

					Convey("Then it successfully deletes the entity", func() {
						So(err, ShouldBeNil)
					})

					Convey("And the data is deleted from the cache", func() {
						_, err := memcache.JSON.Get(c, account.CacheID(), nil)
						So(err, ShouldEqual, memcache.ErrCacheMiss)
					})

					Convey("And the data is deleted from datastore", func() {
						account.Id = 321
						err := ds.Datastore{c}.Load(account)
						So(err, ShouldEqual, datastore.ErrNoSuchEntity)
					})
				})
			})
		})
	})
}