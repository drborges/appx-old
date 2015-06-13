package dsx_test

import (
	. "github.com/smartystreets/goconvey/convey"
	"reflect"
	"testing"
	"github.com/drborges/ds/dsx"
	"github.com/drborges/ds"
)

func TestStringIDExtractor(t *testing.T) {
	t.Parallel()

	type CreditCard struct {
		ds.Model
		Owner  string `ds:"id"`
		Number int    `ds:"id"`
	}

	Convey("StringIDExtractor", t, func() {
		Convey("Given I have a model tagged with a string ID", func() {
			card := &CreditCard{Owner: "Borges"}

			Convey("When I apply StringIdExtractor to the tagged field", func() {
				field := reflect.TypeOf(card).Elem().Field(1)
				value := reflect.ValueOf(card).Elem().Field(1)

				meta := &ds.KeyMetadata{}
				err := dsx.StringIDExtractor{meta}.Extract(card, field, value)

				Convey("Then it succeeds", func() {
					So(err, ShouldBeNil)
				})

				Convey("Then metadata StringID field is extracted from the tagged field", func() {
					So(meta.StringID, ShouldEqual, card.Owner)
				})
			})
		})

		Convey(`Given I have a model missing its string ID`, func() {
			card := &CreditCard{}

			Convey(`When I apply StringIdExtractor to the tagged field`, func() {
				field := reflect.TypeOf(card).Elem().Field(1)
				value := reflect.ValueOf(card).Elem().Field(1)

				meta := &ds.KeyMetadata{}
				err := dsx.StringIDExtractor{meta}.Extract(card, field, value)

				Convey("Then it fails with ErrMissingStringID", func() {
					So(err, ShouldEqual, dsx.ErrMissingStringID)
				})
			})
		})

		Convey(`Given I have a model with string and int fields tagged as ds:"id"`, func() {
			card := &CreditCard{}
			stringField := reflect.TypeOf(card).Elem().Field(1)
			intField := reflect.TypeOf(card).Elem().Field(2)

			Convey(`When I apply StringIdExtractor.Accept to the fields`, func() {

				meta := &ds.KeyMetadata{}
				acceptsString := dsx.StringIDExtractor{meta}.Accept(stringField)
				acceptsInt := dsx.StringIDExtractor{meta}.Accept(intField)

				Convey("Then it accepts string field", func() {
					So(acceptsString, ShouldBeTrue)
				})

				Convey("Then it does not accept int field", func() {
					So(acceptsInt, ShouldBeFalse)
				})
			})
		})
	})
}

