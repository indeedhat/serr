package serr

import (
	"encoding/json"
	"strings"
)

// Wrap an error in a structured context
func Wrap(err error, formatter ...ContextFormatter) StructuredError {
	return Context{}.
		Wrap(err, formatter...)
}

// StructuredError provides an error type with structured data
type StructuredError struct {
	err error
	ctx Context
	fmt ContextFormatter
}

// Error implements the error interface
func (se StructuredError) Error() string {
	var builder strings.Builder

	builder.WriteString("Error: ")
	builder.WriteString(se.err.Error())
	builder.WriteByte('\n')
	builder.WriteString("Context: ")
	builder.WriteString(se.fmt.Render(se.ctx))

	return builder.String()
}

// AddContext adds structured data to the error context for the provided key
func (se StructuredError) AddContext(key string, vals ...any) StructuredError {
	se.ctx[key] = append(se.ctx[key], vals...)
	return se
}

// AddContextMap adds a map of structured data to the error context
func (se StructuredError) AddContextMap(ctx Context) StructuredError {
	for k, v := range ctx {
		se.ctx[k] = append(se.ctx[k], v...)
	}
	return se
}

// Context returns the structured data map for the Structured error
func (se StructuredError) Context() Context {
	return se.ctx
}

// Context provides a container for structured data
type Context map[string][]any

// Wrap an error with the current context
func (c Context) Wrap(err error, formatter ...ContextFormatter) StructuredError {
	var f ContextFormatter
	if len(formatter) > 0 {
		f = formatter[0]
	} else {
		f = JsonFormatter{}
	}

	ctx := make(Context)
	for k, v := range c {
		ctx[k] = append(ctx[k], v...)
	}

	return StructuredError{
		err: err,
		ctx: ctx,
		fmt: f,
	}
}

// ContextFormatter is used by the StructuredError when formatting as a string,
// It transforms the structured data map into a string to be appended to the error string
type ContextFormatter interface {
	Render(Context) string
}

// JsonFormatter formats the errors structured data as a json string
type JsonFormatter struct{}

// Render implements the ContextFormatter interface
func (jf JsonFormatter) Render(ctx Context) string {
	d, _ := json.Marshal(ctx)
	return string(d)
}
