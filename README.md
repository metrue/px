## px

a simple cli to manipulate process.

## Installation

```
brew tap metrue/homebrew-px
brew install px
```

## Usage

### CLI

```
NAME:
   px - manipulate processes like a boss

USAGE:
   px [global options] command [command options] [arguments...]

VERSION:
   0.6.6

COMMANDS:
     start    start process
     list     list processes
     inspect  inspect a processes
     kill     kill a process
     down     terminate a process
     notify   notify a process with signal
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

For example,
```
$ px start '/bin/sleep 1000'
2019/04/21 19:37:52 pid is 55126

$ px inspect 55126
{
  "pid": 55126,
  "ppid": 1,
  "executable": "sleep",
  "path": "/bin/sleep",
  "state": "S"
  }

$ px kill 55126
2019/04/21 19:39:17 process 55126 was killed
```

### Import as a module

```
import github.com/metrue/px/agent

...

agent.List()                # list process with state
agent.Kill(pid)             # kill a process
agent.Down(pid)             # terminate a process
agent.Inspect(pid)          # inspect a process
agent.Signal(pid, signal)   # notify a process with signal
```
