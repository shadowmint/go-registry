package simple_test

import (
	"ntoolkit/registry"
	"ntoolkit/registry/simple"
	"testing"
)

func TestSimpleRegistry(T *testing.T) {
	registry.TestRegistryWithImpl(T, func() registry.Registry { return simple.New() })
}
