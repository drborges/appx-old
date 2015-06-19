package ds

import (
	"appengine"
	"appengine/datastore"
	"reflect"
)

type pageIterator struct {
	query   *datastore.Query
	context appengine.Context
	nextCursor  datastore.Cursor
	prevCursor datastore.Cursor
}

func NewPagesIterator(q *datastore.Query, c appengine.Context) Iterator {
	return &pageIterator{
		query:   q,
		context: c,
	}
}

// TODO refactor this mess :~
func (this *pageIterator) LoadNext(slice interface{}) error {
	sv := reflect.ValueOf(slice)
	if sv.Kind() != reflect.Ptr || sv.IsNil() || sv.Elem().Kind() != reflect.Slice {
		return datastore.ErrInvalidEntityType
	}
	sv = sv.Elem()

	elemType := sv.Type().Elem()
	if elemType.Kind() != reflect.Ptr || elemType.Elem().Kind() != reflect.Struct {
		return datastore.ErrInvalidEntityType
	}

	iter := this.query.Run(this.context)
	for {
		dstValue := reflect.New(elemType.Elem())
		dst := dstValue.Interface()

		entity, ok := dst.(Entity)
		if !ok {
			return datastore.ErrInvalidEntityType
		}

		key, err := iter.Next(entity)
		if err == datastore.Done {
			cursor, cursorErr := iter.Cursor()
			if cursorErr != nil {
				return cursorErr
			}
			this.prevCursor = this.nextCursor
			this.nextCursor = cursor
			this.query = this.query.Start(this.nextCursor)
			return err
		}

		if err != nil {
			return err
		}

		entity.SetKey(key)
		sv.Set(reflect.Append(sv, dstValue))
	}

	return nil
}

func (this *pageIterator) HasNext() bool {
	return this.prevCursor.String() != this.nextCursor.String()
}

func (this *pageIterator) Cursor() string {
	return this.nextCursor.String()
}
