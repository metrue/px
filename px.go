package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "px"
	app.Usage = "manipulate processes like a boss"
	app.Version = "0.6.4"

	app.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "run a binary, eg px start '/bin/sleep' '10000'",
			Action: func(c *cli.Context) error {
				executable := c.Args().First()
				if executable == "" {
					return fmt.Errorf("full path of executable is required")
				}
				pid, err := Start(executable, c.Args())
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
				processes, err := List()
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
				p, err := Inspect(int(pid))
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
				if err := Kill(int(pid)); err != nil {
					return fmt.Errorf("process %d could be killed: %v", pid, err)
				}
				log.Printf("process %d was down", pid)

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
				if err := Down(int(pid)); err != nil {
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

				return Signal(int(pid), signal)
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
