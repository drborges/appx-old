package ds_test

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"appengine/aetest"
	"github.com/drborges/ds"
	"appengine/datastore"
)

func TestItemsIterator(t *testing.T) {
	c, _ := aetest.NewContext(nil)
	defer c.Close()

	tags := []*Tag{
		&Tag{Name: "golang", Owner: "Borges"},
		&Tag{Name: "ruby", Owner: "Borges"},
		&Tag{Name: "scala", Owner: "Borges"},
		&Tag{Name: "swift", Owner: "Diego"},
	}

	createAll(c, tags...)

	Convey("ItemsIterator", t, func() {
		Convey("Given I have an items iterator with 3 pages each with 1 item", func() {
			q := datastore.NewQuery(Tag{}.KeyMetadata().Kind).Filter("Owner=", "Borges").Limit(1)
			iter := ds.Datastore{c}.Query(q).ItemsIterator()
			tagsFromIterator := []*Tag{&Tag{}, &Tag{}, &Tag{}}

			Convey("Then...", func () {
				Convey("I can load the first item", func() {
					So(iter.HasNext(), ShouldBeTrue)
					So(iter.LoadNext(tagsFromIterator[0]), ShouldBeNil)
					So(iter.HasNext(), ShouldBeTrue)
					So(tagsFromIterator[0], ShouldResemble, tags[0])

					Convey("I can load the second item", func() {
						So(iter.LoadNext(tagsFromIterator[1]), ShouldBeNil)
						So(iter.HasNext(), ShouldBeTrue)
						So(tagsFromIterator[1], ShouldResemble, tags[1])

						Convey("I can load the third item", func () {
							So(iter.LoadNext(tagsFromIterator[2]), ShouldBeNil)
							So(iter.HasNext(), ShouldBeTrue)
							So(tagsFromIterator[2], ShouldResemble, tags[2])

							Convey("I cannot load more itemss", func() {
								So(iter.LoadNext(&Tag{}), ShouldEqual, datastore.Done)
								So(iter.HasNext(), ShouldBeFalse)
							})
						})
					})

					Convey("I can create a new iterator using the cursor from the previous one", func () {
						iterWithCursor := ds.Datastore{c}.Query(q).StartFrom(iter.Cursor()).ItemsIterator()

						Convey("I can load the second item", func() {
							So(iterWithCursor.LoadNext(tagsFromIterator[1]), ShouldBeNil)
							So(iterWithCursor.HasNext(), ShouldBeTrue)
							So(tagsFromIterator[1], ShouldResemble, tags[1])
						})
					})
				})
			})
		})
	})
}
