package appx_test

import (
	"appengine"
	"appengine/aetest"
	"appengine/datastore"
	"github.com/drborges/appx"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func createAll(c appengine.Context, tags ...*Tag) {
	keys := make([]*datastore.Key, len(tags))
	for i, tag := range tags {
		appx.ResolveKey(c, tag)
		keys[i] = tag.EntityKey()
	}
	datastore.PutMulti(c, keys, tags)
	time.Sleep(1 * time.Second)
}

func TestQueryRunner(t *testing.T) {
	c, _ := aetest.NewContext(nil)
	defer c.Close()

	golang := &Tag{Name: "golang", Owner: "Borges"}
	swift := &Tag{Name: "swift", Owner: "Borges"}
	ruby := &Tag{Name: "ruby", Owner: "Diego"}

	createAll(c, golang, swift, ruby)

	Convey("Given I have a QueryRunner", t, func() {
		q := appx.From(&Tag{}).Filter("Owner=", "Borges")
		runner := appx.NewQueryRunner(c, q)

		Convey("When I run Results", func() {
			result := []*Tag{}
			err := runner.Results(&result)

			Convey("Then it succeeds", func() {
				So(err, ShouldBeNil)

				Convey("Then it loads the matched entities into the given slice", func() {
					So(result, ShouldResemble, []*Tag{golang, swift})

					Convey("Then it sets keys back to the entities", func() {
						So(result[0].EntityKey(), ShouldNotBeNil)
						So(result[1].EntityKey(), ShouldNotBeNil)
					})
				})
			})
		})

		Convey("When I run Result", func() {
			tag := &Tag{}
			err := runner.Result(tag)

			Convey("Then it succeeds", func() {
				So(err, ShouldBeNil)

				Convey("Then it loads data into the given entity", func() {
					So(tag, ShouldResemble, golang)

					Convey("Then it sets the key back to the entity", func() {
						So(tag.EntityKey(), ShouldNotBeNil)
					})
				})
			})
		})

		Convey("When I run Count", func() {
			count, err := runner.Count()

			Convey("Then it succeeds", func() {
				So(err, ShouldBeNil)

				Convey("Then count is 2", func() {
					So(count, ShouldEqual, 2)
				})
			})
		})
	})
}
