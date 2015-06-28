package appx

import (
	"errors"
)

var Done = errors.New("Done")

type EntityHandler func(e Entity) error
type SliceHandler func(slice interface{}) error

type sliceHandlerChain []SliceHandler
type entityHandlerChain []EntityHandler

func NewEntityHandlerChain(handlers ...EntityHandler) entityHandlerChain {
	return entityHandlerChain(handlers)
}

func (this entityHandlerChain) With(h EntityHandler) entityHandlerChain {
	this = append(this, h)
	return this
}

func (this entityHandlerChain) Handle(e Entity) error {
	for _, handler := range this {
		err := handler(e)
		if err == Done {
			return nil
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func NewSliceHandlerChain(handlers ...SliceHandler) sliceHandlerChain {
	return sliceHandlerChain(handlers)
}

func (this sliceHandlerChain) With(h SliceHandler) sliceHandlerChain {
	this = append(this, h)
	return this
}

func (this sliceHandlerChain) Handle(slice interface{}) error {
	for _, handler := range this {
		err := handler(slice)
		if err == Done {
			return nil
		}

		if err != nil {
			return err
		}
	}

	return nil
}
