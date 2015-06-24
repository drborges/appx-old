package appx_test

import (
	"appengine/aetest"
	"appengine/memcache"
	"github.com/drborges/appx"
	. "github.com/smartystreets/goconvey/convey"
	"strconv"
	"testing"
	"time"
)

func TestCachedQueryRunner(t *testing.T) {
	c, _ := aetest.NewContext(nil)
	defer c.Close()

	Convey("CachedQueryRunner", t, func() {
		Convey("Given I have a few entities in datastore", func() {
			tag1 := &Tag{Name: "golang", Owner: "Borges"}
			tag2 := &Tag{Name: "swift", Owner: "Borges"}
			tags := []*Tag{tag1, tag2}
			appx.NewDatastore(c).CreateAll(tags)
			time.Sleep(1 * time.Second)

			Convey("And a query that matches them all", func() {
				q := appx.From(&Tag{})

				Convey("When I run Count it caches the result", func() {
					runner := appx.NewCachedQueryRunner(c, q).
						CachedAs("count tags").
						ExpiresIn(5 * time.Second)

					count, err := runner.Count()

					Convey("Then the operation succeeds", func() {
						So(err, ShouldBeNil)

						Convey("Then the count matches the total number of entities", func() {
							So(count, ShouldEqual, len(tags))

							Convey("Then the count value is cached with expiration time", func() {
								item, _ := memcache.Get(c, "count tags")
								So(string(item.Value), ShouldEqual, strconv.Itoa(count))
								// Seems like appengine dev server does not set the expiration
								// time field of a cache item :~
								//So(item.Expiration, ShouldEqual, 2*time.Second)
								Convey("Then subsequent counts will hit the cache rather than datastore", func() {
									appx.NewDatastore(c).Delete(tag1)

									subsequentCount, err := runner.Count()
									So(err, ShouldBeNil)
									So(subsequentCount, ShouldEqual, count)
								})
							})
						})
					})
				})
			})
		})
	})
}
