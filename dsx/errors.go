package dsx

import "errors"

var (
	ErrMissingStringID = errors.New(`Model is missing StringId. String field tagged with db:"id" cannot be empty`)
	ErrMissingIntID    = errors.New(`Model is missing IntId. Integer field tagged with db:"id" cannot be zero`)
)
