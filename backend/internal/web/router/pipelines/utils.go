package pipelines

import (
	"fmt"
	"net/http"
)

func summarize(r *http.Request) string {
	return fmt.Sprintf("(\033[34m%s\033[0m, \033[35m%s\033[0m) \033[31m%s\033[0m \033[32m%s\033[0m \033[33m%s\033[0m", r.Host, r.RemoteAddr, r.Proto, r.Method, r.URL.Path)
}

func lastsix[S ~string](id S) S {
	return id[max(0, len(id)-6):]
}
