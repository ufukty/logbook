package logger

import (
	"crypto/sha512"
	"fmt"
	"log"
)

func init() {
	log.SetFlags(log.LstdFlags | log.LUTC)
}

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

type Logger struct {
	name string
}

func NewLogger(name string) *Logger {
	sum := sha512.Sum512([]byte(name))
	var color = colors[sum[len(sum)-1]%12]
	return &Logger{
		name: fmt.Sprintf("%s%s:%s", color, name, colorReset),
	}
}

func (l Logger) Print(args ...any) {
	log.Print(append([]any{l.name}, args...)...)
}

func (l Logger) Println(args ...any) {
	log.Println(append([]any{l.name}, args...)...)
}

func (l Logger) Printf(args ...any) {
	if format, ok := args[0].(string); ok {
		log.Printf("%s "+format, append([]any{l.name}, args[1:]...)...)
	} else {
		panic("First argument Logger.Printf call should be the format string")
	}
}

func (l Logger) Fatal(args ...any) {
	log.Fatal(append([]any{l.name}, args...)...)
}

func (l Logger) Fatalln(args ...any) {
	log.Fatalln(append([]any{l.name}, args...)...)
}

func (l Logger) Fatalf(args ...any) {
	if format, ok := args[0].(string); ok {
		log.Fatalf("%s "+format, append([]any{l.name}, args[1:]...)...)
	} else {
		panic("First argument Logger.Fatalf call should be the format string")
	}
}
