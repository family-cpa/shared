package directive

import (
	"context"
	"reflect"

	"github.com/99designs/gqlgen/graphql"
)

func Ref(ctx context.Context, obj any, next graphql.Resolver, of string) (any, error) {
	return next(ctx)
}

func Entity[T any, ID comparable](id ID) *T {
	if emptyID(id) {
		return nil
	}
	var v T
	f := reflect.ValueOf(&v).Elem().FieldByName("ID")
	if !f.IsValid() || !f.CanSet() {
		panic("directive.Entity: type is missing settable ID")
	}
	val := reflect.ValueOf(id)
	if !val.Type().AssignableTo(f.Type()) {
		panic("directive.Entity: ID type mismatch")
	}
	f.Set(val)
	return &v
}

func EntityPtr[T any, ID comparable](id *ID) *T {
	if id == nil {
		return nil
	}
	return Entity[T](*id)
}

func emptyID[ID comparable](id ID) bool {
	var zero ID
	if id == zero {
		return true
	}
	switch v := any(id).(type) {
	case int:
		return v < 0
	case int8:
		return v < 0
	case int16:
		return v < 0
	case int32:
		return v < 0
	case int64:
		return v < 0
	default:
		return false
	}
}
