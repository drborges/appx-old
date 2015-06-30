package appx_test

import (
	"appengine/aetest"
	"fmt"
	"github.com/drborges/appx"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestPipeline(t *testing.T) {
	c, _ := aetest.NewContext(nil)
	defer c.Close()

	Convey("Given I have a few entities", t, func() {
		golang := &Tag{Name: "golang"}
		swift := &Tag{Name: "swift"}
		scala := &Tag{Name: "scala"}

		Convey("When I use a Producer to emit them into a channel", func() {
			errs := make(chan error)

			emitted := []appx.Entity{}
			for entity := range appx.NewProducer(golang, swift, scala)(errs) {
				emitted = append(emitted, entity)
			}

			Convey("Then producer finishes successfully", func() {
				So(<-errs, ShouldEqual, appx.Done)

				Convey("Then all entities are emitted", func() {
					So(emitted, ShouldResemble, []appx.Entity{golang, swift, scala})
				})
			})
		})

		Convey("When I use a Producer to emit entities from a nil slice", func() {
			errs := make(chan error)

			for item := range appx.NewProducer()(errs) {
				panic(fmt.Sprintf("Should not process %+v", item))
			}

			Convey("Then an error is emitted", func() {
				So(<-errs, ShouldEqual, appx.ErrEmptySlice)
			})
		})

		Convey("When I successfully consume produced entities with a given consumer", func() {
			errs := make(chan error)
			out := appx.NewProducer(golang, swift, scala)(errs)

			Convey("Then all emitted entities are consumed", func() {
				for item := range appx.NewKeyResolver(c, 3)(out, errs) {
					So(item.EntityKey(), ShouldNotBeNil)
				}

				Convey("Then all emitted entities are consumed", func() {
					So(<-errs, ShouldEqual, appx.Done)
					So(<-errs, ShouldEqual, appx.Done)
				})
			})
		})
	})
}
