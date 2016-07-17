package threaded

import (
	"ntoolkit/registry"
	"ntoolkit/registry/simple"
	"reflect"
	"sync"
)

// Registry is mutex wrapper around simple.Registry
type Registry struct {
	lock     *sync.Mutex
	registry registry.Registry
}

// New returns a new Registry instance.
func New() *Registry {
	return &Registry{
		lock:     &sync.Mutex{},
		registry: simple.New()}
}

// Register registers a type to a type factory.
// The factory should return an instance of the given type when called.
func (registry *Registry) Register(T interface{}, factory func(registry.Registry) (interface{}, error)) error {
	registry.lock.Lock()
	err := registry.registry.Register(T, factory)
	registry.lock.Unlock()
	return err
}

// Clear clears any existing binding to T and discards any instance.
// Anyone already using the instance will retain their instance.
func (registry *Registry) Clear(T interface{}) error {
	registry.lock.Lock()
	err := registry.registry.Clear(T)
	registry.lock.Unlock()
	return err
}

// Bind should reflect over the given instance 'target', and bind to
// any public properties which are of a known type T.
func (registry *Registry) Bind(target interface{}) error {
	registry.lock.Lock()
	err := registry.registry.Bind(target)
	registry.lock.Unlock()
	return err
}

// Resolve directly resolves a single interface type.
func (registry *Registry) Resolve(T reflect.Type) (interface{}, error) {
	registry.lock.Lock()
	instance, err := registry.registry.Resolve(T)
	registry.lock.Unlock()
	return instance, err
}
