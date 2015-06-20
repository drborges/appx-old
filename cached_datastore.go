package appx

import (
	"appengine"
	"appengine/datastore"
	"appengine/memcache"
)

type CacheableEntity struct {
	Cacheable Cacheable
	Key       *datastore.Key
}

type CachedDatastore struct {
	ds Datastore
}

func NewCachedDatastore(c appengine.Context) *CachedDatastore {
	return &CachedDatastore{Datastore{c}}
}

// TODO any fallback to CacheMissQuery should be avoided in case the given entity
// already has a key set

func (this CachedDatastore) Load(entity Cacheable) error {
	queryable, isQueryable := entity.(CacheMissQueryable)

	if !isQueryable {
		if err := ResolveKey(this.ds.Context, entity); err != nil {
			return err
		}
	}

	// TODO check for empty CacheID?

	// Workaround to persist the not exported entity's key in the memcache
	cacheableEntity := &CacheableEntity{entity, entity.Key()}
	_, err := memcache.JSON.Get(this.ds.Context, entity.CacheID(), cacheableEntity)

	if err == memcache.ErrCacheMiss {
		// Falls back to look up by key
		// in case the key is present
		if entity.Key() != nil {
			return this.ds.Load(entity)
		}

		// In case the given cacheable is also queryable,
		// the data is retrieved by executing the
		// cache miss query provided by the entity
		if isQueryable {
			return this.ds.Query(queryable.CacheMissQuery()).Result(entity)
		}

		// Looks up the data by datastore key
		return this.ds.Load(entity)
	}

	// Sets back the key to the cacheable
	entity.SetKey(cacheableEntity.Key)
	return err
}

func (this CachedDatastore) Create(cacheable Cacheable) error {
	if err := this.ds.Create(cacheable); err != nil {
		return err
	}

	// Saves the cacheable as an entity with the key set
	// to an exported field so it may also be saved
	return memcache.JSON.Set(this.ds.Context, &memcache.Item{
		Key:    cacheable.CacheID(),
		Object: CacheableEntity{cacheable, cacheable.Key()},
	})
}

func (this CachedDatastore) Update(cacheable Cacheable) error {
	if err := this.ds.Update(cacheable); err != nil {
		return err
	}

	// Saves the cacheable as an entity with the key set
	// to an exported field so it may also be saved
	return memcache.JSON.Set(this.ds.Context, &memcache.Item{
		Key:    cacheable.CacheID(),
		Object: CacheableEntity{cacheable, cacheable.Key()},
	})
}

func (this CachedDatastore) Delete(cacheable Cacheable) error {
	// Fetches the cacheable key using the provided cache miss query
	// so it may be deleted from datastore
	if queryable, ok := cacheable.(CacheMissQueryable); ok {
		this.ds.Query(queryable.CacheMissQuery()).Result(cacheable)
	}

	if err := this.ds.Delete(cacheable); err != nil {
		return err
	}

	return memcache.Delete(this.ds.Context, cacheable.CacheID())
}
