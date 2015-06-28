package appx

import (
	"errors"
)

var Done = errors.New("Done")

type Handler func(e interface{}) (interface{}, error)
type HandlerChain []Handler

func NewHandlerChain(handlers ...Handler) HandlerChain {
	return HandlerChain(handlers)
}

func (this HandlerChain) With(h Handler) HandlerChain {
	this = append(this, h)
	return this
}

func (this HandlerChain) Handle(item interface{}) error {
	for _, handler := range this {
		processedEntity, err := handler(item)
		item = processedEntity

		if err == Done {
			return nil
		}

		if err != nil {
			return err
		}
	}

	return nil
}
