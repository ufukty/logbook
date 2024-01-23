package database

import (
	"bytes"
	"fmt"
	"logbook/internal/utilities/strw"
	"os"
	"os/exec"
)

func run(program string, args ...string) string {
	cmd := exec.Command(program, args...)
	stdout, stderr := bytes.NewBuffer([]byte{}), bytes.NewBuffer([]byte{})
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("running %q: %s\n", cmd.String(), err.Error())
		fmt.Println("    stderr:")
		fmt.Println(strw.IndentLines(stderr.String(), 8))
		fmt.Println("    stdout:")
		fmt.Println(strw.IndentLines(stdout.String(), 8))
		os.Exit(1)
	}

	return stdout.String()
}
