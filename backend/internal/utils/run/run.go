package run

import (
	"bufio"
	"fmt"
	"logbook/internal/utils/lines"
	"os"
	"os/exec"
	"sync"
)

// it exits if stderr receives anything or command returns error
func ExitAfterStderr(program string, args ...string) string {
	cmd := exec.Command(program, args...)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("could not get the stdout for command %q\n", cmd.String())
		os.Exit(1)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		fmt.Printf("could not get the stderr for command %q\n", cmd.String())
		os.Exit(1)
	}

	var wg sync.WaitGroup
	outputChan := make(chan string)
	stderrReceived := false

	reader := func(pipe *bufio.Reader, isStderr bool) {
		defer wg.Done()
		for {
			line, err := pipe.ReadString('\n')
			outputChan <- line
			if err != nil {
				break
			}
			if isStderr {
				stderrReceived = true
			}
		}
	}

	if err := cmd.Start(); err != nil {
		fmt.Println(fmt.Errorf("could not start the command %q: %w", cmd.String(), err))
		os.Exit(1)
	}

	wg.Add(2)
	go reader(bufio.NewReader(stdoutPipe), false)
	go reader(bufio.NewReader(stderrPipe), true)

	go func() {
		wg.Wait()
		close(outputChan)
	}()

	output := ""
	for line := range outputChan {
		output += line
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Printf("running %q: %s\n", cmd.String(), err.Error())
		fmt.Println(lines.Prefix(output, "    "))
		os.Exit(1)
	} else if stderrReceived {
		fmt.Printf("running %q\n", cmd.String())
		fmt.Println(lines.Prefix(output, "    "))
		os.Exit(1)
	}

	return output
}
