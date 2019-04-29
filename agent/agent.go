package agent

import (
	"os"
	"syscall"

	"github.com/gin-gonic/gin"
	gops "github.com/keybase/go-ps"
)

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
