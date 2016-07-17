package tools

// ErrBadType is raised when an operation (eg. TypeOf) is called on an invalid type.
type ErrBadType struct{}

// ErrBadValue is raised when an operation (eg. BindProperties) is called on
// an invalid value; eg. trying to bind properties on an interface.
type ErrBadValue struct{}

// ErrFactory is raised when a factory function returns an error.
// A missing binding is not ErrFactory; this is only raised when the factory
// implementation fails, eg. due to a panic in the service constructor.
type ErrFactory struct{}

// ErrBindingConflict is raised when an attempt is made to create a binding when
// an existing binding exists.
type ErrBindingConflict struct{}

// ErrNoBinding is raised when an attempt is made to resolve an object
// when no binding for that type exists.
type ErrNoBinding struct{}
