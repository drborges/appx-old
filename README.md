# Features

- [X] Support to model datastore key descriptor
- [X] Support Models with parent keys
- [X] Support to memcache for basic operations (`Load`, `Create`, `Update`, `Delete`)
- [X] Support to Queries runner for fetching multiple results, single result and couting matches
- [X] Support to overwrite cache miss datastore fall back algorithm with a user defined query
- [X] Rich set of Datastore interfaces for decoupling dependency with clients (`Loader`, `Updater`, `Creator` and `Deleter`)
- [X] Support to ItemsIterator
- [X] Support to PagesIterator
- [ ] Support to Iterator as a generator function
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