package ds

type TaggedModel struct {
	Entity
}

func (this TaggedModel) KeyMetadata() *KeyMetadata {
	metadata := &KeyMetadata{}
	// TODO handle extraction errors (defer, recover, panic?)
	NewTaggedModelExtractorsChain(metadata).ExtractFrom(this.Entity)
	return metadata
}
