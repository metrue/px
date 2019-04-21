package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/metrue/px/agent"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "px"
	app.Usage = "manipulate processes like a boss"
	app.Version = "0.6.6"

	app.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "start process",
			Action: func(c *cli.Context) error {
				executable := c.Args().First()
				if executable == "" {
					return fmt.Errorf("full path of executable is required")
				}

				args := strings.Split(strings.Trim(executable, " "), " ")
				pid, err := agent.Start(args[0], args)
				if err != nil {
					return err
				}
				log.Printf("pid is %d", pid)
				return nil
			},
		},
		{
			Name:  "list",
			Usage: "list processes",
			Action: func(c *cli.Context) error {
				processes, err := agent.List()
				if err != nil {
					return err
				}

				for _, p := range processes {
					fmt.Println(p)
				}
				return nil
			},
		},
		{
			Name:  "inspect",
			Usage: "inspect a processes",
			Action: func(c *cli.Context) error {
				num := c.Args().First()
				if num == "" {
					return fmt.Errorf("pid is required")
				}

				pid, err := strconv.ParseInt(num, 10, 64)
				p, err := agent.Inspect(int(pid))
				if err != nil {
					return err
				}
				fmt.Println(p)
				return nil
			},
		},
		{
			Name:  "kill",
			Usage: "kill a process",
			Action: func(c *cli.Context) error {
				p := c.Args().First()
				if p == "" {
					return fmt.Errorf("pid is required")
				}

				pid, err := strconv.ParseInt(p, 10, 64)
				if err != nil {
					return err
				}
				if err := agent.Kill(int(pid)); err != nil {
					return fmt.Errorf("process %d could be killed: %v", pid, err)
				}
				log.Printf("process %d was killed", pid)

				return nil
			},
		},
		{
			Name:  "down",
			Usage: "terminate a process",
			Action: func(c *cli.Context) error {
				p := c.Args().First()
				if p == "" {
					return fmt.Errorf("pid is required")
				}

				pid, err := strconv.ParseInt(p, 10, 64)
				if err != nil {
					return err
				}
				if err := agent.Down(int(pid)); err != nil {
					return fmt.Errorf("process %d could be down: %v", pid, err)
				}
				log.Printf("process %d was down", pid)
				return nil
			},
		},
		{
			Name:  "notify",
			Usage: "notify a process with signal",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "signal, s",
					Usage: "signal number",
				},
			},
			Action: func(c *cli.Context) error {
				p := c.Args().First()
				if p == "" {
					return fmt.Errorf("pid is required")
				}

				pid, err := strconv.ParseInt(p, 10, 64)
				if err != nil {
					return err
				}

				signal := c.Int("signal")
				if signal == 0 {
					return fmt.Errorf("signal number is required")
				}

				return agent.Signal(int(pid), signal)
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
