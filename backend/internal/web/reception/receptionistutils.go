package reception

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"time"

	"logbook/config/deployment"
	"logbook/internal/logger/colors"
	"logbook/internal/web/captured"
)

func summarize(deplcfg *deployment.Config, r *http.Request) string {
	if deplcfg.Environment == "local" {
		return fmt.Sprintf("%s %s %s (%s, %s)",
			colors.Green(r.Method),
			colors.Yellow(r.URL.Path),
			colors.Red(r.Proto),
			colors.Blue(r.Host),
			colors.Magenta(r.RemoteAddr),
		)
	}
	return fmt.Sprintf("%s %s %s (%s, %s)",
		r.Method,
		r.URL.Path,
		r.Proto,
		r.Host,
		r.RemoteAddr,
	)
}

func summarizeW(deplcfg *deployment.Config, w *captured.ResponseWriter, t time.Time) string {
	if deplcfg.Environment == "local" {
		return fmt.Sprintf("%s %s %s bytes",
			colors.Magenta(w.Status),
			colors.Green(time.Since(t)),
			colors.Cyan(w.Header().Get("Content-Length")),
		)
	}
	sc, written := w.Status()
	format := "%d %s %s bytes"
	if !written {
		format = "%d* %s %s bytes"
	}
	return fmt.Sprintf(format, sc, time.Since(t), w.Header().Get("Content-Length"))
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
