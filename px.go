package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/urfave/cli"
)

const endpoint = "http://127.0.0.1:8080"

func request(path string, qs map[string]string) (string, error) {
	req, _ := http.NewRequest("GET", endpoint+path, nil)
	client := &http.Client{Timeout: time.Second * 10}
	q := req.URL.Query()
	for k, v := range qs {
		q.Set(k, v)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func main() {
	app := cli.NewApp()
	app.Name = "px"
	app.Usage = "manipulate processes like a boss"
	app.Version = "0.7.0"

	app.Commands = []cli.Command{
		{
			Name:  "ping",
			Usage: "ping px daemon",
			Action: func(c *cli.Context) error {
				res, err := request("/ping", map[string]string{})
				if err != nil {
					return err
				}
				fmt.Println(res)
				return nil
			},
		},
		{
			Name:  "start",
			Usage: "start process",
			Action: func(c *cli.Context) error {
				cmd := c.Args().First()
				res, err := request(
					"/start",
					map[string]string{
						"cmd": cmd,
					},
				)
				if err != nil {
					return err
				}
				fmt.Println(res)
				return nil
			},
		},
		{
			Name:  "inspect",
			Usage: "inspect a processes",
			Action: func(c *cli.Context) error {
				num := c.Args().First()
				res, err := request(
					"/inspect",
					map[string]string{
						"pid": num,
					},
				)
				if err != nil {
					return err
				}
				fmt.Println(res)
				return nil
			},
		},
		{
			Name:  "kill",
			Usage: "kill a process",
			Action: func(c *cli.Context) error {
				num := c.Args().First()
				res, err := request(
					"/kill",
					map[string]string{
						"pid": num,
					},
				)
				if err != nil {
					return err
				}
				fmt.Println(res)
				return nil
			},
		},
		{
			Name:  "down",
			Usage: "terminate a process",
			Action: func(c *cli.Context) error {
				num := c.Args().First()
				res, err := request(
					"/down",
					map[string]string{
						"pid": num,
					},
				)
				if err != nil {
					return err
				}
				fmt.Println(res)
				return nil
			},
		},
		{
			Name:  "notify",
			Usage: "notify a process with signal",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "signal, s",
					Usage: "signal number",
				},
			},
			Action: func(c *cli.Context) error {
				pid := c.Args().First()
				signal := c.String("signal")
				res, err := request(
					"/notify",
					map[string]string{
						"pid":    pid,
						"signal": signal,
					},
				)
				if err != nil {
					return err
				}
				fmt.Println(res)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
