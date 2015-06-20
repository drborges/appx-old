package appx_test

import (
	"appengine/aetest"
	"appengine/datastore"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestModel(t *testing.T) {
	c, _ := aetest.NewContext(nil)
	defer c.Close()

	Convey("appx.Model", t, func() {
		Convey("Given I have a model with parent key", func() {
			parentKey := datastore.NewKey(c, "User", "Borges", 0, nil)
			commentKey := datastore.NewKey(c, "Comments", "", 0, parentKey)
			comment := Comment{}

			Convey("When I set its ID", func() {
				err := comment.SetID(commentKey.Encode())

				Convey("Then it succeeds", func() {
					So(err, ShouldBeNil)

					Convey("Then its key is set", func() {
						So(comment.Key(), ShouldResemble, commentKey)

						Convey("Then its parent key is set", func() {
							So(comment.ParentKey(), ShouldResemble, parentKey)
						})
					})
				})
			})
		})
	})
}
