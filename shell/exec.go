package shell

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
)

func CallReverse(command string, args ...string) ([]byte, error) {
	cmd := exec.Command("zsh", "-c", command)

	var outb bytes.Buffer

	var errb bytes.Buffer

	// cmd.Stdin = os.Stdin
	cmd.Stdout = &errb
	cmd.Stderr = &outb

	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, err
	}
	err1 := cmd.Wait()
	if err1 != nil {
		return nil, err1
	}
	s := errb.String()
	var err error

	if s == "" {
		err = nil
	} else {
		err = errors.New(s)
	}

	return outb.Bytes(), err
}

func CallAsync(command string) error {
	go func() {
		cmd := exec.Command("zsh", "-c", command)
		// cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stderr

		if err := cmd.Start(); err != nil {
			return
		}
		err1 := cmd.Wait()
		if err1 != nil {
			return
		}
	}()
	return nil
}

func Call(command string, args ...string) ([]byte, error) {
	cmd := exec.Command("zsh", "-c", command)

	var outb bytes.Buffer

	var errb bytes.Buffer

	// cmd.Stdin = os.Stdin
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, err
	}
	err1 := cmd.Wait()
	if err1 != nil {
		return nil, err1
	}
	s := errb.String()
	var err error

	if s == "" {
		err = nil
	} else {
		err = errors.New(s)
	}

	return outb.Bytes(), err
}
