package summarizer

import (
	"fmt"
	"net/http"
	"time"

	"logbook/internal/web/reception/captured"
)

type colorizer interface {
	Blue(s any) any
	Cyan(s any) any
	Green(s any) any
	Magenta(s any) any
	Red(s any) any
	Yellow(s any) any
}

type albino struct{}

func (albino) Blue(s any) any    { return s }
func (albino) Cyan(s any) any    { return s }
func (albino) Green(s any) any   { return s }
func (albino) Magenta(s any) any { return s }
func (albino) Red(s any) any     { return s }
func (albino) Yellow(s any) any  { return s }

type color struct{}

func (color) Blue(s any) any    { return fmt.Sprintf("\033[34m%v\033[0m", s) }
func (color) Cyan(s any) any    { return fmt.Sprintf("\033[36m%v\033[0m", s) }
func (color) Green(s any) any   { return fmt.Sprintf("\033[32m%v\033[0m", s) }
func (color) Magenta(s any) any { return fmt.Sprintf("\033[35m%v\033[0m", s) }
func (color) Red(s any) any     { return fmt.Sprintf("\033[31m%v\033[0m", s) }
func (color) Yellow(s any) any  { return fmt.Sprintf("\033[33m%v\033[0m", s) }

func newColorizer(colorize bool) colorizer {
	if colorize {
		return color{}
	}
	return albino{}
}

type Summarizer struct {
	c colorizer
}

func New(colorize bool) *Summarizer {
	return &Summarizer{
		c: newColorizer(colorize),
	}
}

func (s Summarizer) Pre(r *http.Request) string {
	return fmt.Sprintf("%s %s %s (%s, %s)",
		s.c.Green(r.Method),
		s.c.Yellow(r.URL.Path),
		s.c.Red(r.Proto),
		s.c.Blue(r.Host),
		s.c.Magenta(r.RemoteAddr),
	)
}

func (s Summarizer) Post(crw *captured.ResponseWriter, start time.Time) string {
	return fmt.Sprintf("%s %s %s",
		s.c.Magenta(crw.StatusRepresentation()),
		s.c.Green(fmt.Sprintf("%dµs", time.Since(start).Microseconds())),
		s.c.Cyan(crw.SizeRepresentation()),
	)
}
