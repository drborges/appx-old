package ds

import (
	"appengine"
	"appengine/datastore"
	"reflect"
)

type QueryRunner struct {
	Context appengine.Context
	Query   *datastore.Query
}

func (this QueryRunner) Count() (int, error) {
	return this.Query.Count(this.Context)
}

func (this QueryRunner) Results(slice interface{}) error {
	keys, err := this.Query.GetAll(this.Context, slice)

	for i, key := range keys {
		this.EntityAt(slice, i).SetKey(key)
	}

	return err
}

func (this QueryRunner) Result(dst Entity) error {
	iter := this.Query.Run(this.Context)
	key, err := iter.Next(dst)

	if err == nil {
		dst.SetKey(key)
	}

	return err
}

func (this QueryRunner) EntityAt(slice interface{}, i int) Entity {
	s := reflect.ValueOf(slice)

	if s.Kind() == reflect.Slice {
		return s.Index(i).Interface().(Entity)
	}

	if s.Kind() == reflect.Ptr && s.Type().Elem().Kind() == reflect.Slice {
		return s.Elem().Index(i).Interface().(Entity)
	}

	panic(datastore.ErrInvalidEntityType)
}