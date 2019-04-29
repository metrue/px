package agent

import (
	"fmt"
	"os"
	"strconv"
	"syscall"

	"github.com/gin-gonic/gin"
)

func notify(store IStore) func(c *gin.Context) {
	return func(c *gin.Context) {
		pid := c.Query("pid")
		if pid == "" {
			c.JSON(400, gin.H{
				"message": "pid is required",
			})
			return
		}

		pidNum, err := strconv.Atoi(pid)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "pid invalid",
			})
			return
		}

		signal := c.Query("signal")
		if signal == "" {
			c.JSON(400, gin.H{
				"message": "signal is required",
			})
			return
		}

		signalNum, err := strconv.Atoi(signal)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "signal invalid",
			})
			return
		}

		pd, err := store.Get(pid)
		if err != nil {
			msg := fmt.Sprintf("find job with pid=%d failed: %v", pidNum, err)
			c.JSON(500, gin.H{
				"message": msg,
			})
			return
		}

		if pd == nil {
			msg := fmt.Sprintf("no job with pid=%d", pidNum)
			c.JSON(404, gin.H{
				"message": msg,
			})
			return
		}

		p, err := os.FindProcess(pidNum)
		if err != nil {
			msg := fmt.Sprintf("could find job with pid=%d, %v", pidNum, err)
			c.JSON(500, gin.H{
				"message": msg,
			})
			return
		}

		if err := p.Signal(syscall.Signal(signalNum)); err != nil {
			msg := fmt.Sprintf("could notify process (pid %d), %v", pidNum, err)
			c.JSON(500, gin.H{
				"message": msg,
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "ok",
		})
	}
}
