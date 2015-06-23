package appx_test

import (
	"appengine/aetest"
	"appengine/datastore"
	"github.com/drborges/appx"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestKeyMetadata(t *testing.T) {
	c, _ := aetest.NewContext(nil)
	defer c.Close()

	Convey("Key metadata", t, func() {
		Convey("NewKey", func() {
			Convey("returns key with parent key properly set", func() {
				parentKey := datastore.NewKey(c, "Account", "", 0, nil)
				post := &Post{Description: "Super cool!"}
				post.SetEntityParentKey(parentKey)

				key, err := appx.NewKey(c, post)

				So(err, ShouldBeNil)
				So(key.Parent(), ShouldResemble, parentKey)
			})

			Convey("returns ErrMissingParentKey", func() {
				parentKey := datastore.NewKey(c, "Account", "", 0, nil)
				comment := &Comment{Content: "Super cool!"}
				comment.SetEntityParentKey(parentKey)

				key, err := appx.NewKey(c, comment)

				So(err, ShouldBeNil)
				So(key.Parent(), ShouldResemble, parentKey)
			})
		})

		Convey("ResolveKey", func() {
			Convey("succesfully sets model key", func() {
				parentKey := datastore.NewKey(c, "Account", "", 0, nil)
				tag := &Tag{Name: "golang"}
				tag.SetEntityParentKey(parentKey)

				err := appx.ResolveKey(c, tag)

				So(err, ShouldBeNil)
				So(tag.EntityKey(), ShouldNotBeNil)
				So(tag.EntityKey().StringID(), ShouldEqual, tag.Name)
				So(tag.EntityKey().Parent(), ShouldResemble, parentKey)
			})

			Convey("returns ErrUnresolvableKey if key is not set", func() {
				comment := &Comment{Content: "Super cool!"}
				comment.SetEntityParentKey(datastore.NewKey(c, "Account", "", 0, nil))

				err := appx.ResolveKey(c, comment)

				So(comment.EntityKey(), ShouldBeNil)
				So(err, ShouldEqual, appx.ErrUnresolvableKey)
			})

			Convey("returns ErrUnresolvableKey if key set is incomplete", func() {
				comment := &Comment{Content: "Super cool!"}
				comment.SetEntityKey(datastore.NewIncompleteKey(c, "Posts", nil))
				comment.SetEntityParentKey(datastore.NewKey(c, "Account", "", 0, nil))

				err := appx.ResolveKey(c, comment)

				So(err, ShouldEqual, appx.ErrUnresolvableKey)
			})

			Convey("returns ErrMissingParentKey if parent key is not set", func() {
				comment := &Comment{Content: "Super cool!"}
				comment.SetEntityKey(datastore.NewIncompleteKey(c, "Posts", nil))

				err := appx.ResolveKey(c, comment)

				So(err, ShouldEqual, appx.ErrMissingParentKey)
			})
		})
	})
}
