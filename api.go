package ds

import "appengine/datastore"

type (
	Persistable interface {
		KeyMetadata() *KeyMetadata
		Key() *datastore.Key
		SetKey(*datastore.Key)
	}

	CacheableModel interface {
		Persistable
		CacheKey() string
	}

	Putter interface {
		Put(Persistable) error
	}

	Loader interface {
		Load(Persistable) error
	}

	Deleter interface {
		Delete(Persistable) error
	}

	Datasource interface {
		Putter
		Loader
		Deleter
	}
)
