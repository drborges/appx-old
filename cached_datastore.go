package ds

import (
	"appengine/memcache"
	"appengine"
)

type CachedDatastore struct {
	ds Datastore
}

func NewCachedDatastore(c appengine.Context) *CachedDatastore {
	return &CachedDatastore{Datastore{c}}
}

func (this CachedDatastore) Load(cacheable Cacheable) error {
	queryable, implementsQueryable := cacheable.(CacheMissQueryable)

	if !implementsQueryable {
		if err := ResolveKey(this.ds.Context, cacheable); err != nil {
			return err
		}
	}

	// TODO check for empty CacheID

	_, err := memcache.JSON.Get(this.ds.Context, cacheable.CacheID(), cacheable)
	if err == memcache.ErrCacheMiss {
		// In case the given cacheable is also a queryable,
		// the data is retrieved by executing the query
		// provided by the queryable
		if implementsQueryable {
			return this.ds.Query(queryable.CacheMissQuery()).Result(cacheable)
		}

		// Looks up the data by datastore key
		return this.ds.Load(cacheable)
	}

	return err
}

