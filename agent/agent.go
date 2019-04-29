package agent

import (
	"os"
	"syscall"

	"github.com/gin-gonic/gin"
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

const StateUnknown = "Unknown"

// Agent agent to maninpulate jobs
type Agent struct {
	store IStore
}

func New(store IStore) *Agent {
	return &Agent{store: store}
}

func (a *Agent) Run(addr ...string) error {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong pong",
		})
	})

	r.GET("/inspect", inspect(a.store))
	r.GET("/start", start(a.store))
	r.GET("/notify", notify(a.store))

	r.GET("/kill", func(c *gin.Context) {
		qs := c.Request.URL.Query()
		qs.Set("signal", "9")
		c.Request.URL.RawQuery = qs.Encode()
		notify(a.store)(c)
	})

	r.GET("/down", func(c *gin.Context) {
		qs := c.Request.URL.Query()
		qs.Set("signal", "15")
		c.Request.URL.RawQuery = qs.Encode()
		notify(a.store)(c)
	})

	return r.Run(addr...) // listen and serve on 0.0.0.0:8080
}

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
