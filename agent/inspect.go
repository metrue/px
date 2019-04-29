package agent

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

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

		c.JSON(200, gin.H{
			"data":    string(data),
			"message": "ok",
		})
	}
}
