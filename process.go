package main

import "encoding/json"

// Process a process
type Process struct {
	Pid        int    `json:"pid"`
	PPid       int    `json:"ppid"`
	Executable string `json:"executable"`
	Path       string `json:"path"`
	State      string `json:"state"`
}

func (p Process) String() string {
	s, _ := json.MarshalIndent(p, "", "  ")
	return string(s)
}
