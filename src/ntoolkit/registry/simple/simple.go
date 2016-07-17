package simple

import (
	"container/list"
	"ntoolkit/errors"
	"ntoolkit/registry"
	"ntoolkit/registry/tools"
	"reflect"
)

type record struct {
	T        reflect.Type
	instance interface{}
	factory  func(registry.Registry) (interface{}, error)
}

// Registry is a simple naive binding registry.
type Registry struct {
	records *list.List
}

// New returns a new Registry instance.
func New() *Registry {
	return &Registry{
		records: list.New()}
}

// recordFor returns the record for a given type.
func (registry *Registry) recordFor(T reflect.Type) (*record, *list.Element) {
	for e := registry.records.Front(); e != nil; e = e.Next() {
		val := e.Value.(*record)
		if val.T == T {
			return val, e
		}
	}
	return nil, nil
}

// insertRecordFor creates and returns a record for the given type.
func (registry *Registry) createRecordFor(T reflect.Type, factory func(registry.Registry) (interface{}, error)) *record {
	r := &record{T, nil, factory}
	registry.records.PushBack(r)
	return r
}

// Register registers a type to a type factory.
// The factory should return an instance of the given type when called.
func (registry *Registry) Register(T interface{}, factory func(registry.Registry) (interface{}, error)) error {
	tvalue, err := tools.TypeOf(T)
	if err != nil {
		return errors.Fail(tools.ErrBadType{}, err, "Invalid interface type")
	}
	record, _ := registry.recordFor(tvalue)
	if record != nil {
		return errors.Fail(tools.ErrBindingConflict{}, nil, "An existing binding for %v already exists", tvalue)
	}
	_ = registry.createRecordFor(tvalue, factory)
	return nil
}

// Clear clears any existing binding to T and discards any instance.
// Anyone already using the instance will retain their instance.
func (registry *Registry) Clear(T interface{}) error {
	tvalue, err := tools.TypeOf(T)
	if err != nil {
		return errors.Fail(tools.ErrBadType{}, err, "Invalid interface type")
	}
	_, e := registry.recordFor(tvalue)
	if e != nil {
		registry.records.Remove(e)
	}
	return nil
}

// Bind should reflect over the given instance 'target', and bind to
// any public properties which are of a known type T.
func (registry *Registry) Bind(target interface{}) error {
	_, err := tools.BindProperties(target, func(T reflect.Type) (interface{}, error) {
		instance, err := registry.Resolve(T)
		if err != nil {
			if !errors.Is(err, tools.ErrNoBinding{}) {
				return instance, err
			}
			return nil, nil
		}
		return instance, nil
	})
	return err
}

// Resolve directly resolves a single interface type.
func (registry *Registry) Resolve(T reflect.Type) (interface{}, error) {

	// Check binding
	record, _ := registry.recordFor(T)
	if record == nil {
		return nil, errors.Fail(tools.ErrNoBinding{}, nil, "No binding for %v", T)
	}

	// Create instance?
	if record.instance == nil {
		value, err := record.factory(registry)
		if err != nil {
			return nil, errors.Fail(tools.ErrFactory{}, err, "Failed to create instance of %v", T)
		} else if value == nil {
			return nil, errors.Fail(tools.ErrFactory{}, nil, "Factory returned nil instance of %v", T)
		}
		record.instance = value
	}

	return record.instance, nil
}
