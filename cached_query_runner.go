package appx

import (
	"appengine"
	"appengine/datastore"
	"time"
	"appengine/memcache"
	"strconv"
)

type CachedQueryRunner struct {
	runner        *QueryRunner
	cacheID       string
	cacheDuration time.Duration
}

func NewCachedQueryRunner(c appengine.Context, q *datastore.Query) *CachedQueryRunner {
	return &CachedQueryRunner{runner: NewQueryRunner(c, q)}
}

func (this *CachedQueryRunner) CachedAs(id string) *CachedQueryRunner {
	this.cacheID = id
	return this
}

func (this *CachedQueryRunner) ExpiresIn(duration time.Duration) *CachedQueryRunner {
	this.cacheDuration = duration
	return this
}

func (this *CachedQueryRunner) Count() (int, error) {
	if this.cacheID == "" {
		return this.runner.Count()
	}

	item, err := memcache.Get(this.runner.context, this.cacheID)
	if err == memcache.ErrCacheMiss {
		count, err := this.runner.Count()
		if err != nil {
			return 0, err
		}

		err = memcache.Set(this.runner.context, &memcache.Item{
			Key: this.cacheID,
			Expiration: this.cacheDuration,
			Value: []byte(strconv.Itoa(count)),
		})
		return count, err
	}
	count, _ := strconv.Atoi(string(item.Value))
	return count, err
}

func (this *CachedQueryRunner) Results(slice interface{}) error {
	if this.cacheID == "" {
		return this.runner.Results(slice)
	}

	_, err := memcache.JSON.Get(this.runner.context, this.cacheID, slice)

	if err == memcache.ErrCacheMiss {
		if err = this.runner.Results(slice); err != nil {
			return err
		}

		return memcache.JSON.Set(this.runner.context, &memcache.Item{
			Key: this.cacheID,
			Expiration: this.cacheDuration,
			Object: slice,
		})
	}

	if err != nil {
		return err
	}

	return nil
}

func (this *CachedQueryRunner) Result(e Entity) error {
	if this.cacheID == "" {
		return this.runner.Result(e)
	}

	_, err := memcache.JSON.Get(this.runner.context, this.cacheID, e)

	if err == memcache.ErrCacheMiss {
		if err = this.runner.Result(e); err != nil {
			return err
		}

		return memcache.JSON.Set(this.runner.context, &memcache.Item{
			Key: this.cacheID,
			Expiration: this.cacheDuration,
			Object: e,
		})
	}

	if err != nil {
		return err
	}

	return nil
}
