package logger

import "sync"

var colors = []string{
	"\033[31m", // Red
	"\033[32m", // Green
	"\033[33m", // Yellow
	"\033[34m", // Blue
	"\033[35m", // Magenta
	"\033[36m", // Cyan
	"\033[91m", // LightRed
	"\033[92m", // LightGreen
	"\033[93m", // LightYellow
	"\033[94m", // LightBlue
	"\033[95m", // LightMagenta
	"\033[96m", // LightCyan
}

const colorReset = "\033[0m"

type colormanager struct {
	picks map[string]string
	m     sync.Mutex
}

func (m *colormanager) pickcolor(s string) string {
	m.m.Lock()
	defer m.m.Unlock()
	if color, ok := m.picks[s]; ok {
		return color
	}
	next := colors[len(m.picks)%12]
	m.picks[s] = next
	return next
}

var m = colormanager{picks: map[string]string{}}
