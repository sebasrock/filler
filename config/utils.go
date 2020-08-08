package config

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func GetValueFromProvider(providers []ProviderByTag, ctx *Context) (string, error) {
	var valueToSet string
	for _, provider := range providers {
		if tag, ok := ctx.Tags.Lookup(provider.Tag); ok {
			ctx.TagValue = tag
			valueFromProvider, err := provider.Provider.Execute(ctx)
			if err != nil {
				return "", fmt.Errorf("load-config: unable to execuete proovider for tag  %v. err=%v", provider.Tag, err)
			}
			if strings.TrimSpace(valueFromProvider) != "" {
				valueToSet = valueFromProvider
				break
			}
		}
	}
	return valueToSet, nil
}

func createInfancyPtr(elem reflect.Value) reflect.Value {
	if elem.Kind() == reflect.Ptr {
		if elem.IsNil() {
			elem.Set(reflect.New(elem.Type().Elem()))
		}
		return elem.Elem()
	}
	return elem
}

func isDurationField(t reflect.Type) bool {
	return t.AssignableTo(durationType)
}

func isUnmarshaler(t reflect.Type) bool {
	return t.Implements(unmarshalerType) || reflect.PtrTo(t).Implements(unmarshalerType)
}

func parseBoolValue(v *reflect.Value, str string) error {
	val, err := strconv.ParseBool(str)
	if err != nil {
		return err
	}
	v.SetBool(val)

	return nil
}

func parseIntValue(v *reflect.Value, str string) error {
	val, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return err
	}
	v.SetInt(val)

	return nil
}

func parseUintValue(v *reflect.Value, str string) error {
	val, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return err
	}
	v.SetUint(val)

	return nil
}

func parseFloatValue(v *reflect.Value, str string) error {
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return err
	}
	v.SetFloat(val)

	return nil
}

func parseWithUnmarshal(v *reflect.Value, str string) error {
	var u = v.Addr().Interface().(Unmarshaler)
	return u.Unmarshal(str)
}

func parseDuration(v *reflect.Value, str string) error {
	d, err := time.ParseDuration(str)
	if err != nil {
		return err
	}

	v.SetInt(int64(d))

	return nil
}
