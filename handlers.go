package appx

import (
	"appengine"
	"appengine/datastore"
	"reflect"
	"appengine/memcache"
	"encoding/json"
)

func KeyResolver(c appengine.Context) Handler {
	return func(item interface{}) (interface{}, error) {
		e := item.(Entity)
		_, isQueryable := e.(CacheMissQueryable)

		err := ResolveKey(c, e)

		if err == ErrUnresolvableKey && isQueryable {
			return e, nil
		}

		return e, err
	}
}

func KeyAssigner(c appengine.Context) Handler {
	return func(item interface{}) (interface{}, error) {
		e := item.(Entity)
		key, err := NewKey(c, e)
		if err != nil {
			return e, err
		}
		e.SetEntityKey(key)
		return e, nil
	}
}

func DatastoreLoader(c appengine.Context) Handler {
	return func(item interface{}) (interface{}, error) {
		e := item.(Entity)
		return e, datastore.Get(c, e.EntityKey(), e)
	}
}

func DatastoreDeleter(c appengine.Context) Handler {
	return func(item interface{}) (interface{}, error) {
		e := item.(Entity)
		return e, datastore.Delete(c, e.EntityKey())
	}
}

func DatastoreUpdater(c appengine.Context) Handler {
	return func(item interface{}) (interface{}, error) {
		e := item.(Entity)
		key, err := datastore.Put(c, e.EntityKey(), e)
		if err != nil {
			return e, err
		}

		e.SetEntityKey(key)
		return e, nil
	}
}

func CacheSetter(c appengine.Context) Handler {
	return func(item interface{}) (interface{}, error) {
		e := item.(Entity)
		return e, memcache.JSON.Set(c, NewCacheItem(e.(Cacheable)))
	}
}

func CacheLoader(c appengine.Context) Handler {
	return func(item interface{}) (interface{}, error) {
		e := item.(Entity)
		cacheable, ok := e.(Cacheable)

		if !ok {
			return e, nil
		}

		_, err := memcache.JSON.Get(c, cacheable.CacheID(), cacheable)

		if err == nil {
			return e, Done
		}

		if err == memcache.ErrCacheMiss {
			return e, nil
		}

		return e, err
	}
}

func CacheDeleter(c appengine.Context) Handler {
	return func(item interface{}) (interface{}, error) {
		e := item.(Entity)
		cacheable, ok := e.(Cacheable)

		if !ok {
			return e, nil
		}

		if err := memcache.Delete(c, cacheable.CacheID()); err != nil && err != memcache.ErrCacheMiss {
			return e, err
		}

		return e, nil
	}
}

func DatastoreThrottler(cds *CachedDatastore) Handler {
	return func (item interface{}) (interface{}, error) {
		e := item.(Entity)
		cacheable, ok := e.(Cacheable)

		if !ok {
			return e, nil
		}

		throttledCacheID := "throttled" + cacheable.CacheID()
		_, err := memcache.Get(cds.context, throttledCacheID)
		throttled := err != memcache.ErrCacheMiss

		if throttled {
			return e, Done
		}

		if cds.rateLimit > 0 {
			memcache.Set(cds.context, &memcache.Item{
				Key:        throttledCacheID,
				Expiration: cds.rateLimit,
				Value:      []byte(throttledCacheID),
			})
		}

		return e, nil
	}
}

func CacheMissQuerier(ds *CachedDatastore) Handler {
	return func(item interface{}) (interface{}, error) {
		e := item.(Entity)
		queryable, isQueryable := e.(CacheMissQueryable)
		needQueryFallback := e.EntityKey() == nil && isQueryable

		// In case the given cacheable is also queryable,
		// the data is retrieved by executing the
		// cache miss query provided by the entity
		if needQueryFallback {
			return e, ds.Query(queryable.CacheMissQuery()).Result(e)
		}

		return e, nil
	}
}

func DatastoreBatchLoader(c appengine.Context) Handler {
	return func(slice interface{}) (interface{}, error) {
		s := reflect.ValueOf(slice)
		keys, err := ResolveAllKeys(c, s)
		if err != nil {
			return slice, err
		}
		return slice, datastore.GetMulti(c, keys, slice)
	}
}

func DatastoreBatchCreator(c appengine.Context) Handler {
	return func(slice interface{}) (interface{}, error) {
		s := reflect.ValueOf(slice)
		keys := make([]*datastore.Key, s.Len())
		for i := 0; i < s.Len(); i++ {
			key, err := NewKey(c, s.Index(i).Interface().(Entity))
			if err != nil {
				return slice, err
			}
			keys[i] = key
		}

		keys, err := datastore.PutMulti(c, keys, slice)
		if err != nil {
			return slice, err
		}

		for i, key := range keys {
			s.Index(i).Interface().(Entity).SetEntityKey(key)
		}

		return slice, nil
	}
}

func CacheBatchSetter(c appengine.Context) Handler {
	return func(slice interface{}) (interface{}, error) {
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
		return slice, memcache.JSON.SetMulti(c, items)
	}
}

func CacheBatchLoader(cds *CachedDatastore) Handler {
	return func(slice interface{}) (interface{}, error) {
		s := reflect.ValueOf(slice)

		// Creates a slice of cache keys and
		// makes sure every entity in the given
		// slice is appx.Cacheable
		cacheKeys := make([]string, s.Len())
		for i, _ := range cacheKeys {
			cacheable, ok := s.Index(i).Interface().(Cacheable)
			if !ok {
				return slice, ErrNonCacheableEntity
			}
			cacheKeys[i] = cacheable.CacheID()
		}

		items, err := memcache.GetMulti(cds.context, cacheKeys)
		if err != nil {
			return slice, err
		}

		// No entity is cached
		if len(items) == 0 {
			return slice, nil
		}

		nonCachedEntities := make([]Entity, s.Len() - len(items))
		j := 0
		for i := 0; i < s.Len(); i++ {
			entity := s.Index(i).Interface().(Cacheable)
			// At this point we are safe to assume the entity
			// implements appx.Cacheable since it was already
			// verified while creating the slice of cache keys
			cacheID := entity.CacheID()
			item, itemCached := items[cacheID]

			if !itemCached {
				nonCachedEntities[j] = entity
				j++
			} else {
				if err := json.Unmarshal(item.Value, &entity); err != nil {
					return slice, err
				}
			}
		}

		// All entities are cached
		if len(items) == s.Len() {
			return slice, Done
		}

		if len(nonCachedEntities) > 0 {
			return nonCachedEntities, nil
		}

		return slice, Done
	}
}

func SliceValidator() Handler {
	return func(slice interface{}) (interface{}, error) {
		s := reflect.ValueOf(slice)

		if s.Kind() != reflect.Slice {
			return slice, datastore.ErrInvalidEntityType
		}

		return slice, nil
	}
}

