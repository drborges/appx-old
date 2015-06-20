package appx

import "appengine/datastore"

type Resource interface {
	Entity
	ResourceID() string
	SetResourceID(string) error
}

type Entity interface {
	Key() *datastore.Key
	SetKey(*datastore.Key)
	ParentKey() *datastore.Key
	SetParentKey(*datastore.Key)
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
