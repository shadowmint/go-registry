# Registry

A reflecting property based IOC library.

## Usage

    import "ntoolkit/registry/threaded"

    type IFoo interface { ... }
    type IBar interface { ... }

    type FooType { ... }
    type BarType { foo IFoo }

    type MyType { Bar IBar }

    ...

    reg := threaded.New()

    reg.Register((*IFoo)(nil), func(_ Registry) (interface{}, error) {
        return FooType{}, nil
    })

    reg.Register((*IBar)(nil), func(R Registry) (interface{}, error) {
        rtn := BarType{nil}

        R.Bind(&rtn) // <-- Manually apply recursive resolution on struct objects

        return rtn, nil
    })

    target := MyType{nil}
    reg.Bind(&target)
    target.Bar.Bar()
