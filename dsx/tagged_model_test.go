package dsx_test

import (
	"github.com/drborges/ds"
	"github.com/drborges/ds/dsx"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTaggedModel(t *testing.T) {
	t.Parallel()

	Convey("ds.TaggedModel", t, func() {

		Convey(`Given I have a model tagged with ds:"KindName"`, func() {
			type User struct {
				ds.Model `ds:"Users"`
				Name     string
				Twitter  string
			}

			user := User{Name: "Diego", Twitter: "@drborges"}

			Convey("When I wrap it with TaggedModel", func() {
				taggedUser := dsx.TaggedModel{&user}

				Convey("Then KeyMetadata extracts kind information from the tag", func() {
					metadata := taggedUser.KeyMetadata()
					So(metadata.Kind, ShouldEqual, "Users")
				})
			})
		})

		Convey(`Given I have a model with a string field tagged with ds:"id"`, func() {
			type User struct {
				ds.Model
				Name    string `ds:"id"`
				Twitter string
			}

			user := User{Name: "Diego", Twitter: "@drborges"}

			Convey("When I wrap it with TaggedModel", func() {
				taggedUser := dsx.TaggedModel{&user}

				Convey("Then KeyMetadata extracts string id from the tagged field", func() {
					metadata := taggedUser.KeyMetadata()
					So(metadata.StringID, ShouldEqual, user.Name)
				})
			})
		})

		Convey(`Given I have a model with a int field tagged with ds:"id"`, func() {
			type Account struct {
				ds.Model
				Id   int `ds:"id"`
				Name string
			}

			account := Account{Name: "Diego", Id: 123}

			Convey("When I wrap it with TaggedModel", func() {
				taggedAccount := dsx.TaggedModel{&account}

				Convey("Then KeyMetadata extracts string id from the tagged field", func() {
					metadata := taggedAccount.KeyMetadata()
					So(metadata.IntID, ShouldEqual, account.Id)
				})
			})
		})
	})
}
