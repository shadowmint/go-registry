package tools

import (
	stderr "errors"
	"ntoolkit/errors"
	"reflect"
)

// TypeOf returns the type from an interface
func TypeOf(T interface{}) (rtn reflect.Type, err error) {
	defer func() {
		if praw := recover(); praw != nil {
			perr := stderr.New("Unknown failure")
			switch praw := praw.(type) {
			case string:
				perr = stderr.New(praw)
			case error:
				perr = praw
			}
			err = errors.Fail(ErrBadType{}, perr, "Invalid type")
			rtn = nil
		}
	}()
	return reflect.TypeOf(T).Elem(), nil
}

// BindProperties runs a callback on every public property of the given
// object, and if the object returns a value, assigns it to the object.
// The return is a count of bound values and an error code, if any error.
func BindProperties(target interface{}, handler func(reflect.Type) (interface{}, error)) (int, error) {
	bound := 0

	// Deference pointer types
	V := reflect.ValueOf(target)
	for V.Kind() == reflect.Ptr {
		V = V.Elem()
	}

	// Only structs can be bound
	if V.Kind() != reflect.Struct {
		return 0, errors.Fail(ErrBadValue{}, nil, "Cannot bind properties to non-struct types")
	}

	// Reflect over each property on the type and bind it to the value
	T := V.Type()
	for i := 0; i < T.NumField(); i++ {
		var field = T.Field(i)
		var value = V.Field(i)

		// If possible, bind a new value
		if value.CanSet() && value.IsNil() {
			resolved, rerr := handler(field.Type)
			if rerr != nil {
				return 0, errors.Fail(ErrFactory{}, rerr, "An error occurred during binding")
			} else if resolved != nil {
				rvalue := reflect.ValueOf(resolved)
				value.Set(rvalue)
				bound++
			}
		}
	}

	return bound, nil
}
