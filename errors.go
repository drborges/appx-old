package appx

import "errors"

var (
	ErrInvalidEntityType  = errors.New("Invalid entity type. Make sure your model implements appx.Entity (watch out for pointer receivers)")
	ErrInvalidSliceType   = errors.New("Invalid slice type. Make sure you pass a pointer to a slice of appx.Entity")
	ErrUnresolvableKey    = errors.New("Cannot resolve incomplete keys. Make sure you set the model's key first")
	ErrMissingParentKey   = errors.New("Parent key is missing. Make sure you set the parent key on your model first")
	ErrNonCacheableEntity = errors.New("Non cacheable entity. Make sure the entity implements appx.Cacheable")
)
