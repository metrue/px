package main

import "testing"

func TestProcess(t *testing.T) {
	p := Process{
		Pid:        1,
		PPid:       0,
		Executable: "/bin/echo",
		Path:       "/bin/echo",
		State:      "S",
	}
	if p.String() != `{
  "pid": 1,
  "ppid": 0,
  "executable": "/bin/echo",
  "path": "/bin/echo",
  "state": "S"
}` {
		t.Fatal("incorrect string value")
	}
}
