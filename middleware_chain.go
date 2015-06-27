package appx

import "errors"

var Done = errors.New("Done")

type entityHandlerChain []EntityHandler
type EntityHandler func(e Entity) error

func NewEntityHandlerChain(handlers ...EntityHandler) entityHandlerChain {
	return entityHandlerChain(handlers)
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
