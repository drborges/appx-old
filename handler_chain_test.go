package appx_test

import (
	"errors"
	"github.com/drborges/appx"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestEntityHandlerChain(t *testing.T) {
	FailingHandler := func(e appx.Entity) error {
		return errors.New("Name cannot be empty")
	}

	DoneHandler := func(e appx.Entity) error {
		return appx.Done
	}

	TagNameResolver := func(e appx.Entity) error {
		tag := e.(*Tag)
		tag.Name = "Name Resolved"
		return nil
	}

	TagOwnerResolver := func(e appx.Entity) error {
		tag := e.(*Tag)
		tag.Owner = "Owner Resolved"
		return nil
	}

	Convey("EntityHandlerChain", t, func() {
		Convey("Given I have a chain with two handlers that never fails", func() {
			chain := appx.NewEntityHandlerChain(
				TagNameResolver,
				TagOwnerResolver,
			)

			Convey("When I apply the chain to an entity", func() {
				tag := &Tag{}
				err := chain.Handle(tag)

				Convey("Then the execution succeeds", func() {
					So(err, ShouldBeNil)

					Convey("Then second handler is called", func() {
						So(tag.Name, ShouldEqual, "Name Resolved")
						So(tag.Owner, ShouldEqual, "Owner Resolved")
					})
				})
			})
		})

		Convey("Given I have a chain with two handlers where the first one fails", func() {
			chain := appx.NewEntityHandlerChain(
				FailingHandler,
				TagOwnerResolver,
			)

			Convey("When I apply the chain to an entity", func() {
				tag := &Tag{}
				err := chain.Handle(tag)

				Convey("Then the execution fails", func() {
					So(err, ShouldNotBeNil)

					Convey("Then second handler is not called", func() {
						So(tag.Owner, ShouldBeEmpty)
					})
				})
			})
		})

		Convey("Given I have a chain with two handlers where the first one returns appx.Done", func() {
			chain := appx.NewEntityHandlerChain(
				DoneHandler,
				TagOwnerResolver,
			)

			Convey("When I apply the chain to an entity", func() {
				tag := &Tag{}
				err := chain.Handle(tag)

				Convey("Then the execution succeeds", func() {
					So(err, ShouldBeNil)

					Convey("Then second handler is not called", func() {
						So(tag.Owner, ShouldBeEmpty)
					})
				})
			})
		})
	})
}
