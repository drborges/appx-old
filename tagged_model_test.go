package ds_test

import (
	"github.com/drborges/ds"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTaggedModel(t *testing.T) {
	Convey("ds.TaggedModel", t, func() {
		Convey(`Given I have a model with a string field tagged with ds:"id"`, func() {
			type User struct {
				ds.Model
				Name    string `ds:"id"`
				Twitter string
			}

			user := User{Name: "Diego", Twitter: "@drborges"}

			Convey("When I wrap it with TaggedModel", func() {
				taggedUser := ds.TaggedModel{&user}

				Convey("Then KeyMetadata extracts string id from the tagged field", func() {
					metadata := taggedUser.KeyMetadata()
					So(metadata.StringID, ShouldEqual, user.Name)
				})
			})
		})
	})
}
