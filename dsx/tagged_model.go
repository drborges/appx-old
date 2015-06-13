package dsx

import "github.com/drborges/ds"

type TaggedModel struct {
	ds.Entity
}

func (this TaggedModel) KeyMetadata() *ds.KeyMetadata {
	metadata := &ds.KeyMetadata{}
	// TODO handle extraction errors (defer, recover, panic?)
	NewTaggedModelExtractorsChain(metadata).ExtractFrom(this.Entity)
	return metadata
}
