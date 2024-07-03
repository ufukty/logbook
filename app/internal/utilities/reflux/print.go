package reflux

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

func print(dst io.Writer, v reflect.Value, paths []string) {
	switch t := v.Type(); t.Kind() {
	case reflect.Pointer:
		print(dst, reflect.Indirect(v), paths)
	case reflect.Struct:
		var nFields = v.Type().NumField()
		for i := 0; i < nFields; i++ {
			ft := t.Field(i)
			if ft.Type.Kind() != reflect.Pointer {
				print(dst, v.Field(i), append(paths, v.Type().Field(i).Name))
			}
		}
	default:
		fmt.Fprintf(dst, "    %s = %s\n", strings.Join(paths, "."), v)
	}
}

func Print(subj any) {
	print(os.Stdout, reflect.ValueOf(subj), []string{})
}

func String(subj any) string {
	b := bytes.NewBuffer([]byte{})
	print(b, reflect.ValueOf(subj), []string{})
	return b.String()
}
