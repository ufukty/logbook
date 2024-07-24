package urls

import "strings"

func Join(segments ...string) string {
	url := ""
	for i, segment := range segments {
		if i != 0 && !strings.HasPrefix(segment, "/") {
			url += "/"
		}
		url += segment
	}
	return url
}
