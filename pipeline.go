package appx

import (
	"appengine"
	"errors"
)

var ErrEmptySlice = errors.New("Cannot produce entities from an empty slice")

type Producer interface {
	BufferSize() int
	Produce(status chan<- error) <-chan Entity
}

type Consumer interface {
	SetBufferSize(size int)
	Consume(in <-chan Entity)
}

type ConsumerProducer interface {
	Consumer
	Producer
}

type Pipeline struct {
	producer  Producer
	consumers []ConsumerProducer
	status    <-chan error
	out       <-chan Entity
}

func NewPipeline(producer Producer) *Pipeline {
	return &Pipeline{
		producer:  producer,
		consumers: []ConsumerProducer{},
	}
}

func (this *Pipeline) Then(consumer ConsumerProducer) *Pipeline {
	this.consumers = append(this.consumers, consumer)
	return this
}

func (this *Pipeline) Run() *Pipeline {
	status := make(chan error, len(this.consumers)+1)
	out := this.producer.Produce(status)
	for _, consumer := range this.consumers {
		consumer.SetBufferSize(this.producer.BufferSize())
		consumer.Consume(out)
		out = consumer.Produce(status)
	}
	this.status = status
	this.out = out
	return this
}

func (this *Pipeline) CollectStatus() error {
	return NewStatusCollector(this.status).Collect()
}

func (this *Pipeline) CollectEntities() []Entity {
	return NewEntitiesCollector(this.out).Collect()
}

type EntityProducer struct {
	Stage
	entities []Entity
}

func NewEntityProducer(entities ...Entity) *EntityProducer {
	return &EntityProducer{entities: entities}
}

func (this *EntityProducer) BufferSize() int {
	return len(this.entities)
}

func (this *EntityProducer) Produce(status chan<- error) <-chan Entity {
	return NewPipelineStage(len(this.entities), func(out chan<- Entity) {
		if len(this.entities) == 0 {
			this.Emit(status, ErrEmptySlice)
			return
		}

		for _, e := range this.entities {
			out <- e
		}

		this.Emit(status, Done)
	})
}

type Stage struct {
	buffSize int
	in       <-chan Entity
}

func (this *Stage) SetBufferSize(size int) {
	this.buffSize = size
}

func (this *Stage) BufferSize() int {
	return this.buffSize
}

func (this *Stage) Consume(in <-chan Entity) {
	this.in = in
}

// Emit emits a pipeline stage result within a
// goroutine to be handled thereafter
func (this *Stage) Emit(status chan<- error, s error) {
	go func() {
		status <- s
	}()
}

type KeyResolverStage struct {
	Stage
	context  appengine.Context
}

func NewKeyResolverStage(c appengine.Context) ConsumerProducer {
	return &KeyResolverStage{
		context:  c,
	}
}

func (this *KeyResolverStage) Produce(status chan<- error) <-chan Entity {
	return NewPipelineStage(this.buffSize, func(out chan<- Entity) {
		for entity := range this.in {
			if err := ResolveKey(this.context, entity); err != nil {
				this.Emit(status, err)
				return
			}
			out <- entity
		}

		this.Emit(status, Done)
	})
}

type StatusCollector struct {
	statuses   <-chan error
	stageCount int
}

func NewStatusCollector(statuses <-chan error) *StatusCollector {
	return &StatusCollector{
		statuses:   statuses,
		stageCount: 1,
	}
}

func (this *StatusCollector) SetStageCount(count int) *StatusCollector {
	this.stageCount = count
	return this
}

func (this *StatusCollector) Collect() error {
	finalStatus := Done

	for i := 0; i < this.stageCount; i++ {
		if err := <-this.statuses; err != Done {
			finalStatus = err
		}
	}

	return finalStatus
}

type EntitiesCollector struct {
	entities <-chan Entity
}

func NewEntitiesCollector(entities <-chan Entity) *EntitiesCollector {
	return &EntitiesCollector{
		entities: entities,
	}
}

func (this *EntitiesCollector) Collect() []Entity {
	entities := []Entity{}

	for entity := range this.entities {
		entities = append(entities, entity)
	}

	return entities
}

func NewPipelineStage(buffSize int, handler func(out chan<- Entity)) chan Entity {
	out := make(chan Entity, buffSize)

	go func() {
		defer close(out)
		handler(out)
	}()

	return out
}

