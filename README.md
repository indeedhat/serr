# Serr - Structured errors for go

## Add contextual data to an error
Wrapping an error to add context is simple
```go
err := errors.New("base error")
se := serr.Wrap(err)

// add a single value
se.AddContext("key", "value")

// Add multiple values under a single key
se.AddContext("key2", "value1", "value2", 3)

// Add a map of values
se.AddContextMap(serr.Context{
    "single": {"data"},
    "multi": {"multiple", "data", "value"},
})

fmt.Print(se)
// Error: base error
// Context: {"key":["value"],"key2":["value1","value2",3],"multi":["multiple","data","value"],"single":["data"]}
```

## Use a shared context for multiple errors
If you want to add shared context to all errors from a specific function/package etc then this can be done
via `serr.Context`
```go
ectx := serr.Context{"base": {"data"}}

err := errors.New("base error")
se := ectx.Wrap(err)

se.AddContext("key", "value")

// Add multiple values under a single key
se.AddContext("key2", "value1", "value2", 3)

// Add a map of values
se.AddContextMap(serr.Context{
    "single": {"data"},
    "multi":  {"multiple", "data", "value"},
})

fmt.Print(se)
// Error: base error
// Context: {"base":["data"],"key":["value"],"key2":["value1","value2",3],"multi":["multiple","data","value"],"single":["data"]}
```

## Custom formatting
By default when formatting the error string json will be used to encode the context data, however
you can create your own formatter to display it as you like

```go

type CustomFormatter struct{}

func (c CustomFormatter) Render(ctx serr.Context) string {
	var builder strings.Builder
	builder.WriteByte('\n')

	for k, vs := range ctx {
		for _, v := range vs {
			builder.WriteString(fmt.Sprintf("\t%s = %v\n", k, v))
		}
	}

	return builder.String()
}

var _ serr.ContextFormatter = (*CustomFormatter)(nil)

...

err := errors.New("base error")
se := serr.Wrap(err, CustomFormatter{})

se.AddContext("key", "value")

// Add multiple values under a single key
se.AddContext("key2", "value1", "value2", 3)

// Add a map of values
se.AddContextMap(serr.Context{
	"single": {"data"},
	"multi":  {"multiple", "data", "value"},
})

fmt.Print(se)

// Error: base error
// Context:
//         key = value
//         key2 = value1
//         key2 = value2
//         key2 = 3
//         single = data
//         multi = multiple
//         multi = data
//         multi = value
```

## Access error context programatically
The error context is just a map so you can access it and make decisions based on how an error fails
```go
func myProc() error {
	return serr.Wrap(errors.New("some error")).
		AddContext("GoodEnough", {false})
}

...

err := myProc()
if goodEnough, found := err.Context()["GoodEnough"]; !found || !goodEnough[0] {
	panic(err)
}
```
