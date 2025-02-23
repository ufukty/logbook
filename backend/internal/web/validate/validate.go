package validate

import (
	"fmt"
	"reflect"
	"strings"
)

type Validator interface {
	Validate() error
}

func All(vs map[string]Validator) error {
	errs := []string{}
	for k, v := range vs {
		if err := v.Validate(); err != nil {
			errs = append(errs, fmt.Sprintf("%s (%s)", k, err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("invalid value(s) for %s", strings.Join(errs, ", "))
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
				return fmt.Errorf("%s.Validate: %s", ft.Name, err)
			}
		}
	}
	return nil
}
