package config

import (
	"errors"
	"reflect"
	"time"
)

var (
	// ErrUnexportedField is the error returned by the Init* functions when a field of the config struct is not
	// exported, and the option AllowUnexported is not used.
	ErrUnexportedField = errors.New("load-config: unexported field")
	// ErrNotAPointer is the error returned by the Init* functions when the configuration object is not a pointer.
	ErrNotAPointer = errors.New("load-config: value is not a pointer")
	// ErrProviderMandatory is the error returned by the NewConfig
	ErrProviderMandatory = errors.New("load-config: required at least a provider")
	// ErrInvalidValueKind is the error returned by the Init* functions when the configuration object is not a struct.
	ErrInvalidValueKind = errors.New("load-config: invalid value kind, only works on structs")
)

var (
	durationType    = reflect.TypeOf(new(time.Duration)).Elem()
	unmarshalerType = reflect.TypeOf(new(Unmarshaler)).Elem()
)

type (
	// Unmarshaler is the interface implemented by objects that can unmarshal
	// an environment variable string of themselves.
	Unmarshaler interface {
		Unmarshal(s string) error
	}
)

// context is represent the basic information need for process a field with a tag
type Context struct {
	// name field read
	Name string
	// tags associate to field
	Tags reflect.StructTag
	// tags associate to field
	TagValue string
}

type ProviderByTag struct {
	Provider Provider
	Tag      string
}

type Provider interface {
	Execute(*Context) (string, error)
}
