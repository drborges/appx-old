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

	Convey("Pipeline", t, func() {
		Convey("Given I have a Producer ready to produce entities", func() {
			golang := &Tag{Name: "golang"}
			swift := &Tag{Name: "swift"}
			scala := &Tag{Name: "scala"}

			producer := appx.NewEntityProducer(golang, swift, scala)

			Convey("When I produce the entities", func() {
				emitted := []appx.Entity{}
				status := make(chan error)
				for entity := range producer.Produce(status) {
					emitted = append(emitted, entity)
				}

				Convey("Then all entities are emitted", func() {
					So(emitted, ShouldResemble, []appx.Entity{
						golang,
						swift,
						scala,
					})

					Convey("Then the producer emits the Done status", func() {
						So(<-status, ShouldEqual, appx.Done)
					})
				})
			})
		})

		Convey("Given I have a Producer with no entities to produce", func() {
			producer := appx.NewEntityProducer()

			Convey("When I attempt to produce entities", func() {
				status := make(chan error)
				for entity := range producer.Produce(status) {
					panic(fmt.Sprintf("Should not process %+v", entity))
				}

				Convey("Then an error is emitted", func() {
					So(<-status, ShouldEqual, appx.ErrEmptySlice)
				})
			})
		})

		Convey("Given I have a ConsumerProducer ready to consume entities", func() {
			golang := &Tag{Name: "golang"}
			swift := &Tag{Name: "swift"}
			scala := &Tag{Name: "scala"}

			status := make(chan error)
			entities := []appx.Entity{golang, swift, scala}
			out := appx.NewEntityProducer(entities...).Produce(status)
			consumerProducer := appx.NewKeyResolverStage(c)
			consumerProducer.SetBufferSize(len(entities))
			consumerProducer.Consume(out)

			Convey("When I call produce on the ConsumerProducer", func() {
				out = consumerProducer.Produce(status)

				Convey("And I collect the pipeline resulting successful status", func() {
					stageCount := 2
					statusCollector := appx.NewStatusCollector(status).SetStageCount(stageCount)
					result := statusCollector.Collect()

					Convey("Then all entities are consumed", func() {
						for entity := range out {
							So(entity.EntityKey(), ShouldNotBeNil)
						}

						Convey("Then all stages emit their status (Producer + Consumer)", func() {
							So(result, ShouldEqual, appx.Done)
						})
					})
				})

				Convey("And I successfully collect all pipeline resulting entities", func() {
					entities := appx.NewEntitiesCollector(out).Collect()

					Convey("Then all emitted entities are consumed", func() {
						So(entities, ShouldResemble, []appx.Entity{
							golang, swift, scala,
						})
					})
				})

				Convey("Then all entities are consumed", func() {
					for entity := range out {
						So(entity.EntityKey(), ShouldNotBeNil)
					}

					Convey("Then all stages emit their status (Producer + Consumer)", func() {
						So(<-status, ShouldEqual, appx.Done)
						So(<-status, ShouldEqual, appx.Done)
					})
				})
			})
		})

		Convey("Given I have a new pipeline with a prducer and a consumer", func() {
			golang := &Tag{Name: "golang"}
			swift := &Tag{Name: "swift"}
			scala := &Tag{Name: "scala"}
			entities := []appx.Entity{golang, swift, scala}

			producer := appx.NewEntityProducer(entities...)
			consumerProducer := appx.NewKeyResolverStage(c)
			consumerProducer.SetBufferSize(len(entities))

			pipeline := appx.NewPipeline(producer).Then(consumerProducer).Run()

			Convey("When I collect the pipeline resulting status", func() {
				status := pipeline.CollectStatus()

				Convey("Then the pipeline finishes successfully", func() {
					So(status, ShouldEqual, appx.Done)
				})
			})

			Convey("When I collect the pipeline resulting entities", func() {
				collectedEntities := pipeline.CollectEntities()

				Convey("Then all entities pass through the pipeline", func() {
					So(collectedEntities, ShouldResemble, entities)
				})
			})
		})
	})
}
