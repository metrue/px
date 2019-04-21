package agent

import (
	"os"
	"testing"
)

func TestAgent(t *testing.T) {
	t.Run("List", func(t *testing.T) {
		_, err := List()
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Start", func(t *testing.T) {
		name := "/bin/sleep"
		args := []string{"/bin/sleep", "1"}
		pid, err := Start(name, args)
		if err != nil {
			t.Fatal(err)
		}
		if pid == 0 {
			t.Fatal("pid should not be 0")
		}
	})

	t.Run("Signal", func(t *testing.T) {
		pid := -1 // a pid would not exist
		signal := 14
		err := Signal(pid, signal)
		if err == nil {
			t.Fatal("should not find any process")
		}
		pid = os.Getpid()
		err = Signal(pid, signal)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Down", func(t *testing.T) {
		name := "/bin/sleep"
		args := []string{"/bin/sleep", "1"}
		pid, err := Start(name, args)
		if err != nil {
			t.Fatal(err)
		}
		err = Down(pid)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Kill", func(t *testing.T) {
		name := "/bin/sleep"
		args := []string{"/bin/sleep", "1"}
		pid, err := Start(name, args)
		if err != nil {
			t.Fatal(err)
		}

		err = Kill(pid)
		if err != nil {
			t.Fatal(err)
		}
	})
}
