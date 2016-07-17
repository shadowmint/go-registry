package tools_test

import (
	"ntoolkit/assert"
	"ntoolkit/errors"
	"ntoolkit/registry/tools"
	"reflect"
	"testing"
)

type IFoo interface {
}

type Foo struct {
}

type IBar interface {
}

type Bar struct {
}

type FooBar struct {
	Foo IFoo
	Bar IBar
	foo IFoo
	bar IBar
}

func TestTypeOf(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		tv, err := tools.TypeOf((*IFoo)(nil))
		T.Assert(err == nil)
		T.Assert(tv != nil)

		tv2, _ := tools.TypeOf((*IFoo)(nil))
		T.Assert(tv == tv2)

		tv, err = tools.TypeOf(IFoo(nil))
		T.Assert(tv == nil)
		T.Assert(err != nil)
	})
}

func TestBindProperties(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		fixture := FooBar{}
		count, err := tools.BindProperties(&fixture, func(T reflect.Type) (interface{}, error) {
			if value, err := tools.TypeOf((*IFoo)(nil)); err == nil && T == value {
				return Foo{}, nil
			}
			if value, err := tools.TypeOf((*IBar)(nil)); err == nil && T == value {
				return Bar{}, nil
			}
			return nil, nil
		})
		T.Assert(err == nil)
		T.Assert(count == 2)

		// Public properties are set
		T.Assert(fixture.Foo != nil)
		T.Assert(fixture.Bar != nil)

		// Private properties are not set
		T.Assert(fixture.foo == nil)
		T.Assert(fixture.bar == nil)
	})
}

func TestBindPropertiesWithMissingType(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		fixture := FooBar{}
		count, err := tools.BindProperties(&fixture, func(T reflect.Type) (interface{}, error) {
			if value, err := tools.TypeOf((*IFoo)(nil)); err == nil && T == value {
				return Foo{}, nil
			}
			return nil, nil
		})
		T.Assert(err == nil)
		T.Assert(count == 1)

		// Public properties are set
		T.Assert(fixture.Foo != nil)
		T.Assert(fixture.Bar == nil)

		// Private properties are not set
		T.Assert(fixture.foo == nil)
		T.Assert(fixture.bar == nil)
	})
}

func TestBindPropertiesWithPrimitiveType(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		fixture := 100
		count, err := tools.BindProperties(&fixture, func(T reflect.Type) (interface{}, error) {
			return nil, nil
		})
		T.Assert(errors.Is(err, tools.ErrBadValue{}))
		T.Assert(count == 0)
	})
}

func TestBindPropertiesWithInterfaceType(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		var fixture IFoo = Foo{}
		count, err := tools.BindProperties(&fixture, func(T reflect.Type) (interface{}, error) {
			return nil, nil
		})
		T.Assert(errors.Is(err, tools.ErrBadValue{}))
		T.Assert(count == 0)
	})
}
