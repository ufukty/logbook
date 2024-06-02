package api

import "path/filepath"

func joinPaths(segment ...Path) Path {
	j := ""
	for _, s := range segment {
		j = filepath.Join(j, string(s))
	}
	return Path(j)
}

func (p Path) Join(n Path) Path {
	return joinPaths(p, n)
}
