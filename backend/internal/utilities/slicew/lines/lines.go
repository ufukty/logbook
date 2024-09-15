package lines

import "strings"

func Prefix(src string, prefix string) string {
	lines := strings.Split(src, "\n")
	for i := range lines {
		lines[i] = prefix + lines[i]
	}
	return strings.Join(lines, "\n")
}

// joins string slice which is lines of text with each line prefixed
func Join(src []string, prefix string) string {
	for i := range src {
		src[i] = prefix + src[i]
	}
	return strings.Join(src, "\n")
}
