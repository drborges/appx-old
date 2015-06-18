package ds_test

import (
	"appengine/aetest"
	"appengine/datastore"
	"github.com/drborges/ds"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"appengine"
	"time"
)

func createAll(c appengine.Context, tags ...*Tag) {
	keys := make([]*datastore.Key, len(tags))
	for i, tag := range tags {
		ds.ResolveKey(c, tag)
		keys[i] = tag.Key()
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

	Convey("Given I have a QueryRunner setup", t, func() {
		q := datastore.NewQuery(Tag{}.KeyMetadata().Kind).Filter("Owner=", "Borges")
		runner := ds.QueryRunner{c, q}

		Convey("When I run Results", func() {
			result := []*Tag{}
			err := runner.Results(&result)

			Convey("Then it succeeds", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then it loads the matched entities into the given slice", func() {
				So(result, ShouldResemble, []*Tag{golang, swift})
			})

			Convey("Then it sets keys back to the entities", func() {
				So(result[0].Key(), ShouldNotBeNil)
				So(result[1].Key(), ShouldNotBeNil)
			})
		})

		Convey("When I run Result", func() {
			tag := &Tag{}
			err := runner.Result(tag)

			Convey("Then it succeeds", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then it loads data into the given entity", func() {
				So(tag, ShouldResemble, golang)
			})

			Convey("Then it sets the key back to the entity", func() {
				So(tag.Key(), ShouldNotBeNil)
			})
		})

		Convey("When I run Count", func() {
			count, err := runner.Count()

			Convey("Then it succeeds", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then count is 2", func() {
				So(count, ShouldEqual, 2)
			})
		})
	})
}

