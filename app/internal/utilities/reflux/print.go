package reflux

import (
	"fmt"
	"reflect"
	"strings"
)

func print(v reflect.Value, paths []string) {
	switch t := v.Type(); t.Kind() {
	case reflect.Pointer:
		print(reflect.Indirect(v), paths)
	case reflect.Struct:
		var nFields = v.Type().NumField()
		for i := 0; i < nFields; i++ {
			print(v.Field(i), append(paths, v.Type().Field(i).Name))
		}
	default:
		fmt.Printf("    %s = %s\n", strings.Join(paths, "."), v)
	}
}

func Print(subj any) {
	print(reflect.ValueOf(subj), []string{})
}
