package main

import (
	"fmt"
	"logbook/config/deployment"
	"logbook/internal/logger"
	"logbook/tools/lbctl/objectives"
	"os"
)

type Run func(l *logger.Logger) error

var commands = map[string]Run{
	"objectives": objectives.Run,
}

func dispatch() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("needs more arguments. run: lbctl help")
	}

	arg := os.Args[1]
	os.Args = os.Args[1:]
	command, ok := commands[arg]
	if !ok {
		return fmt.Errorf("command %q doesn't exist. run: lbctl help", arg)
	}
	l := logger.New(&deployment.Config{Environment: "local"}, "lbctl")
	err := command(l)
	if err != nil {
		return fmt.Errorf("%s: %w", arg, err)
	}
	return nil
}

func main() {
	if err := dispatch(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
