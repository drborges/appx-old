package appx

import (
	"errors"
	"appengine"
)

var ErrEmptySlice = errors.New("Cannot produce entities from an empty slice")

// Produce emits entities in the output channel until
// there is no more closing the channel at the end
type Producer func(status chan<- error) <-chan Entity

// Consume consumes incoming entities from the in channel
// processes the entity emitting the result into the output
// channel which gets closed when there is no more incoming entities
type Consumer func(in <-chan Entity, status chan<- error) <-chan Entity

// Collect collects the pipeline result returning an error in case
// any consumer emits an error message after closing its output channel
type Collect func(in <-chan Entity, status <-chan error) error

type Consume func (in <- chan Entity, out chan <- Entity, status chan <- error)

func NewProducer(entities ...Entity) Producer {
	return func(status chan<- error) <-chan Entity {
		out := make(chan Entity, len(entities))

		go func() {
			defer close(out)
			if len(entities) == 0 {
				Emit(status, ErrEmptySlice)
				return
			}

			for _, e := range entities {
				out <- e
			}

			Emit(status, Done)
		}()

		return out
	}
}

func NewKeyResolver(context appengine.Context, bufsize int) Consumer {
	return NewConsumer(bufsize, func (in <- chan Entity, out chan <- Entity, status chan <- error) {
		for entity := range in {
			if err := ResolveKey(context, entity); err != nil {
				Emit(status, err)
				return
			}
			out <- entity
		}

		Emit(status, Done)
	})
}

func NewConsumer(bufsize int, consume Consume) Consumer {
	return func(in <-chan Entity, status chan<- error) <-chan Entity {
		out := make(chan Entity, bufsize)

		go func() {
			defer close(out)
			consume(in, out, status)
		}()

		return out
	}
}

// Emit emits a pipeline stage result within a
// goroutine to be handled thereafter
func Emit(status chan <- error, s error) {
	go func () {
		status <- s
	}()
}

