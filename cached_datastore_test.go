package appx_test

import (
	"appengine/aetest"
	"appengine/datastore"
	"appengine/memcache"
	"github.com/drborges/appx"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestCachedDatastore(t *testing.T) {
	c, _ := aetest.NewContext(nil)
	defer c.Close()

	Convey("CachedDatastore", t, func() {
		Convey("Load", func() {
			Convey("Given I have a cached model", func() {
				tag := &Tag{Name: "golang", Owner: "Borges"}
				appx.ResolveKey(c, tag)
				memcache.JSON.Set(c, &memcache.Item{
					Key:    tag.CacheID(),
					Object: appx.CacheableEntity{tag, tag.EntityKey()}},
				)

				Convey("When I load it with CachedDatastore", func() {
					tagFromCache := &Tag{Name: tag.Name}
					err := appx.NewCachedDatastore(c).Load(tagFromCache)

					Convey("Then it succeeds", func() {
						So(err, ShouldBeNil)

						Convey("Then it loads the model's data", func() {
							So(tagFromCache, ShouldResemble, tag)
						})
					})
				})
			})

			Convey("Given I have a not cached model in datastore", func() {
				tag := &Tag{Name: "golang", Owner: "Borges"}
				appx.NewDatastore(c).Create(tag)

				Convey("When I load it with CachedDatastore", func() {
					tagFromCache := &Tag{Name: tag.Name}
					err := appx.NewCachedDatastore(c).Load(tagFromCache)

					Convey("Then it successfully falls back to datastore look up by key", func() {
						So(err, ShouldBeNil)

						Convey("Then it loads the model's data", func() {
							So(tagFromCache, ShouldResemble, tag)
						})
					})
				})
			})

			Convey("Given I have a queryable model saved in datastore but not cached", func() {
				account := &Account{Name: "Borges", Token: "my-auth-token"} // datastore key not resolved
				appx.NewDatastore(c).Create(account)
				time.Sleep(1 * time.Second) // gives datastore some time to index the data before querying

				Convey("When I load it missing its key", func() {
					accountFromCache := &Account{Token: account.Token}
					err := appx.NewCachedDatastore(c).Load(accountFromCache)

					Convey("Then it successfully falls back to the provided CacheMissQuery", func() {
						So(err, ShouldBeNil)

						Convey("Then it loads the model's data", func() {
							So(accountFromCache, ShouldResemble, account)
						})
					})
				})

				Convey("When I load it with its key set", func() {
					accountFromCache := &Account{}
					accountFromCache.SetEntityKey(account.EntityKey())

					err := appx.NewCachedDatastore(c).Load(accountFromCache)

					Convey("Then it successfully falls back to look up by key", func() {
						So(err, ShouldBeNil)

						Convey("Then it loads the model's data", func() {
							So(accountFromCache, ShouldResemble, account)
						})
					})
				})
			})
		})

		Convey("Create", func() {
			Convey("Given I have a not cached entity", func() {
				tag := &Tag{Name: "golang", Owner: "Borges"}

				Convey("When I create it with CachedDatastore", func() {
					err := appx.NewCachedDatastore(c).Create(tag)

					Convey("Then the operation succeeds", func() {
						So(err, ShouldBeNil)

						Convey("And I can load a cacheable entity from the cache", func() {
							cachableEntity := &appx.CacheableEntity{Cacheable: &Tag{Name: tag.Name}}
							memcache.JSON.Get(c, tag.CacheID(), cachableEntity)
							cachableEntity.Cacheable.SetEntityKey(cachableEntity.Key)

							So(cachableEntity.Cacheable, ShouldResemble, tag)
						})
					})
				})
			})

			Convey("Given I have a cached entity", func() {
				tag := &Tag{Name: "golang", Owner: "Borges"}
				appx.ResolveKey(c, tag)

				memcache.JSON.Set(c, &memcache.Item{
					Key:    tag.CacheID(),
					Object: appx.CacheableEntity{tag, tag.EntityKey()},
				})

				Convey("When I create the entity with CachedDatastore", func() {
					tag.Owner = "Diego"
					err := appx.NewCachedDatastore(c).Create(tag)

					Convey("Then the operation succeeds", func() {
						So(err, ShouldBeNil)

						Convey("And I the cache information is overwritten", func() {
							cachedEntity := &appx.CacheableEntity{Cacheable: &Tag{Name: tag.Name}}
							memcache.JSON.Get(c, tag.CacheID(), cachedEntity)
							cachedEntity.Cacheable.SetEntityKey(cachedEntity.Key)

							So(cachedEntity.Cacheable, ShouldResemble, tag)
						})
					})
				})
			})
		})

		Convey("Update", func() {
			Convey("Given I have a cached entity", func() {
				tag := &Tag{Name: "golang", Owner: "Borges"}
				cds := appx.NewCachedDatastore(c)
				cds.Create(tag)

				Convey("When I update it with CachedDatastore", func() {
					tag.Owner = "Diego"
					err := cds.Update(tag)

					Convey("Then the operation succeeds", func() {
						So(err, ShouldBeNil)

						Convey("And cache information is updated", func() {
							cachableEntity := &appx.CacheableEntity{Cacheable: &Tag{Name: tag.Name}}
							memcache.JSON.Get(c, tag.CacheID(), cachableEntity)
							cachableEntity.Cacheable.SetEntityKey(cachableEntity.Key)

							So(cachableEntity.Cacheable, ShouldResemble, tag)

							Convey("And datastore information is updated", func() {
								tagFromDatastore := &Tag{Name: tag.Name}
								appx.NewDatastore(c).Load(tagFromDatastore)

								So(tagFromDatastore, ShouldResemble, tag)
							})
						})
					})
				})
			})

			Convey("Given I have a cached queryable entity", func() {
				account := &Account{Id: 12, Name: "Borges", Token: "my-auth-token"}
				cds := appx.NewCachedDatastore(c)
				cds.Create(account)

				Convey("When I update it with CachedDatastore", func() {
					account.Name = "Diego"
					err := cds.Update(account)

					Convey("Then the operation succeeds", func() {
						So(err, ShouldBeNil)

						Convey("And cache information is updated", func() {
							cachableEntity := &appx.CacheableEntity{Cacheable: &Account{Token: account.Token}}
							memcache.JSON.Get(c, account.CacheID(), cachableEntity)
							cachableEntity.Cacheable.SetEntityKey(cachableEntity.Key)

							So(cachableEntity.Cacheable, ShouldResemble, account)

							Convey("And datastore information is updated", func() {
								accountFromDatastore := &Account{Id: 12}
								appx.NewDatastore(c).Load(accountFromDatastore)

								So(accountFromDatastore, ShouldResemble, account)
							})
						})
					})
				})
			})
		})

		Convey("Delete", func() {
			Convey("Given I have an entity cached by its cache ID", func() {
				cds := appx.NewCachedDatastore(c)
				tag := &Tag{Name: "golang", Owner: "Borges"}
				cds.Create(tag)

				Convey("When I delete the entity", func() {
					err := cds.Delete(tag)

					Convey("Then the operation succeeds", func() {
						So(err, ShouldBeNil)

						Convey("And the data is deleted from the cache", func() {
							_, err := memcache.JSON.Get(c, tag.CacheID(), nil)
							So(err, ShouldEqual, memcache.ErrCacheMiss)

							Convey("And the data is deleted from datastore", func() {
								err := appx.NewDatastore(c).Load(tag)
								So(err, ShouldEqual, datastore.ErrNoSuchEntity)
							})
						})
					})
				})
			})

			Convey("Given I have a queryable entity saved in datastore and cached", func() {
				cds := appx.NewCachedDatastore(c)
				account := &Account{Id: 321, Name: "Borges", Token: "my-auth-token"}
				cds.Create(account)
				account.Id = 0

				Convey("When I delete the entity with CachedDatastore", func() {
					err := cds.Delete(account)

					Convey("Then it successfully deletes the entity", func() {
						So(err, ShouldBeNil)

						Convey("And the data is deleted from the cache", func() {
							_, err := memcache.JSON.Get(c, account.CacheID(), nil)
							So(err, ShouldEqual, memcache.ErrCacheMiss)

							Convey("And the data is deleted from datastore", func() {
								account.Id = 321
								err := appx.NewDatastore(c).Load(account)
								So(err, ShouldEqual, datastore.ErrNoSuchEntity)
							})
						})
					})
				})
			})

			Convey("Given I have a queryable entity saved in datastore but not cached", func() {
				account := &Account{Id: 321, Name: "Borges", Token: "my-auth-token"}
				appx.NewDatastore(c).Create(account)

				Convey("When I delete the entity with its key present", func() {
					account.Token = "" // can no longer by queried on a cache miss
					err := appx.NewCachedDatastore(c).Delete(account)

					Convey("Then it successfully deletes the entity by its key", func() {
						So(err, ShouldBeNil)

						Convey("And the data is deleted from the cache", func() {
							_, err := memcache.JSON.Get(c, account.CacheID(), nil)
							So(err, ShouldEqual, memcache.ErrCacheMiss)

							Convey("And the data is deleted from datastore", func() {
								account.Id = 321
								err := appx.NewDatastore(c).Load(account)
								So(err, ShouldEqual, datastore.ErrNoSuchEntity)
							})
						})
					})
				})
			})
		})

		Convey("CreateAll", func() {
			Convey("Given I have a slice of entities", func() {
				account1 := &Account{Id: 123, Name: "Borges", Token: "borges-auth-token"}
				account2 := &Account{Id: 321, Name: "Diego", Token: "diego-auth-token"}
				accounts := []*Account{account1, account2}

				Convey("When I create them all", func() {
					err := appx.NewCachedDatastore(c).CreateAll(accounts)

					Convey("Then the operation succeeds", func() {
						So(err, ShouldBeNil)

						Convey("And the data is stored in the cache", func() {
							// Accounts are cached using their Token field as cache ID
							// See Account implementation for more detailed information
							account1FromCache := &Account{Token: account1.Token}
							account2FromCache := &Account{Token: account2.Token}
							appx.NewCachedDatastore(c).Load(account1FromCache)
							appx.NewCachedDatastore(c).Load(account2FromCache)

							So(account1FromCache, ShouldResemble, account1)
							So(account2FromCache, ShouldResemble, account2)

							Convey("And the data is stored in datastore", func() {
								account1FromDatastore := &Account{Id: account1.Id}
								account2FromDatastore := &Account{Id: account2.Id}
								appx.NewDatastore(c).Load(account1FromDatastore)
								appx.NewDatastore(c).Load(account2FromDatastore)

								So(account1FromDatastore, ShouldResemble, account1)
								So(account2FromDatastore, ShouldResemble, account2)
							})
						})
					})
				})
			})
		})

		Convey("LoadAll", func() {

			Convey("Given I have a few not cached entities in datastore", func() {
				tag1 := &Tag{Name:"golang", Owner: "Borges"}
				tag2 := &Tag{Name:"swift", Owner: "Borges"}
				tags := []*Tag{tag1, tag2}
				appx.NewDatastore(c).CreateAll(tags)

				Convey("When I load them all with CachedDatastore", func() {
					tag1FromCache := &Tag{Name: "golang"}
					tag2FromCache := &Tag{Name: "swift"}
					tagsFromCache := []*Tag{tag1FromCache, tag2FromCache}
					err := appx.NewCachedDatastore(c).LoadAll(tagsFromCache)

					Convey("Then the operation succeeds", func() {
						So(err, ShouldBeNil)

						Convey("And the data is entirely loaded from datastore", func() {
							So(tags, ShouldResemble, tagsFromCache)
						})
					})
				})
			})

			Convey("Given I have a few cached entities (not present in datastore)", func() {
				tag1 := &Tag{Name:"golang", Owner: "Borges"}
				tag2 := &Tag{Name:"swift", Owner: "Borges"}
				tags := []*Tag{tag1, tag2}
				appx.ResolveKey(c, tag1)
				appx.ResolveKey(c, tag2)

				memcache.JSON.Set(c, &memcache.Item{
					Key:    tag1.CacheID(),
					Object: appx.CacheableEntity{tag1, tag1.EntityKey()},
				})

				memcache.JSON.Set(c, &memcache.Item{
					Key:    tag2.CacheID(),
					Object: appx.CacheableEntity{tag2, tag2.EntityKey()},
				})

				Convey("And one not cached entity in datastore", func() {
					tag3 := &Tag{Name: "scala", Owner: "Borges"}
					tags = append(tags, tag3)
					appx.NewDatastore(c).Create(tag3)

					Convey("When I load them all with CachedDatastore", func() {
						tag1FromCache := &Tag{Name: "golang"}
						tag2FromCache := &Tag{Name: "swift"}
						tag3FromCache := &Tag{Name: "scala"}
						tagsFromCache := []*Tag{tag1FromCache, tag2FromCache, tag3FromCache}
						err := appx.NewCachedDatastore(c).LoadAll(tagsFromCache)

						Convey("Then the operation succeeds", func() {
							So(err, ShouldBeNil)

							Convey("And the data is loaded from the cache as well as from datastore", func() {
								So(tags, ShouldResemble, tagsFromCache)
							})
						})
					})
				})

				Convey("When I load them all with CachedDatastore", func() {
					tag1FromCache := &Tag{Name: "golang"}
					tag2FromCache := &Tag{Name: "swift"}
					tagsFromCache := []*Tag{tag1FromCache, tag2FromCache}
					err := appx.NewCachedDatastore(c).LoadAll(tagsFromCache)

					Convey("Then the operation succeeds", func() {
						So(err, ShouldBeNil)

						Convey("And the data is loaded from the cache", func() {
							So(tags, ShouldResemble, tagsFromCache)
						})
					})
				})
			})
		})
	})
}
