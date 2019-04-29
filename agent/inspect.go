package agent

import (
	"fmt"
	"strconv"

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

func inspect(store IStore) func(c *gin.Context) {
	return func(c *gin.Context) {
		pid := c.Query("pid")
		if pid == "" {
			c.JSON(400, gin.H{
				"message": "pid is required",
			})
			return
		}

		data, err := store.Get(pid)
		if err != nil {
			msg := fmt.Sprintf("query process wit pid = %s failed: %v", pid, err)
			c.JSON(500, gin.H{
				"message": msg,
			})
			return
		}

		if data == nil {
			msg := fmt.Sprintf("no such job wit pid = %s ", pid)
			c.JSON(404, gin.H{
				"message": msg,
			})
			return
		}

		pidNum, err := strconv.Atoi(pid)
		if err != nil {
			msg := fmt.Sprintf("inval pid %s", pid)
			c.JSON(400, gin.H{
				"message": msg,
			})
			return
		}
		p, err := gops.FindProcess(pidNum)
		if err != nil {
			msg := fmt.Sprintf("query process wit pid = %s failed: %v", pid, err)
			c.JSON(500, gin.H{
				"message": msg,
			})
			return
		}

		if p == nil {
			msg := fmt.Sprintf("no such job with pid = %s", pid)
			c.JSON(404, gin.H{
				"message": msg,
			})
			return
		}

		path, _ := p.Path()
		process := Process{
			Pid:        pidNum,
			PPid:       p.PPid(),
			Executable: p.Executable(),
			Path:       path,
			State:      getState(pidNum),
		}

		c.JSON(200, gin.H{
			"data":    process,
			"message": "ok",
		})
	}
}
