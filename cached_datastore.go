package appx

import (
	"appengine"
	"appengine/datastore"
	"appengine/memcache"
	"encoding/json"
	"reflect"
	"time"
)

type CachedDatastore struct {
	ds        *Datastore
	rateLimit time.Duration
}

func NewCachedDatastore(c appengine.Context) *CachedDatastore {
	return &CachedDatastore{
		ds: NewDatastore(c),
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

func (this *CachedDatastore) Load(entity Cacheable) error {
	queryable, isQueryable := entity.(CacheMissQueryable)

	if !isQueryable {
		if err := ResolveKey(this.ds.context, entity); err != nil {
			return err
		}
	}

	needQueryFallback := entity.EntityKey() == nil && isQueryable
	_, err := memcache.JSON.Get(this.ds.context, entity.CacheID(), entity)

	if err == memcache.ErrCacheMiss {
		// Falls back to look up by key
		// in case the key is present
		if !needQueryFallback {
			return this.ds.Load(entity)
		}

		// In case the given cacheable is also queryable,
		// the data is retrieved by executing the
		// cache miss query provided by the entity
		if needQueryFallback {
			return this.ds.Query(queryable.CacheMissQuery()).Result(entity)
		}
	}
	return err
}

func (this *CachedDatastore) Create(cacheable Cacheable) error {
	if err := this.ds.Create(cacheable); err != nil {
		return err
	}

	// Saves the cacheable as an entity with the key set
	// to an exported field so it may also be saved
	return memcache.JSON.Set(this.ds.context, NewCacheItem(cacheable))
}

func (this *CachedDatastore) CreateAll(slice interface{}) error {
	// Creates all entities in datastore first so in case of
	// error no data would be cached
	if err := this.ds.CreateAll(slice); err != nil {
		return err
	}

	// At this point we are safe to assume
	// slice is actually a slice since otherwise
	// this.ds.CreateAll would return an error
	s := reflect.ValueOf(slice)

	// Create a memcache.Item for each cacheable
	// in the given slice
	items := make([]*memcache.Item, s.Len())
	for i := 0; i < s.Len(); i++ {
		cacheable := s.Index(i).Interface().(Cacheable)
		items[i] = NewCacheItem(cacheable)
	}

	// Saves the cacheable as an entity with the key set
	// to an exported field so it may also be saved
	return memcache.JSON.SetMulti(this.ds.context, items)
}

func (this *CachedDatastore) LoadAll(slice interface{}) error {
	s := reflect.ValueOf(slice)

	if s.Kind() != reflect.Slice {
		return datastore.ErrInvalidEntityType
	}

	// Creates a slice of cache keys and
	// makes sure every entity in the given
	// slice is appx.Cacheable
	cacheKeys := make([]string, s.Len())
	for i, _ := range cacheKeys {
		cacheable, ok := s.Index(i).Interface().(Cacheable)
		if !ok {
			return ErrNonCacheableEntity
		}
		cacheKeys[i] = cacheable.CacheID()
	}

	items, err := memcache.GetMulti(this.ds.context, cacheKeys)
	if err != nil {
		return err
	}

	// Falls back to datastore batch load
	// in case no entity is cached
	if len(items) == 0 {
		return this.ds.LoadAll(slice)
	}

	for i := 0; i < s.Len(); i++ {
		entity := s.Index(i).Interface().(Cacheable)
		// At this point we are safe to assume the entity
		// implements appx.Cacheable since it was already
		// verified while creating the slice of cache keys
		cacheID := entity.CacheID()
		item, itemCached := items[cacheID]

		if !itemCached {
			// Falls back to datastore
			if err := this.Load(entity); err == datastore.ErrNoSuchEntity {
				// TODO handle partial loads = not in cache nor datastore
				// return error with the failed entity's index? so the user
				// may decide whether to remove the entity from the list or not?
				return err
			}
		} else {
			if err := json.Unmarshal(item.Value, &entity); err != nil {
				return err
			}
		}
	}

	return nil
}

func (this *CachedDatastore) Update(cacheable Cacheable) error {
	if err := memcache.JSON.Set(this.ds.context, NewCacheItem(cacheable)); err != nil {
		return err
	}

	throttledCacheID := "throttled" + cacheable.CacheID()
	_, err := memcache.Get(this.ds.context, throttledCacheID)
	throttled := err != memcache.ErrCacheMiss

	if !throttled {
		if err := this.ds.Update(cacheable); err != nil {
			return err
		}
	}

	if !throttled && this.rateLimit > 0 {
		memcache.Set(this.ds.context, &memcache.Item{
			Key:        throttledCacheID,
			Expiration: this.rateLimit,
			Value:      []byte(throttledCacheID),
		})
	}

	return nil
}

func (this *CachedDatastore) Delete(cacheable Cacheable) error {
	// Fetches the cacheable key using the provided cache miss query
	// so it may be deleted from datastore
	if queryable, ok := cacheable.(CacheMissQueryable); ok && cacheable.EntityKey() == nil {
		if err := this.ds.Query(queryable.CacheMissQuery()).Result(cacheable); err != nil {
			return err
		}
	}

	if err := this.ds.Delete(cacheable); err != nil {
		return err
	}

	// don't really care about cache misses errors
	if err := memcache.Delete(this.ds.context, cacheable.CacheID()); err != nil && err != memcache.ErrCacheMiss {
		return err
	}

	return nil
}

// Query provides a appx.CachedQueryRunner for the given datastore query
func (this *CachedDatastore) Query(q *datastore.Query) *CachedQueryRunner {
	return NewCachedQueryRunner(this.ds.context, q)
}
