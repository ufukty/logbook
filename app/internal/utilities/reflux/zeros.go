package reflux

import (
	"reflect"
	"strings"
)

func findZeroValues(v reflect.Value, path []string) (errs []string) {
	switch t := v.Type(); t.Kind() {
	case reflect.Pointer:
		errs = append(errs, findZeroValues(reflect.Indirect(v), append(path, t.Name()))...)
	case reflect.Struct:
		var fields = t.NumField()
		for i := 0; i < fields; i++ {
			fv := v.Field(i)
			errs = append(errs, findZeroValues(fv, append(path, fv.Type().Name()))...)
		}
	default:
		if v.IsZero() {
			errs = append(errs, strings.Join(append(path, t.Name()), "."))
		}
	}
	return
}

// returns a list of errors for zero-value fields recursively found
func FindZeroValues(src any) []string {
	return findZeroValues(reflect.ValueOf(src), []string{})
}
