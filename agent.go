package main

import (
	"os"
	"syscall"

	gops "github.com/keybase/go-ps"
	gopsutil "github.com/shirou/gopsutil/process"
)

// D    uninterruptible sleep (usually IO)
// R    running or runnable (on run queue)
// S    interruptible sleep (waiting for an event to complete)
// T    stopped, either by a job control signal or because it is being traced.
// W    paging (not valid since the 2.6.xx kernel)
// X    dead (should never be seen)
// Z    defunct ("zombie") process, terminated but not reaped by its parent.

// For BSD formats and when the stat keyword is used, additional characters may be displayed:
// <    high-priority (not nice to other users)
// N    low-priority (nice to other users)
// L    has pages locked into memory (for real-time and custom IO)
// s    is a session leader
// l    is multi-threaded (using CLONE_THREAD, like NPTL pthreads do)
// +    is in the foreground process group.

// Process a process
type Process struct {
	Pid        int
	PPid       int
	Executable string
	Path       string
	State      string
}

const StateUnknown = "Unknown"

func getState(pid int) string {
	p, err := gopsutil.NewProcess(int32(pid))
	if err != nil {
		return StateUnknown
	}
	state, err := p.Status()
	if err != nil {
		return StateUnknown
	}
	return state
}

// Inspect inspect a process
func Inspect(pid int) (Process, error) {
	p, err := gops.FindProcess(pid)
	if err != nil {
		return Process{}, err
	}

	path, _ := p.Path()
	return Process{
		Pid:        pid,
		PPid:       p.PPid(),
		Executable: p.Executable(),
		Path:       path,
		State:      getState(pid),
	}, nil
}

// List list processes
func List() ([]Process, error) {
	ps, err := gops.Processes()
	if err != nil {
		return []Process{}, err
	}

	list := []Process{}
	for _, p := range ps {
		pid := p.Pid()
		path, _ := p.Path()
		list = append(list, Process{
			Pid:        p.Pid(),
			PPid:       p.PPid(),
			Path:       path,
			Executable: p.Executable(),
			State:      getState(pid),
		})
	}

	return list, err
}

// Start start a binary with args
func Start(name string, args []string) (int, error) {
	procAttr := new(os.ProcAttr)
	procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
	p, err := os.StartProcess(name, args, procAttr)
	if err != nil {
		return -1, err
	}
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
