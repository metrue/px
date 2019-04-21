package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/urfave/cli"
)

var agent *Agent

func init() {
	agent = &Agent{}
}

func main() {
	app := cli.NewApp()
	app.Name = "px"
	app.Usage = "manipulate processes like a boss"
	app.Version = "0.0.0"

	app.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "run a binary, eg px start '/bin/sleep' '10000'",
			Action: func(c *cli.Context) error {
				binary := c.Args().First()
				procAttr := new(os.ProcAttr)
				procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
				err := agent.StartProcess(binary, c.Args(), procAttr)
				return err
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
					path, _ := p.Path()
					fmt.Println(p.Pid(), p.PPid(), path)
				}
				return nil
			},
		},
		{
			Name:  "kill",
			Usage: "kill a process",
			Action: func(c *cli.Context) error {
				p := c.Args().First()
				pid, err := strconv.ParseInt(p, 10, 64)
				if err != nil {
					return err
				}
				return agent.Kill(int(pid))
			},
		},
		{
			Name:  "down",
			Usage: "terminate a process",
			Action: func(c *cli.Context) error {
				p := c.Args().First()
				pid, err := strconv.ParseInt(p, 10, 64)
				if err != nil {
					return err
				}
				return agent.Down(int(pid))
			},
		},
		{
			Name:  "notify",
			Usage: "notify a process with signal",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "signal, S",
					Usage: "signal number",
				},
			},
			Action: func(c *cli.Context) error {
				p := c.Args().First()
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
