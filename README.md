# Features

- [ ] Support to model datastore key descriptor
- [ ] Support Models key resolution through struct tags (`dsx` package)
- [ ] Datastore wrapper with model's key resolution for operations such as `Load`, `Update`, `Create` and `Delete`
- [ ] Rich set of Datastore interfaces for decoupling dependency with clients (`Loader`, `Updater`, `Creator` and `Deleter`)

# TODO

- [ ] Support Models with parent keys (`ds.Model` needs `ParentKey` and `SetParentKey` methods)
- [ ] CachableDatastore with support to custom datastore fallback query
