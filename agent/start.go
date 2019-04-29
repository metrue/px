package agent

import (
	"fmt"
	"os/exec"
	"strings"

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
		err := cmd.Start()
		if err != nil {
			msg := fmt.Sprintf("start job %s failed: %v", rawCmd, err)
			c.JSON(500, gin.H{
				"message": msg,
			})
			return
		}

		// TODO we may need update status when it end
		// log.Printf("Waiting for command to finish...")
		// err = cmd.Wait()
		// log.Printf("Command finished with error: %v", err)
		pid := cmd.Process.Pid
		process := Process{
			Pid:   pid,
			Path:  cmd.Path,
			State: cmd.ProcessState.String(),
		}
		if err := store.Set(string(pid), process); err != nil {
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
