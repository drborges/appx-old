package appx_test

import (
	"appengine/aetest"
	"appengine/datastore"
	"appengine/memcache"
	"github.com/drborges/appx"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCounterUpdater(t *testing.T) {
	c, _ := aetest.NewContext(nil)
	defer c.Close()

	Convey("CounterUpdater", t, func() {
		Convey("Given I have a cached model with a counter field", func() {
			tag := &Tag{Name: "golang", Owner: "Borges"}
			Convey("When I Increment its counter field twice", func() {
				counterUpdater := appx.NewCachedDatastore(c).CounterUpdaterFor(tag)
				err1 := counterUpdater.Increment(&tag.PostCount)
				err2 := counterUpdater.Increment(&tag.PostCount)

				Convey("Then the operations succeed", func() {
					So(err1, ShouldBeNil)
					So(err2, ShouldBeNil)

					Convey("Then the entity is updated", func() {
						So(tag.PostCount, ShouldEqual, 2)

						Convey("Then memcache is updated", func() {
							item, _ := memcache.Get(c, "Counter,Tags,golang")
							So(string(item.Value), ShouldEqual, "2")

							Convey("Then datastore is updated", func() {
								tagFromDatastore := &Tag{Name: "golang"}
								appx.NewDatastore(c).Load(tagFromDatastore)

								So(tagFromDatastore, ShouldResemble, tag)
							})
						})
					})
				})
			})
		})

		Convey("Given I have another cached model with a counter field", func() {
			tag := &Tag{Name: "android", Owner: "Borges"}

			Convey("When I Increment its counter field twice on memcache only", func() {
				counterUpdater := appx.NewCachedDatastore(c).CounterUpdaterFor(tag)
				err1 := counterUpdater.CacheOnly().Increment(&tag.PostCount)
				err2 := counterUpdater.CacheOnly().Increment(&tag.PostCount)

				Convey("Then the operations succeed", func() {
					So(err1, ShouldBeNil)
					So(err2, ShouldBeNil)

					Convey("Then the entity is updated", func() {
						So(tag.PostCount, ShouldEqual, 2)

						Convey("Then memcache is updated", func() {
							item, _ := memcache.Get(c, "Counter,Tags,android")
							So(string(item.Value), ShouldEqual, "2")

							Convey("Then datastore is not updated", func() {
								tagFromDatastore := &Tag{Name: "android"}
								err := appx.NewDatastore(c).Load(tagFromDatastore)
								So(err, ShouldEqual, datastore.ErrNoSuchEntity)
							})
						})
					})
				})
			})

			Convey("When I Increment its counter field with a custom key", func() {
				err := appx.NewCachedDatastore(c).
					CounterUpdaterFor(tag).
					Key("post count").
					Increment(&tag.PostCount)

				Convey("Then the operations succeed", func() {
					So(err, ShouldBeNil)

					Convey("Then I can retrieve the counter from memcache using the custom key", func() {
						item, _ := memcache.Get(c, "post count")
						So(string(item.Value), ShouldEqual, "1")

					})
				})
			})

			Convey("When I Decrement an exitent counter", func() {
				appx.NewCachedDatastore(c).
					CounterUpdaterFor(tag).
					Key("a new counter").
					Increment(&tag.PostCount)

				err := appx.NewCachedDatastore(c).
					CounterUpdaterFor(tag).
					Key("a new counter").
					Decrement(&tag.PostCount)

				Convey("Then the operations succeed", func() {
					So(err, ShouldBeNil)

					Convey("Then the entity is updated", func() {
						So(tag.PostCount, ShouldEqual, 0)

						Convey("Then memcache is updated", func() {
							item, _ := memcache.Get(c, "a new counter")
							So(string(item.Value), ShouldEqual, "0")

							Convey("Then datastore is updated", func() {
								tagFromDatastore := &Tag{Name: "android"}
								appx.NewDatastore(c).Load(tagFromDatastore)

								So(tagFromDatastore, ShouldResemble, tag)
							})
						})
					})
				})
			})

			Convey("When I Add a delta to an exitent counter", func() {
				appx.NewCachedDatastore(c).
					CounterUpdaterFor(tag).
					Key("counter").
					Increment(&tag.PostCount)

				err := appx.NewCachedDatastore(c).
					CounterUpdaterFor(tag).
					Key("counter").
					Add(&tag.PostCount, 3)

				Convey("Then the operations succeed", func() {
					So(err, ShouldBeNil)

					Convey("Then the entity is updated", func() {
						So(tag.PostCount, ShouldEqual, 4)

						Convey("Then memcache is updated", func() {
							item, _ := memcache.Get(c, "counter")
							So(string(item.Value), ShouldEqual, "4")

							Convey("Then datastore is updated", func() {
								tagFromDatastore := &Tag{Name: "android"}
								appx.NewDatastore(c).Load(tagFromDatastore)

								So(tagFromDatastore, ShouldResemble, tag)
							})
						})
					})
				})
			})
		})
	})
}
