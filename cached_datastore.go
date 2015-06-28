package appx

import (
	"appengine"
	"appengine/datastore"
	"appengine/memcache"
	"time"
)

type CachedDatastore struct {
	context   appengine.Context
	rateLimit time.Duration
}

func NewCachedDatastore(c appengine.Context) *CachedDatastore {
	return &CachedDatastore{
		context: c,
	}
}

func NewCacheItem(cacheable Cacheable) *memcache.Item {
	return &memcache.Item{
		Key:    cacheable.CacheID(),
		Object: cacheable,
	}
}

func (this *CachedDatastore) CounterUpdaterFor(entity Cacheable) *CounterUpdater {
	return NewCounterUpdater(this, entity)
}

func (this *CachedDatastore) ThrottleBy(limit time.Duration) *CachedDatastore {
	this.rateLimit = limit
	return this
}

func (this *CachedDatastore) Load(cacheable Cacheable) error {
	return NewHandlerChain().
		With(KeyResolver(this.context)).
		With(CacheLoader(this.context)).
		With(CacheMissQuerier(this)).
		With(DatastoreLoader(this.context)).
		Handle(cacheable)
}

func (this *CachedDatastore) Create(cacheable Cacheable) error {
	return NewHandlerChain().
		With(KeyAssigner(this.context)).
		With(CacheSetter(this.context)).
		With(DatastoreUpdater(this.context)).
		Handle(cacheable)
}

func (this *CachedDatastore) Update(cacheable Cacheable) error {
	return NewHandlerChain().
		With(KeyResolver(this.context)).
		With(CacheSetter(this.context)).
		With(DatastoreThrottler(this)).
		With(DatastoreUpdater(this.context)).
		Handle(cacheable)
}

func (this *CachedDatastore) Delete(cacheable Cacheable) error {
	return NewHandlerChain().
		With(KeyResolver(this.context)).
		With(CacheDeleter(this.context)).
		With(CacheMissQuerier(this)).
		With(DatastoreDeleter(this.context)).
		Handle(cacheable)
}

func (this *CachedDatastore) CreateAll(slice interface{}) error {
	return NewHandlerChain().
		With(SliceValidator()).
		With(DatastoreBatchCreator(this.context)).
		With(CacheBatchSetter(this.context)).
		Handle(slice)
}

func (this *CachedDatastore) LoadAll(slice interface{}) error {
	return NewHandlerChain().
		With(SliceValidator()).
		With(CacheBatchLoader(this)).
		With(DatastoreBatchLoader(this.context)).
		Handle(slice)
}

// Query provides a appx.CachedQueryRunner for the given datastore query
func (this *CachedDatastore) Query(q *datastore.Query) *CachedQueryRunner {
	return NewCachedQueryRunner(this.context, q)
}
