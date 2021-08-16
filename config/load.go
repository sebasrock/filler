package config

import (
	"fmt"
	"reflect"
	"strings"
)

// Config context is represented the basic information need for process a field with a tag
type Config struct {
	providers []ProviderByTag
}

func NewConfig(providers ...ProviderByTag) (*Config, error) {
	if len(providers) < 1 {
		return nil, ErrProviderMandatory
	}
	return &Config{providers: providers}, nil
}

func (c *Config) Marshall(structToPopulate interface{}) error {
	rv := reflect.ValueOf(structToPopulate)

	if rv.Kind() != reflect.Ptr {
		return ErrNotAPointer
	}
	elem := createInfancyPtr(rv.Elem())
	return c.readStruct(elem)
}

func (c *Config) readStruct(value reflect.Value) (err error) {
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		if !field.CanSet() {
			return ErrUnexportedField
		}
		ctx := &Context{
			Name: value.Type().Field(i).Name,
			Tags: value.Type().Field(i).Tag,
		}

		field = createInfancyPtr(field)
		switch field.Kind() {
		case reflect.Struct:
			err = c.readStruct(field)
		case reflect.Interface:
			if field.Elem().Kind() == reflect.Struct {
				fieldNew := reflect.New(field.Elem().Type())
				err = c.readStruct(fieldNew.Elem())
				field.Set(fieldNew)
			}
		default:
			err = c.setField(&field, ctx)
		}

		if err != nil {
			return err
		}
	}

	return err
}

func (c *Config) setField(value *reflect.Value, ctx *Context) (err error) {
	valueToSet, err := GetValueFromProvider(c.providers, ctx)
	if err != nil {
		return err
	}
	return parseValue(value, valueToSet)
}

func parseValue(v *reflect.Value, str string) error {
	var err error = nil
	valueType := v.Type()

	if strings.TrimSpace(str) == "" && valueType.Kind() != reflect.String {
		fmt.Println("load-config: unset value ")
		return nil
	}

	// Special case when the type is a map: we need to make the map
	if valueType.Kind() == reflect.Map {
		v.Set(reflect.MakeMap(valueType))
	}

	kind := valueType.Kind()
	switch {
	case isUnmarshaler(valueType):
		// Special case for Unmarshaler
		err = parseWithUnmarshal(v, str)
	case isDurationField(valueType):
		// Special case for time.Duration
		err = parseDuration(v, str)
	case kind == reflect.Bool:
		err = parseBoolValue(v, str)
	case kind == reflect.Int, kind == reflect.Int8, kind == reflect.Int16, kind == reflect.Int32, kind == reflect.Int64:
		err = parseIntValue(v, str)
	case kind == reflect.Uint, kind == reflect.Uint8, kind == reflect.Uint16, kind == reflect.Uint32, kind == reflect.Uint64:
		err = parseUintValue(v, str)
	case kind == reflect.Float32, kind == reflect.Float64:
		err = parseFloatValue(v, str)
	case kind == reflect.Ptr:
		v.Set(reflect.New(valueType.Elem()))
		field := v.Elem()
		return parseValue(&field, str)
	case kind == reflect.String:
		v.SetString(str)
	default:
		return fmt.Errorf("load-config: kind %v not supported", kind)
	}
	return err
}
