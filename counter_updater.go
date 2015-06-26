package appx

import (
	"appengine/memcache"
)

type Counter int64

type CounterUpdater struct {
	cds        *CachedDatastore
	cacheable  Cacheable
	cacheOnly  bool
	counterKey string
}

func NewCounterUpdater(cds *CachedDatastore, cacheable Cacheable) *CounterUpdater {
	return &CounterUpdater{
		cds:       cds,
		cacheable: cacheable,
	}
}

func (this *CounterUpdater) Key(counterKey string) *CounterUpdater {
	this.counterKey = counterKey
	return this
}

func (this *CounterUpdater) CacheOnly() *CounterUpdater {
	this.cacheOnly = true
	return this
}

func (this *CounterUpdater) Decrement(counter *Counter) error {
	return this.update(counter, -1)
}

func (this *CounterUpdater) Increment(counter *Counter) error {
	return this.update(counter, 1)
}

func (this *CounterUpdater) Add(counter *Counter, delta int64) error {
	return this.update(counter, delta)
}

func (this *CounterUpdater) update(counter *Counter, delta int64) error {
	key := "Counter," + this.cacheable.KeyMetadata().Kind + "," + this.cacheable.CacheID()
	if this.counterKey != "" {
		key = this.counterKey
	}
	updatedCount, err := memcache.Increment(this.cds.ds.context, key, delta, 0)
	if err != nil {
		return err
	}

	*counter = Counter(updatedCount)

	if !this.cacheOnly {
		return this.cds.Update(this.cacheable)
	}

	return nil
}
