package main

import (
	"os"
	"syscall"

	gops "github.com/keybase/go-ps"
)

type Agent struct{}

// List list processes
func (a *Agent) List() ([]gops.Process, error) {
	ps, err := gops.Processes()
	if err != nil {
		return []gops.Process{}, err
	}

	return ps, err
}

// StartProcess start a binary with args
func (a *Agent) StartProcess(
	binary string,
	args []string,
	attr *os.ProcAttr,
) error {
	_, err := os.StartProcess(binary, args, attr)
	return err
}

// Kill kill a process by pid
func (a *Agent) Kill(pid int) error {
	return a.Signal(pid, 9)
}

// Down termination a process by pid
func (a *Agent) Down(pid int) error {
	return a.Signal(pid, 15)
}

// Signal notify a process with signal
func (a *Agent) Signal(pid int, signalNum int) error {
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return p.Signal(int2Signal(signalNum))
}

func int2Signal(s int) os.Signal {
	// TODO support more signal here
	l := map[int]os.Signal{
		1:  syscall.SIGABRT,
		2:  syscall.SIGINT,
		9:  syscall.SIGKILL,
		15: syscall.SIGTERM,
	}
	return l[s]
}

// 		syscall.SIGALRM,
// 		syscall.SIGBUS,
// 		syscall.SIGCHLD,
// 		syscall.SIGCONT,
// 		syscall.SIGFPE,
// 		syscall.SIGILL,
// 		syscall.SIGINT,
// 		syscall.SIGIO,
// 		syscall.SIGIOT,
// 		syscall.SIGKILL,
// 		syscall.SIGPIPE,
// 		syscall.SIGPROF,
// 		syscall.SIGQUIT,
// 		syscall.SIGSEGV,
// 		syscall.SIGSTOP,
// 		syscall.SIGSYS,
// 		syscall.SIGTERM,
// 		syscall.SIGTRAP,
// 		syscall.SIGTSTP,
// 		syscall.SIGTTIN,
// 		syscall.SIGTTOU,
// 		syscall.SIGURG,
// 		syscall.SIGUSR1,
// 		syscall.SIGUSR2,
// 		syscall.SIGVTALRM,
// 		syscall.SIGWINCH,
// 		syscall.SIGXCPU,
// 		syscall.SIGXFSZ,
// 	}
// 	return l[s]
// }
//
// SIGHUP	1	Exit	Hangup
// SIGINT	2	Exit	Interrupt
// SIGQUIT	3	Core	Quit
// SIGILL	4	Core	Illegal Instruction
// SIGTRAP	5	Core	Trace/Breakpoint Trap
// SIGABRT	6	Core	Abort
// SIGEMT	7	Core	Emulation Trap
// SIGFPE	8	Core	Arithmetic Exception
// SIGKILL	9	Exit	Killed
// SIGBUS	10	Core	Bus Error
// SIGSEGV	11	Core	Segmentation Fault
// SIGSYS	12	Core	Bad System Call
// SIGPIPE	13	Exit	Broken Pipe
// SIGALRM	14	Exit	Alarm Clock
// SIGTERM	15	Exit	Terminated
// SIGUSR1	16	Exit	User Signal 1
// SIGUSR2	17	Exit	User Signal 2
// SIGCHLD	18	Ignore	Child Status
// SIGPWR	19	Ignore	Power Fail/Restart
// SIGWINCH	20	Ignore	Window Size Change
// SIGURG	21	Ignore	Urgent Socket Condition
// SIGPOLL	22	Ignore	Socket I/O Possible
// SIGSTOP	23	Stop	Stopped (signal)
// SIGTSTP	24	Stop	Stopped (user)
// SIGCONT	25	Ignore	Continued
// SIGTTIN	26	Stop	Stopped (tty input)
// SIGTTOU	27	Stop	Stopped (tty output)
// SIGVTALRM	28	Exit	Virtual Timer Expired
// SIGPROF	29	Exit	Profiling Timer Expired
// SIGXCPU	30	Core	CPU time limit exceeded
// SIGXFSZ	31	Core	File size limit exceeded
// SIGWAITING	32	Ignore	All LWPs blocked
// SIGLWP	33	Ignore	Virtual Interprocessor Interrupt for Threads Library
// SIGAIO	34	Ignore	Asynchronous I/O
