package dsx_test

import (
	"github.com/drborges/ds"
	"github.com/drborges/ds/dsx"
	. "github.com/smartystreets/goconvey/convey"
	"reflect"
	"testing"
)

func TestKindExtractor(t *testing.T) {
	t.Parallel()

	Convey("KindExtractor", t, func() {

		Convey(`Given I have a model tagged with db:"KindName"`, func() {
			type CreditCard struct {
				ds.Model `ds:"CreditCards"`
				Owner    string
				Number   int
			}

			card := &CreditCard{Owner: "Borges"}

			Convey("When I apply the extractor to the tagged field", func() {
				field := reflect.TypeOf(card).Elem().Field(0)
				value := reflect.ValueOf(card).Elem().Field(0)

				meta := &ds.KeyMetadata{}
				err := dsx.KindExtractor{meta}.Extract(card, field, value)

				Convey("Then it succeeds", func() {
					So(err, ShouldBeNil)
				})

				Convey("Then metadata Kind field is extracted from the tagged field", func() {
					So(meta.Kind, ShouldEqual, "CreditCards")
				})
			})
		})

		Convey(`Given I have a model without kind tag`, func() {
			type CreditCard struct {
				ds.Model
				Owner  string
				Number int
			}

			card := &CreditCard{}

			Convey(`When I apply the extractor to the model field`, func() {
				field := reflect.TypeOf(card).Elem().Field(0)
				value := reflect.ValueOf(card).Elem().Field(0)

				meta := &ds.KeyMetadata{}
				err := dsx.KindExtractor{meta}.Extract(card, field, value)

				Convey("Then it succeeds", func() {
					So(err, ShouldBeNil)
				})

				Convey("Then its kind is derived from the struct's name", func() {
					So(meta.Kind, ShouldEqual, "CreditCard")
				})
			})
		})
	})
}
