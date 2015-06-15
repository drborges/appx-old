package ds_test

import (
	"appengine/aetest"
	"appengine/datastore"
	"github.com/drborges/ds"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type Comment struct {
	ds.Model
	Content string
}

func (this Comment) KeyMetadata() *ds.KeyMetadata {
	return &ds.KeyMetadata{
		Kind: "Comments",
		HasParent: true,
	}
}

func TestKeyMetadata(t *testing.T) {
	t.Parallel()
	c, _ := aetest.NewContext(nil)
	defer c.Close()

	Convey("Key metadata", t, func() {
		Convey("NewKey", func() {
			Convey("returns key with parent key properly set", func() {
				parentKey := datastore.NewKey(c, "Account", "", 0, nil)
				post := &Post{Description: "Super cool!"}
				post.SetParentKey(parentKey)

				key, err := ds.NewKey(c, post)

				So(err, ShouldBeNil)
				So(key.Parent(), ShouldResemble, parentKey)
			})

			Convey("returns ErrMissingParentKey", func() {
				parentKey := datastore.NewKey(c, "Account", "", 0, nil)
				comment := &Comment{Content: "Super cool!"}
				comment.SetParentKey(parentKey)

				key, err := ds.NewKey(c, comment)

				So(err, ShouldBeNil)
				So(key.Parent(), ShouldResemble, parentKey)
			})
		})

		Convey("ResolveKey", func() {
			Convey("succesfully sets model key", func() {
				parentKey := datastore.NewKey(c, "Account", "", 0, nil)
				tag := &Tag{Name: "golang"}
				tag.SetParentKey(parentKey)

				err := ds.ResolveKey(c, tag)

				So(err, ShouldBeNil)
				So(tag.Key(), ShouldNotBeNil)
				So(tag.Key().StringID(), ShouldEqual, tag.Name)
				So(tag.Key().Parent(), ShouldResemble, parentKey)
			})

			Convey("returns ErrUnresolvableKey if key is not set", func() {
				comment := &Comment{Content: "Super cool!"}
				comment.SetParentKey(datastore.NewKey(c, "Account", "", 0, nil))

				err := ds.ResolveKey(c, comment)

				So(comment.Key(), ShouldBeNil)
				So(err, ShouldEqual, ds.ErrUnresolvableKey)
			})

			Convey("returns ErrUnresolvableKey if key set is incomplete", func() {
				comment := &Comment{Content: "Super cool!"}
				comment.SetKey(datastore.NewIncompleteKey(c, "Posts", nil))
				comment.SetParentKey(datastore.NewKey(c, "Account", "", 0, nil))

				err := ds.ResolveKey(c, comment)

				So(err, ShouldEqual, ds.ErrUnresolvableKey)
			})

			Convey("returns ErrMissingParentKey if parent key is not set", func() {
				comment := &Comment{Content: "Super cool!"}
				comment.SetKey(datastore.NewIncompleteKey(c, "Posts", nil))

				err := ds.ResolveKey(c, comment)

				So(err, ShouldEqual, ds.ErrMissingParentKey)
			})
		})
	})
}
