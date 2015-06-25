package appx

import "appengine/datastore"

type Entity interface {
	EntityKey() *datastore.Key
	SetEntityKey(*datastore.Key)
	EncodedKey() string
	SetEncodedKey(string) error
	EntityParentKey() *datastore.Key
	SetEntityParentKey(*datastore.Key)
}

type Persistable interface {
	Entity
	KeyMetadata() *KeyMetadata
}

type Cacheable interface {
	Persistable
	CacheID() string
}

type CacheMissQueryable interface {
	CacheMissQuery() *datastore.Query
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

type Iterator interface {
	LoadNext(interface{}) error
	HasNext() bool
	Cursor() string
}
