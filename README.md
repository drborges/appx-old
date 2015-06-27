# appx

## appx.Datastore

### Create

```go
err := NewDatastore(c).Create(entity)
```

Algorithm:

1. `Assigns` a new key to the entity based on its KeyMetadata() implementation regardless of whether its an incomplete key or not.
2. Creates the entity in datastore

### Load

```go
err := NewDatastore(c).Load(entity)
```

Algorithm:

1. `Resolves` entity's key based on its KeyMetadata() implementation. It fails if the key is incomplete - see [datastore docs](https://cloud.google.com/appengine/docs/go/datastore/reference) for more information.
2. Loads entity's data from datastore

### Update

```go
err := NewDatastore(c).Update(entity)
```

Algorithm:

1. `Resolves` the entity's key based on its KeyMetadata() implementation. It fails if the key is incomplete - see [datastore docs](https://cloud.google.com/appengine/docs/go/datastore/reference) for more information.
2. Updates the entity in datastore

### Delete

```go
err := NewDatastore(c).Delete(entity)
```

Algorithm:

1. `Resolves` the entity's key based on its KeyMetadata() implementation. It fails if the key is incomplete - see [datastore docs](https://cloud.google.com/appengine/docs/go/datastore/reference) for more information.
2. Updates the entity in datastore

## Batch Operations

Batch operations are limited by the same constraints limiting datastore itself, e.g. limit of `1000` reads and limit of `500` writes in a single batch.

### CreateAll

```go
err := NewDatastore(c).CreateAll(entities)
```

Algorithm:

1. For each entity:
  - `Assigns` the entity's key based on its KeyMetadata() implementation.
2. Create the entities in datastore in a single batch of up to `500` entities.

### LoadAll

```go
err := NewDatastore(c).LoadAll(entities)
```

Algorithm:

1. For each entity:
  - `Resolves` the entity's key based on its KeyMetadata() implementation.
2. Loads the entities in datastore in a single batch of up to `1000` entities.

### UpdateAll

```go
err := NewDatastore(c).UpdateAll(entities)
```

Algorithm:

1. For each entity:
  - `Resolves` the entity's key based on its KeyMetadata() implementation.
2. Updates the entities in datastore in a single batch of up to `500` entities.

### DeleteAll

```go
err := NewDatastore(c).DeleteAll(entities)
```

Algorithm:

1. For each entity:
  - `Resolves` the entity's key based on its KeyMetadata() implementation.
2. Deletes the entities in datastore in a single batch of up to `1000` entities. **Need to double check this number**

## Cached datastore

### Load

```go
err := NewDatastore(c).Cached(true).Load(entity)
```

Algorithm:

1. Does the entity implement `appx.Cacheable`?
  - Yes: goto to item 2
  - No: Return an error
2. `Resolves` the entity's key based on its KeyMetadata() implementation.
3. Is the entity cached?
  - Yes: load the entity's data from cache
  - No: goto item 4
4. Loads entity's data from datastore

# Features

- [X] Support to model datastore key descriptor
- [X] Support Models with parent keys
- [X] Support to memcache for basic operations (`Load`, `Create`, `Update`, `Delete`)
- [X] Support to Queries runner for fetching multiple results, single result and couting matches
- [X] Support to overwrite cache miss datastore fall back algorithm with a user defined query
- [X] Rich set of Datastore interfaces for decoupling dependency with clients (`Loader`, `Updater`, `Creator` and `Deleter`)
- [X] Support to ItemsIterator with cursors
- [X] Support to PagesIterator with cursors
- [ ] Support to Iterator as a generator function (channels power \m/)
- [ ] Support memcache on batch operations `CreateAll`, `LoadAll`, `UpdateAll`, `DeleteAll` and perhaps to queries
- [ ] Allow user to define expiration time for cached items
- [ ] Implement: ds.Datastore#Exist(Entity) (bool, error)
- [ ] Implement: ds.Datastore#ExistAll([]Entity) (bool, error)
- [ ] Implement: ds.Datastore#ExistAny([]Entity) (bool, error)
- [ ] Better error handling. Review error cases specially for key resolution algorithm
- [ ] Support struct tags for auto generating KeyMetadata?

# Tagged Structs Prototype:

Struct with a ds.Model embedded and a KeyMetadata method definition that through reflection
extracts key metadata information from the tags in the embedded model 

```golang
type TaggedModel struct {
  Model
}

func (this TaggedModel) KeyMetadata() *KeyMetadata {
	return NewKeyMetadataFromTags(this.Model)
}

func (this TaggedModel) CacheID() string {
    return NewCacheIDFromTags(this.Model)
}
```
