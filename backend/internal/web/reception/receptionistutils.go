package reception

import (
	"fmt"
	"logbook/internal/logger/colors"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

func summarize(r *http.Request) string {
	return fmt.Sprintf("%s %s %s (%s, %s)",
		colors.Green(r.Method),
		colors.Yellow(r.URL.Path),
		colors.Red(r.Proto),
		colors.Blue(r.Host),
		colors.Magenta(r.RemoteAddr),
	)
}

func summarizeW(w *response, t time.Time) string {
	return fmt.Sprintf("%s %s %s bytes",
		colors.Magenta(w.Status),
		colors.Green(time.Since(t)),
		colors.Cyan(w.Header().Get("Content-Length")))
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
