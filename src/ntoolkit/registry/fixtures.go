package registry

import (
	"ntoolkit/assert"
	"ntoolkit/errors"
	"ntoolkit/registry/tools"
	"testing"
)

type iFoo interface {
	Foo() int
}

type iBar interface {
	Bar() int
}

type fooA struct {
	value int
}

type fooB struct {
}

type wantsDeps struct {
	Foo iFoo
}

type wantsDepsHasDeps struct {
	Foo iFoo
	Old iFoo
}

type barWithFoo struct {
	Foo iFoo
}

type withBar struct {
	Bar iBar
}

func (bar barWithFoo) Bar() int {
	return bar.Foo.Foo()
}

func (foo fooA) Foo() int {
	return foo.value
}

func (foo fooB) Foo() int {
	return 1
}

// TestRegistryWithImpl runs a test suite against a specific implementation.
func TestRegistryWithImpl(T *testing.T, factory func() Registry) {
	testRegistryIsValid(T, factory())
	testRegisterCantResolveMissingType(T, factory())
	testRegisterCanResolveBoundType(T, factory())
	testRegisterAndBindType(T, factory())
	testBindDoesNotDestroyExistingBindings(T, factory())
	testRecursiveResolution(T, factory())
}

func testRegistryIsValid(T *testing.T, impl Registry) {
	assert.Test(T, func(T *assert.T) {
		T.Assert(impl != nil)
	})
}

func testRegisterCantResolveMissingType(T *testing.T, impl Registry) {
	assert.Test(T, func(T *assert.T) {
		tvalue, _ := tools.TypeOf((*iFoo)(nil))
		instance, err := impl.Resolve(tvalue)
		T.Assert(instance == nil)
		T.Assert(errors.Is(err, tools.ErrNoBinding{}))
	})
}

func testRegisterCanResolveBoundType(T *testing.T, impl Registry) {
	assert.Test(T, func(T *assert.T) {
		tvalue, _ := tools.TypeOf((*iFoo)(nil))

		T.Assert(impl.Register((*iFoo)(nil), func(_ Registry) (interface{}, error) { return fooA{1}, nil }) == nil)
		instance, err := impl.Resolve(tvalue)
		T.Assert(err == nil)
		T.Assert(instance != nil)

		instance2, err := impl.Resolve(tvalue)
		T.Assert(err == nil)
		T.Assert(instance2 != nil)
		T.Assert(instance == instance2)
	})
}

func testRegisterAndBindType(T *testing.T, impl Registry) {
	assert.Test(T, func(T *assert.T) {
		target := wantsDeps{nil}
		T.Assert(impl.Register((*iFoo)(nil), func(_ Registry) (interface{}, error) { return fooA{1}, nil }) == nil)
		T.Assert(impl.Bind(&target) == nil)
		T.Assert(target.Foo != nil)
		T.Assert(target.Foo.Foo() == 1)
	})
}

func testBindDoesNotDestroyExistingBindings(T *testing.T, impl Registry) {
	assert.Test(T, func(T *assert.T) {
		var value iFoo = fooA{2}
		target := wantsDepsHasDeps{nil, value}
		T.Assert(impl.Register((*iFoo)(nil), func(_ Registry) (interface{}, error) { return fooA{3}, nil }) == nil)
		T.Assert(impl.Bind(&target) == nil)
		T.Assert(target.Foo.Foo() == 3)
		T.Assert(target.Old.Foo() == 2)
	})
}

func testRecursiveResolution(T *testing.T, impl Registry) {
	assert.Test(T, func(T *assert.T) {
		target := withBar{nil}

		T.Assert(impl.Register((*iFoo)(nil), func(_ Registry) (interface{}, error) {
			return fooA{99}, nil
		}) == nil)
		T.Assert(impl.Register((*iBar)(nil), func(R Registry) (interface{}, error) {
			rtn := barWithFoo{nil}
			R.Bind(&rtn)
			return rtn, nil
		}) == nil)

		T.Assert(impl.Bind(&target) == nil)
		T.Assert(target.Bar != nil)
		T.Assert(target.Bar.Bar() == 99)
	})
}
