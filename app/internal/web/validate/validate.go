package validate

import (
	"fmt"
	"reflect"
)

type Validator interface {
	Validate() error
}

func All(vs map[string]Validator) error {
	for k, v := range vs {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("invalid value for %q: %q", k, v)
		}
	}
	return nil
}

// it uses field tags for "json" to report validation errors
func RequestFields(i any) error {
	v := reflect.Indirect(reflect.ValueOf(i))
	t := v.Type()
	fs := v.NumField()
	for i := 0; i < fs; i++ {
		fv := v.Field(i)
		if val, ok := fv.Interface().(Validator); ok {
			ft := t.Field(i)
			if err := val.Validate(); err != nil {
				if k, ok := ft.Tag.Lookup("json"); ok {
					return fmt.Errorf("invalid value for %q: %q", k, v)
				} else {
					panic(fmt.Sprintf("%q.%q doesn't have %q key in field tags", t.Name(), ft.Name, "json"))
				}
			}
		}
	}
	return nil
}
