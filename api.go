package appx

import "appengine/datastore"

type Entity interface {
	EntityKey() *datastore.Key
	SetEntityKey(*datastore.Key)
	EncodedKey() string
	SetEncodedKey(string) error
	EntityParentKey() *datastore.Key
	SetEntityParentKey(*datastore.Key)
	KeyMetadata() *KeyMetadata
}

type Cacheable interface {
	Entity
	CacheID() string
}

type CacheMissQueryable interface {
	CacheMissQuery() *datastore.Query
}

type Creator interface {
	Create(Entity) error
}

type Updater interface {
	Update(Entity) error
}

type Loader interface {
	Load(Entity) error
}

type Deleter interface {
	Delete(Entity) error
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
