package reception

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
)

func summarize(r *http.Request) string {
	return fmt.Sprintf("(\033[34m%s\033[0m, \033[35m%s\033[0m) \033[31m%s\033[0m \033[32m%s\033[0m \033[33m%s\033[0m", r.Host, r.RemoteAddr, r.Proto, r.Method, r.URL.Path)
}

func lastsix[S ~string](id S) S {
	return id[max(0, len(id)-6):]
}

func funcname(i any) string {
	value := reflect.ValueOf(i)
	if value.Kind() != reflect.Func {
		return "(Not a function)"
	}

	functionPtr := runtime.FuncForPC(reflect.ValueOf(i).Pointer())
	if functionPtr == nil {
		return "(Unknown function)"
	}

	return functionPtr.Name()
}
