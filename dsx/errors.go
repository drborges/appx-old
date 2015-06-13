package dsx

import "errors"

var (
	ErrMissingStringId = errors.New(`Model is missing StringId. String field tagged with db:"id" cannot be empty`)
	ErrMissingIntId    = errors.New(`Model is missing IntId. Integer field tagged with db:"id" cannot be zero`)
)
