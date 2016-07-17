package threaded_test

import (
	"ntoolkit/registry"
	"ntoolkit/registry/threaded"
	"testing"
)

func TestThreadedRegistry(T *testing.T) {
	registry.TestRegistryWithImpl(T, func() registry.Registry { return threaded.New() })
}
