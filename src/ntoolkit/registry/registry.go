package registry

import "reflect"

// Registry is a basic interface for any dependency resolver.
type Registry interface {

	// Register registers the interface IFoo to a factory for that interface.
	// The factory function is passed a Registry instance for recursive resolution.
	// @param T the type should be a IFoo(nil) for some interface.
	// @param factory The factory should return an instance of IFoo when called.
	Register(T interface{}, factory func(Registry) (interface{}, error)) error

	// Clear clears any existing binding to T and discards any instance.
	// Anyone already using the instance will retain their instance.
	Clear(T interface{}) error

	// Bind should reflect over the given instance 'target', and bind to
	// any public properties which are of a known type T.
	Bind(target interface{}) error

	// Directly resolve a single interface type.
	Resolve(T reflect.Type) (interface{}, error)
}
