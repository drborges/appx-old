package ds

import "errors"

var (
	ErrUnresolvableKey  = errors.New("Cannot resolve incomplete keys. Make sure you set the model's key before using it with ds.Datastore.")
	ErrMissingParentKey = errors.New("Parent key is missing. Make sure you set the parent key of your model first.")
)
