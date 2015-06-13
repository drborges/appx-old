package ds

import "appengine/datastore"

type Entity interface {
	Key() *datastore.Key
	SetKey(*datastore.Key)
}

type Persistable interface {
	Entity
	KeyMetadata() *KeyMetadata
}

type CacheableModel interface {
	Persistable
	CacheKey() string
}

type Creator interface {
	Create(Persistable) error
}

type Updater interface {
	Update(Persistable) error
}

type Loader interface {
	Load(Persistable) error
}

type Deleter interface {
	Delete(Persistable) error
}

type Datasource interface {
	Loader
	Creator
	Updater
	Deleter
}
