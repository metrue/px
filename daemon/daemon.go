package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"px/agent"

	"github.com/boltdb/bolt"
	"github.com/takama/daemon"
)

const (
	name        = "me.minghe.px"
	description = "px_daemon, a process manipulation agent"
)

var dependencies = []string{"dummy.service"}

var stdlog, errlog *log.Logger

// Service px daemon
type Service struct {
	daemon.Daemon
}

// Manage by daemon commands or run the daemon
func (service *Service) Manage() (string, error) {
	usage := "Usage: px_daemon install | remove | start | stop | status"
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			return service.Start()
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		default:
			return usage, nil
		}
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	go startAgent()

	for {
		select {
		case killSignal := <-interrupt:
			stdlog.Println("Got signal:", killSignal)
			if killSignal == os.Interrupt {
				return "Px Daemon was interruped by system signal", nil
			}
			return "Px Daemon was killed", nil
		}
	}

	return usage, nil
}

// Accept a client connection and collect it in a channel
func startAgent() {
	dbFile := "px.db"
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		errlog.Println("Error: ", err)
		os.Exit(1)
	}
	store := agent.NewStore(db)
	agent := agent.New(store)
	agent.Run()
}

func init() {
	stdlog = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	errlog = log.New(os.Stderr, "", log.Ldate|log.Ltime)
}

func main() {
	srv, err := daemon.New(name, description, dependencies...)
	if err != nil {
		errlog.Println("Error: ", err)
		os.Exit(1)
	}
	service := &Service{srv}
	status, err := service.Manage()
	if err != nil {
		errlog.Println(status, "\nError: ", err)
		os.Exit(1)
	}
	fmt.Println(status)
}
