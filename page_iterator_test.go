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
		Convey("Given I have a pages iterator with 2 pages each with 2 items", func() {
			q := datastore.NewQuery(Tag{}.KeyMetadata().Kind).Limit(2)
			iter := ds.Datastore{c}.Query(q).PagesIterator()

			Convey("Then I can load the first page", func() {
				firstPage := []*Tag{}
				So(iter.Cursor(), ShouldBeEmpty)
				So(iter.LoadNext(&firstPage), ShouldEqual, datastore.Done)
				So(iter.HasNext(), ShouldBeTrue)
				So(firstPage, ShouldResemble, tags[0:2])

				Convey("Then I can load the second page", func() {
					secondPage := []*Tag{}
					So(iter.Cursor(), ShouldNotBeEmpty)
					So(iter.LoadNext(&secondPage), ShouldEqual, datastore.Done)
					So(iter.HasNext(), ShouldBeTrue)
					So(secondPage, ShouldResemble, tags[2:])

					Convey("Then I cannot load more pages", func() {
						page := []*Tag{}
						So(iter.Cursor(), ShouldNotBeEmpty)
						So(iter.LoadNext(&page), ShouldEqual, datastore.Done)
						So(iter.HasNext(), ShouldBeFalse)
						So(page, ShouldBeEmpty)
					})
				})
			})
		})

		Convey("Given I have a pages iterator with zero items", func() {
			q := datastore.NewQuery(Tag{}.KeyMetadata().Kind).Filter("Owner=", "non existent").Limit(1)
			iter := ds.Datastore{c}.Query(q).PagesIterator()

			Convey("When I load the next page", func() {
				firstPage := []*Tag{}
				So(iter.Cursor(), ShouldBeEmpty)
				So(iter.LoadNext(&firstPage), ShouldEqual, datastore.Done)
				So(iter.Cursor(), ShouldBeEmpty)

				Convey("Then the page is empty", func() {
					So(firstPage, ShouldBeEmpty)

					Convey("Then it has no more results", func () {
						So(iter.HasNext(), ShouldBeFalse)
					})
				})
			})
		})
	})
}
