package logger

import (
	"fmt"
	"log"
	"logbook/config/deployment"
	"os"
	"slices"
)

type logger interface {
	Print(args ...any)
	Println(args ...any)
	Printf(format string, args ...any)
	Fatal(args ...any)
	Fatalln(args ...any)
	Fatalf(format string, args ...any)
}

type Logger struct {
	name    string
	parent  logger
	deplcfg *deployment.Config
}

func New(deplcfg *deployment.Config, name string) *Logger {
	if deplcfg.Environment == "local" {
		name = fmt.Sprintf("%s%s:%s", m.pickcolor(name), name, colorReset)
	} else {
		name = fmt.Sprintf("%s:", name)
	}
	return &Logger{
		name:    name,
		parent:  log.New(os.Stdout, "", log.LstdFlags|log.LUTC),
		deplcfg: deplcfg,
	}
}

func (l *Logger) Sub(name string) *Logger {
	s := New(l.deplcfg, name)
	s.parent = l
	return s
}

func (l Logger) prefix(args []any) []any {
	return slices.Insert[[]any, any](args, 0, l.name)
}

func (l Logger) Print(args ...any) {
	l.parent.Print(l.prefix(args)...)
}

func (l Logger) Println(args ...any) {
	l.parent.Println(l.prefix(args)...)
}

func (l Logger) Printf(format string, args ...any) {
	l.parent.Printf("%s "+format, l.prefix(args)...)
}

func (l Logger) Fatal(args ...any) {
	l.parent.Fatal(l.prefix(args)...)
}

func (l Logger) Fatalln(args ...any) {
	l.parent.Fatalln(l.prefix(args)...)
}

func (l Logger) Fatalf(format string, args ...any) {
	l.parent.Fatalf("%s "+format, l.prefix(args)...)
}
