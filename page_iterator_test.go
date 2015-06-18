package ds_test

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"appengine/aetest"
	"github.com/drborges/ds"
	"appengine/datastore"
)

func TestPagesIterator(t *testing.T) {
	c, _ := aetest.NewContext(nil)
	defer c.Close()

	tags := []*Tag{
		&Tag{Name: "golang", Owner: "Borges"},
		&Tag{Name: "ruby", Owner: "Borges"},
		&Tag{Name: "scala", Owner: "Borges"},
		&Tag{Name: "swift", Owner: "Diego"},
	}

	createAll(c, tags...)

	Convey("PagesIterator", t, func() {
		Convey("Given I have a pages iterator with 2 pages each with 2 item", func() {
			q := datastore.NewQuery(Tag{}.KeyMetadata().Kind).Limit(2)
			iter := ds.Datastore{c}.Query(q).PagesIterator()

			Convey("Then...", func() {
				Convey("I can load the first page", func() {
					firstPage := []*Tag{}
					So(iter.LoadNext(&firstPage), ShouldBeNil)
					So(iter.HasNext(), ShouldBeTrue)
					So(firstPage, ShouldResemble, tags[0:2])

					Convey("I can load the second page", func() {
						secondPage := []*Tag{}
						So(iter.LoadNext(&secondPage), ShouldBeNil)
						So(iter.HasNext(), ShouldBeTrue)
						So(secondPage, ShouldResemble, tags[2:])

						Convey("I cannot load more pages", func() {
							page := []*Tag{}
							So(iter.LoadNext(&page), ShouldBeNil)
							So(iter.HasNext(), ShouldBeFalse)
						})
					})
				})
			})
		})
	})
}
