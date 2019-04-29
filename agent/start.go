package agent

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
)

func start(store IStore) func(c *gin.Context) {
	return func(c *gin.Context) {
		rawCmd := c.Query("cmd")
		if rawCmd == "" {
			c.JSON(400, gin.H{
				"message": "cmd is required",
			})
			return
		}

		cmds := strings.Split(rawCmd, " ")
		cmd := exec.Command(cmds[0], cmds[1:]...)
		if err := cmd.Start(); err != nil {
			msg := fmt.Sprintf("start job %s failed: %v", rawCmd, err)
			c.JSON(500, gin.H{
				"message": msg,
			})
			return
		}

		pid := cmd.Process.Pid

		go func() {
			if err := cmd.Wait(); err != nil {
				// TODO we may need to log with different strategy when
				// according exit code, but delete from store by now
				if exiterr, ok := err.(*exec.ExitError); ok {
					if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
						log.Printf("Exit Status: %d", status.ExitStatus())
					}
				} else {
					log.Fatalf("cmd.Wait: %v", err)
				}
			} else {
				ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
				if ws.ExitStatus() == 0 {
					// TODO record it if need, but we just delete from store by
					log.Printf("Exit Status: %d", ws.ExitStatus())
				} else {
					log.Fatalf("Exit Status: %d", ws.ExitStatus())
				}
			}
			store.Delete(strconv.Itoa(pid))
		}()

		// TODO we may need update status when it end
		// log.Printf("Waiting for command to finish...")
		// err = cmd.Wait()
		// log.Printf("Command finished with error: %v", err)
		process := Process{
			Pid:   pid,
			Path:  cmd.Path,
			State: cmd.ProcessState.String(),
		}
		if err := store.Set(strconv.Itoa(pid), process); err != nil {
			msg := fmt.Sprintf("save job info %s failed: %v", rawCmd, err)
			c.JSON(500, gin.H{
				"message": msg,
			})
			return
		}

		msg := fmt.Sprintf("job started with %d", pid)
		c.JSON(200, gin.H{
			"message": msg,
		})
	}
}
