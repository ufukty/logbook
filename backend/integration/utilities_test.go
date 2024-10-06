package integration

import (
	"bytes"
	"fmt"
	"logbook/internal/utilities/slicew/lines"
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
		fmt.Println(lines.Prefix(stderr.String(), "        "))
		fmt.Println("    stdout:")
		fmt.Println(lines.Prefix(stdout.String(), "        "))
		os.Exit(1)
	}

	return stdout.String()
}
