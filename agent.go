package main

import (
	"os"
	"syscall"

	gops "github.com/keybase/go-ps"
)

// List list processes
func List() ([]gops.Process, error) {
	ps, err := gops.Processes()
	if err != nil {
		return []gops.Process{}, err
	}

	return ps, err
}

// Start start a binary with args
func Start(name string, args []string) (int, error) {
	procAttr := new(os.ProcAttr)
	procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
	p, err := os.StartProcess(name, args, procAttr)
	return p.Pid, err
}

// Kill kill a process by pid
func Kill(pid int) error {
	return Signal(pid, 9)
}

// Down termination a process by pid
func Down(pid int) error {
	return Signal(pid, 15)
}

// Signal notify a process with signal
func Signal(pid int, signal int) error {
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return p.Signal(syscall.Signal(signal))
}
