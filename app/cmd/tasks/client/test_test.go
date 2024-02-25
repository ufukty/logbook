package objectives

import (
	"fmt"
	"os"
	"os/exec"
)

type server struct {
	cmd        *exec.Cmd
	interrupt  chan bool
	safeToExit chan bool
}

func (s *server) Close() {
	defer func() {
		fmt.Println("killing the process which listens :8082 if there is any...")
		killafter := exec.Command(`bash`, `-c`, `'zombie=$(lsof -i :8082 -t) && kill -9 "$zombie"'`)
		killafter.Output()
	}()

	select {
	case <-s.safeToExit:
		return

	default:
		s.interrupt <- true
		<-s.safeToExit
	}
}

func newTestServer() (*server, error) {
	cmd := exec.Command("go", "run", ".", "-config", "../../../platform/testing/config.yml")
	cmd.Dir = "../"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("starting: %w", err)
	}

	s := &server{cmd: cmd}
	s.interrupt = make(chan bool, 1)
	s.safeToExit = make(chan bool, 1)

	go func() {
		e := make(chan error, 1)
		select {
		case e <- cmd.Wait():
			if err := <-e; err != nil {
				fmt.Println(fmt.Errorf("waiting: %w", <-e))
			}
		case <-s.interrupt:
			if err := cmd.Process.Kill(); err != nil {
				fmt.Println(fmt.Errorf("killing: %w", err))
			}
		}
		s.safeToExit <- true
	}()

	return s, nil
}
